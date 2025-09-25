package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/cherry-pick/pkg/analyzer/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAnalyzerService struct {
	connector   core.MongoConnector
	calculator  core.AnalysisCalculator
	aggregator  core.AnalysisAggregator
	validator   core.AnalysisValidator
}

func NewMongoAnalyzerService(
	connector core.MongoConnector,
	calculator core.AnalysisCalculator,
	aggregator core.AnalysisAggregator,
	validator core.AnalysisValidator,
) *MongoAnalyzerService {
	return &MongoAnalyzerService{
		connector:  connector,
		calculator: calculator,
		aggregator: aggregator,
		validator:  validator,
	}
}

func (mas *MongoAnalyzerService) AnalyzeDatabase(ctx context.Context, request core.AnalysisRequest) (*core.AnalysisResult, error) {
	if err := mas.validator.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid analysis request: %w", err)
	}

	if !mas.connector.IsConnected() {
		return nil, fmt.Errorf("MongoDB not connected")
	}

	startTime := time.Now()
	log.Printf("Starting MongoDB analysis for %s", request.DatabaseType)

	collections, err := mas.AnalyzeCollections(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze collections: %w", err)
	}

	dbStats, err := mas.GetDatabaseStats(ctx, request)
	if err != nil {
		log.Printf("Warning: Could not get database stats: %v", err)
	}

	tables := mas.convertCollectionsToTables(collections)
	summary := mas.generateSummary(collections, dbStats)
	insights := mas.generateInsights(collections, dbStats)
	recommendations := mas.generateRecommendations(insights)

	var performance *core.PerformanceMetrics
	if request.Options.IncludePerformance {
		performance, err = mas.GetPerformanceMetrics(ctx, request)
		if err != nil {
			log.Printf("Warning: Could not get performance metrics: %v", err)
		}
	}

	result := &core.AnalysisResult{
		ID:             generateAnalysisID(),
		DatabaseName:   mas.connector.GetDatabaseName(),
		DatabaseType:   request.DatabaseType,
		AnalysisTime:   time.Now(),
		Summary:        summary,
		Tables:         tables,
		Insights:       insights,
		Recommendations: recommendations,
		Performance:    performance,
	}

	log.Printf("MongoDB analysis completed in %v", time.Since(startTime))
	return result, nil
}

func (mas *MongoAnalyzerService) AnalyzeCollections(ctx context.Context, request core.AnalysisRequest) ([]core.MongoCollectionInfo, error) {
	collectionNames, err := mas.GetCollectionNames(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection names: %w", err)
	}

	var collections []core.MongoCollectionInfo
	for _, name := range collectionNames {
		log.Printf("Analyzing collection: %s", name)

		collection, err := mas.AnalyzeCollection(ctx, name, request)
		if err != nil {
			log.Printf("Warning: Failed to analyze collection %s: %v", name, err)
			continue
		}
		collections = append(collections, *collection)
	}

	return collections, nil
}

func (mas *MongoAnalyzerService) AnalyzeCollection(ctx context.Context, collectionName string, request core.AnalysisRequest) (*core.MongoCollectionInfo, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)
	collection := db.Collection(collectionName)

	var stats bson.M
	err := db.RunCommand(ctx, bson.D{{"collStats", collectionName}}).Decode(&stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection stats: %w", err)
	}

	collInfo := &core.MongoCollectionInfo{
		Name:         collectionName,
		LastModified: time.Now(),
	}

	if count, ok := stats["count"].(int64); ok {
		collInfo.DocumentCount = count
	} else if count32, ok := stats["count"].(int32); ok {
		collInfo.DocumentCount = int64(count32)
	}

	if avgSize, ok := stats["avgObjSize"].(float64); ok {
		collInfo.AvgDocSize = int64(avgSize)
	}

	if size, ok := stats["size"].(int64); ok {
		collInfo.TotalSize = size
	} else if size32, ok := stats["size"].(int32); ok {
		collInfo.TotalSize = int64(size32)
	}

	if storageSize, ok := stats["storageSize"].(int64); ok {
		collInfo.StorageSize = storageSize
	} else if storageSize32, ok := stats["storageSize"].(int32); ok {
		collInfo.StorageSize = int64(storageSize32)
	}

	if request.Options.IncludeIndexes {
		indexes, err := mas.GetIndexes(ctx, collectionName, request)
		if err != nil {
			log.Printf("Warning: Could not get indexes for %s: %v", collectionName, err)
		}
		collInfo.Indexes = indexes
	}

	if request.Options.IncludeSchema {
		fields, err := mas.AnalyzeSchema(ctx, collectionName, request)
		if err != nil {
			log.Printf("Warning: Could not analyze schema for %s: %v", collectionName, err)
		}
		collInfo.Fields = fields
	}

	if request.Options.IncludeData {
		sampleDoc, err := mas.getSampleDocument(ctx, collection)
		if err != nil {
			log.Printf("Warning: Could not get sample document for %s: %v", collectionName, err)
		}
		collInfo.SampleDocument = sampleDoc
	}

	return collInfo, nil
}

func (mas *MongoAnalyzerService) GetCollectionNames(ctx context.Context, request core.AnalysisRequest) ([]string, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)

	names, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %w", err)
	}

	if request.Options.MaxCollections > 0 && len(names) > request.Options.MaxCollections {
		names = names[:request.Options.MaxCollections]
	}

	return names, nil
}

func (mas *MongoAnalyzerService) GetCollectionStats(ctx context.Context, collectionName string, request core.AnalysisRequest) (*core.MongoCollectionInfo, error) {
	return mas.AnalyzeCollection(ctx, collectionName, request)
}

func (mas *MongoAnalyzerService) AnalyzeSchema(ctx context.Context, collectionName string, request core.AnalysisRequest) ([]core.MongoFieldInfo, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)
	collection := db.Collection(collectionName)

	sampleSize := request.Options.SampleSize
	if sampleSize <= 0 {
		sampleSize = 100
	}

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(int64(sampleSize)))
	if err != nil {
		return nil, fmt.Errorf("failed to get sample documents: %w", err)
	}
	defer cursor.Close(ctx)

	fieldMap := make(map[string]*core.MongoFieldInfo)
	totalDocs := 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		totalDocs++
		mas.analyzeDocument(doc, fieldMap, "")
	}

	var fields []core.MongoFieldInfo
	for _, field := range fieldMap {
		field.Frequency = field.Frequency / float64(totalDocs)
		fields = append(fields, *field)
	}

	return fields, nil
}

func (mas *MongoAnalyzerService) GetIndexes(ctx context.Context, collectionName string, request core.AnalysisRequest) ([]core.MongoIndexInfo, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)
	collection := db.Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []core.MongoIndexInfo
	for cursor.Next(ctx) {
		var indexDoc bson.M
		if err := cursor.Decode(&indexDoc); err != nil {
			continue
		}

		index := core.MongoIndexInfo{
			UsageStats: core.MongoIndexUsageStats{
				Since: time.Now(),
			},
		}

		if name, ok := indexDoc["name"].(string); ok {
			index.Name = name
		}

		if keys, ok := indexDoc["key"].(bson.M); ok {
			index.Keys = keys
		}

		if unique, ok := indexDoc["unique"].(bool); ok {
			index.IsUnique = unique
		}

		if sparse, ok := indexDoc["sparse"].(bool); ok {
			index.IsSparse = sparse
		}

		if partialFilterExpression, ok := indexDoc["partialFilterExpression"]; ok && partialFilterExpression != nil {
			index.IsPartial = true
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (mas *MongoAnalyzerService) GetDatabaseStats(ctx context.Context, request core.AnalysisRequest) (*core.MongoDatabaseStats, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)

	var stats bson.M
	err := db.RunCommand(ctx, bson.D{{"dbStats", 1}}).Decode(&stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get database stats: %w", err)
	}

	dbStats := &core.MongoDatabaseStats{
		Name: mas.connector.GetDatabaseName(),
	}

	if collections, ok := stats["collections"].(int32); ok {
		dbStats.Collections = int(collections)
	}

	if views, ok := stats["views"].(int32); ok {
		dbStats.Views = int(views)
	}

	if objects, ok := stats["objects"].(int64); ok {
		dbStats.Objects = objects
	} else if objects32, ok := stats["objects"].(int32); ok {
		dbStats.Objects = int64(objects32)
	}

	if avgObjSize, ok := stats["avgObjSize"].(float64); ok {
		dbStats.AvgObjSize = avgObjSize
	}

	if dataSize, ok := stats["dataSize"].(int64); ok {
		dbStats.DataSize = dataSize
	} else if dataSize32, ok := stats["dataSize"].(int32); ok {
		dbStats.DataSize = int64(dataSize32)
	}

	if storageSize, ok := stats["storageSize"].(int64); ok {
		dbStats.StorageSize = storageSize
	} else if storageSize32, ok := stats["storageSize"].(int32); ok {
		dbStats.StorageSize = int64(storageSize32)
	}

	if indexSize, ok := stats["indexSize"].(int64); ok {
		dbStats.IndexSize = indexSize
	} else if indexSize32, ok := stats["indexSize"].(int32); ok {
		dbStats.IndexSize = int64(indexSize32)
	}

	dbStats.TotalSize = dbStats.DataSize + dbStats.IndexSize

	return dbStats, nil
}

func (mas *MongoAnalyzerService) GetPerformanceMetrics(ctx context.Context, request core.AnalysisRequest) (*core.PerformanceMetrics, error) {
	db := mas.connector.GetDatabase().(*mongo.Database)

	var serverStatus bson.M
	err := db.RunCommand(ctx, bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to get server status: %w", err)
	}

	metrics := &core.PerformanceMetrics{}

	if connections, ok := serverStatus["connections"].(bson.M); ok {
		if current, ok := connections["current"].(int32); ok {
			metrics.Connections.Current = int(current)
		}
		if available, ok := connections["available"].(int32); ok {
			metrics.Connections.Available = int(available)
		}
		if totalCreated, ok := connections["totalCreated"].(int32); ok {
			metrics.Connections.TotalCreated = int(totalCreated)
		}
	}

	if opcounters, ok := serverStatus["opcounters"].(bson.M); ok {
		if insert, ok := opcounters["insert"].(int64); ok {
			metrics.Operations.Insert = insert
		} else if insert32, ok := opcounters["insert"].(int32); ok {
			metrics.Operations.Insert = int64(insert32)
		}

		if query, ok := opcounters["query"].(int64); ok {
			metrics.Operations.Query = query
		} else if query32, ok := opcounters["query"].(int32); ok {
			metrics.Operations.Query = int64(query32)
		}

		if update, ok := opcounters["update"].(int64); ok {
			metrics.Operations.Update = update
		} else if update32, ok := opcounters["update"].(int32); ok {
			metrics.Operations.Update = int64(update32)
		}

		if delete, ok := opcounters["delete"].(int64); ok {
			metrics.Operations.Delete = delete
		} else if delete32, ok := opcounters["delete"].(int32); ok {
			metrics.Operations.Delete = int64(delete32)
		}

		if getmore, ok := opcounters["getmore"].(int64); ok {
			metrics.Operations.GetMore = getmore
		} else if getmore32, ok := opcounters["getmore"].(int32); ok {
			metrics.Operations.GetMore = int64(getmore32)
		}

		if command, ok := opcounters["command"].(int64); ok {
			metrics.Operations.Command = command
		} else if command32, ok := opcounters["command"].(int32); ok {
			metrics.Operations.Command = int64(command32)
		}
	}

	return metrics, nil
}

func (mas *MongoAnalyzerService) getSampleDocument(ctx context.Context, collection *mongo.Collection) (map[string]interface{}, error) {
	var doc bson.M
	err := collection.FindOne(ctx, bson.D{}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return doc, nil
}

func (mas *MongoAnalyzerService) analyzeDocument(doc bson.M, fieldMap map[string]*core.MongoFieldInfo, prefix string) {
	for key, value := range doc {
		fieldName := key
		if prefix != "" {
			fieldName = prefix + "." + key
		}

		field, exists := fieldMap[fieldName]
		if !exists {
			field = &core.MongoFieldInfo{
				Name:        fieldName,
				SampleValue: value,
				Frequency:   0,
			}
			fieldMap[fieldName] = field
		}

		field.Frequency++
		field.Type = mas.getFieldType(value)

		if nestedDoc, ok := value.(bson.M); ok {
			mas.analyzeDocument(nestedDoc, fieldMap, fieldName)
		}
	}
}

func (mas *MongoAnalyzerService) getFieldType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.Bool:
		return "bool"
	case reflect.Slice:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "mixed"
	}
}

func (mas *MongoAnalyzerService) convertCollectionsToTables(collections []core.MongoCollectionInfo) []core.TableInfo {
	var tables []core.TableInfo
	for _, coll := range collections {
		table := core.TableInfo{
			Name:         coll.Name,
			RowCount:     coll.DocumentCount,
			Size:         fmt.Sprintf("%d bytes", coll.TotalSize),
			LastModified: coll.LastModified,
		}

		for _, field := range coll.Fields {
			column := core.ColumnInfo{
				Name:     field.Name,
				DataType: field.Type,
				DataProfile: core.DataProfile{
					Quality: 1.0,
				},
			}
			table.Columns = append(table.Columns, column)
		}

		for _, idx := range coll.Indexes {
			index := core.IndexInfo{
				Name:     idx.Name,
				IsUnique: idx.IsUnique,
				Type:     "btree",
			}
			for key := range idx.Keys {
				index.Columns = append(index.Columns, key)
			}
			table.Indexes = append(table.Indexes, index)
		}

		tables = append(tables, table)
	}
	return tables
}

func (mas *MongoAnalyzerService) generateInsights(collections []core.MongoCollectionInfo, stats *core.MongoDatabaseStats) []core.DatabaseInsight {
	var insights []core.DatabaseInsight

	for _, coll := range collections {
		if coll.DocumentCount > 1000000 {
			insight := core.DatabaseInsight{
				Type:           "performance",
				Severity:       "medium",
				Title:          "Large Collection Detected",
				Description:    fmt.Sprintf("Collection '%s' has %d documents, which may impact performance", coll.Name, coll.DocumentCount),
				Suggestion:     "Consider sharding, archiving old data, or optimizing queries",
				AffectedTables: []string{coll.Name},
				MetricValue:    coll.DocumentCount,
			}
			insights = append(insights, insight)
		}

		if len(coll.Indexes) <= 1 && coll.DocumentCount > 10000 {
			insight := core.DatabaseInsight{
				Type:           "performance",
				Severity:       "high",
				Title:          "Missing Indexes on Large Collection",
				Description:    fmt.Sprintf("Collection '%s' has %d documents but only %d indexes", coll.Name, coll.DocumentCount, len(coll.Indexes)),
				Suggestion:     "Consider adding indexes on frequently queried fields",
				AffectedTables: []string{coll.Name},
				MetricValue:    len(coll.Indexes),
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

func (mas *MongoAnalyzerService) generateSummary(collections []core.MongoCollectionInfo, stats *core.MongoDatabaseStats) core.DatabaseSummary {
	var totalRows int64
	var totalColumns int

	for _, coll := range collections {
		totalRows += coll.DocumentCount
		totalColumns += len(coll.Fields)
	}

	healthScore := mas.calculateHealthScore(collections)
	complexityScore := mas.calculateComplexityScore(collections)

	return core.DatabaseSummary{
		TotalTables:     len(collections),
		TotalColumns:    totalColumns,
		TotalRows:       totalRows,
		TotalSize:       fmt.Sprintf("%d bytes", stats.TotalSize),
		HealthScore:     healthScore,
		ComplexityScore: complexityScore,
	}
}

func (mas *MongoAnalyzerService) generateRecommendations(insights []core.DatabaseInsight) []string {
	var recommendations []string

	highPriorityCount := 0
	for _, insight := range insights {
		if insight.Severity == "high" {
			highPriorityCount++
			recommendations = append(recommendations, fmt.Sprintf("Priority: %s", insight.Suggestion))
		}
	}

	if highPriorityCount == 0 {
		recommendations = append(recommendations, "MongoDB database appears to be in good condition with no critical issues")
	}

	return recommendations
}

func (mas *MongoAnalyzerService) calculateHealthScore(collections []core.MongoCollectionInfo) float64 {
	if len(collections) == 0 {
		return 0.0
	}

	var totalScore float64
	for _, coll := range collections {
		collectionScore := 1.0

		if len(coll.Indexes) <= 1 && coll.DocumentCount > 1000 {
			collectionScore -= 0.2
		}

		if coll.DocumentCount > 10000000 && !coll.IsSharded {
			collectionScore -= 0.3
		}

		if collectionScore < 0 {
			collectionScore = 0
		}

		totalScore += collectionScore
	}

	return totalScore / float64(len(collections))
}

func (mas *MongoAnalyzerService) calculateComplexityScore(collections []core.MongoCollectionInfo) float64 {
	complexity := float64(len(collections)) * 0.1

	for _, coll := range collections {
		complexity += float64(len(coll.Fields)) * 0.05
		complexity += float64(len(coll.Indexes)) * 0.1
	}

	return complexity
}
