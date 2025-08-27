package analyzer

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"github.com/intelligent-algorithm/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAnalyzerImpl struct {
	connector interfaces.MongoConnector
}

func NewMongoAnalyzer(connector interfaces.MongoConnector) interfaces.MongoAnalyzer {
	return &MongoAnalyzerImpl{
		connector: connector,
	}
}

func (ma *MongoAnalyzerImpl) AnalyzeDatabase(ctx context.Context) (*types.DatabaseReport, error) {
	log.Println("Starting comprehensive MongoDB analysis...")

	dbStats, err := ma.GetDatabaseStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database stats: %w", err)
	}

	collections, err := ma.AnalyzeCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze collections: %w", err)
	}

	tables := ma.convertCollectionsToTables(collections)
	insights := ma.generateInsights(collections, dbStats)
	summary := ma.generateSummary(collections, dbStats)
	recommendations := ma.generateRecommendations(insights)

	report := &types.DatabaseReport{
		DatabaseName:    ma.connector.GetDatabaseName(),
		DatabaseType:    "mongodb",
		AnalysisTime:    time.Now(),
		Summary:         summary,
		Tables:          tables,
		Insights:        insights,
		Recommendations: recommendations,
	}

	log.Println("MongoDB analysis completed successfully")
	return report, nil
}

func (ma *MongoAnalyzerImpl) AnalyzeCollections(ctx context.Context) ([]types.MongoCollectionInfo, error) {
	collectionNames, err := ma.GetCollectionNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection names: %w", err)
	}

	var collections []types.MongoCollectionInfo
	for _, name := range collectionNames {
		log.Printf("Analyzing collection: %s", name)

		collection, err := ma.AnalyzeCollection(ctx, name)
		if err != nil {
			log.Printf("Warning: Failed to analyze collection %s: %v", name, err)
			continue
		}
		collections = append(collections, *collection)
	}

	return collections, nil
}

func (ma *MongoAnalyzerImpl) AnalyzeCollection(ctx context.Context, collectionName string) (*types.MongoCollectionInfo, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	collection := db.Collection(collectionName)

	var stats bson.M
	err := db.RunCommand(ctx, bson.D{{"collStats", collectionName}}).Decode(&stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection stats: %w", err)
	}

	collInfo := &types.MongoCollectionInfo{
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

	indexes, err := ma.GetIndexes(ctx, collectionName)
	if err != nil {
		log.Printf("Warning: Could not get indexes for %s: %v", collectionName, err)
	}
	collInfo.Indexes = indexes

	fields, err := ma.AnalyzeSchema(ctx, collectionName, 100)
	if err != nil {
		log.Printf("Warning: Could not analyze schema for %s: %v", collectionName, err)
	}
	collInfo.Fields = fields

	sampleDoc, err := ma.getSampleDocument(ctx, collection)
	if err != nil {
		log.Printf("Warning: Could not get sample document for %s: %v", collectionName, err)
	}
	collInfo.SampleDocument = sampleDoc

	return collInfo, nil
}

func (ma *MongoAnalyzerImpl) GetCollectionNames(ctx context.Context) ([]string, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	names, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %w", err)
	}

	return names, nil
}

func (ma *MongoAnalyzerImpl) GetCollectionStats(ctx context.Context, collectionName string) (*types.MongoCollectionInfo, error) {
	return ma.AnalyzeCollection(ctx, collectionName)
}

func (ma *MongoAnalyzerImpl) AnalyzeSchema(ctx context.Context, collectionName string, sampleSize int) ([]types.MongoFieldInfo, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	collection := db.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(int64(sampleSize)))
	if err != nil {
		return nil, fmt.Errorf("failed to get sample documents: %w", err)
	}
	defer cursor.Close(ctx)

	fieldMap := make(map[string]*types.MongoFieldInfo)
	totalDocs := 0

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		totalDocs++
		ma.analyzeDocument(doc, fieldMap, "")
	}

	var fields []types.MongoFieldInfo
	for _, field := range fieldMap {
		field.Frequency = field.Frequency / float64(totalDocs)
		fields = append(fields, *field)
	}

	return fields, nil
}

func (ma *MongoAnalyzerImpl) GetIndexes(ctx context.Context, collectionName string) ([]types.MongoIndexInfo, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	collection := db.Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []types.MongoIndexInfo
	for cursor.Next(ctx) {
		var indexDoc bson.M
		if err := cursor.Decode(&indexDoc); err != nil {
			continue
		}

		index := types.MongoIndexInfo{
			UsageStats: types.MongoIndexUsageStats{
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

func (ma *MongoAnalyzerImpl) GetDatabaseStats(ctx context.Context) (*types.MongoDatabaseStats, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	var stats bson.M
	err := db.RunCommand(ctx, bson.D{{"dbStats", 1}}).Decode(&stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get database stats: %w", err)
	}

	dbStats := &types.MongoDatabaseStats{
		Name: ma.connector.GetDatabaseName(),
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

func (ma *MongoAnalyzerImpl) GetPerformanceMetrics(ctx context.Context) (*types.MongoPerformanceMetrics, error) {
	db := ma.connector.GetDatabase("")
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	var serverStatus bson.M
	err := db.RunCommand(ctx, bson.D{{"serverStatus", 1}}).Decode(&serverStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to get server status: %w", err)
	}

	metrics := &types.MongoPerformanceMetrics{}

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

func (ma *MongoAnalyzerImpl) getSampleDocument(ctx context.Context, collection *mongo.Collection) (map[string]interface{}, error) {
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

func (ma *MongoAnalyzerImpl) analyzeDocument(doc bson.M, fieldMap map[string]*types.MongoFieldInfo, prefix string) {
	for key, value := range doc {
		fieldName := key
		if prefix != "" {
			fieldName = prefix + "." + key
		}

		field, exists := fieldMap[fieldName]
		if !exists {
			field = &types.MongoFieldInfo{
				Name:        fieldName,
				SampleValue: value,
				Frequency:   0,
			}
			fieldMap[fieldName] = field
		}

		field.Frequency++
		field.Type = ma.getFieldType(value)

		if nestedDoc, ok := value.(bson.M); ok {
			ma.analyzeDocument(nestedDoc, fieldMap, fieldName)
		}
	}
}

func (ma *MongoAnalyzerImpl) getFieldType(value interface{}) string {
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

func (ma *MongoAnalyzerImpl) convertCollectionsToTables(collections []types.MongoCollectionInfo) []types.TableInfo {
	var tables []types.TableInfo
	for _, coll := range collections {
		table := types.TableInfo{
			Name:         coll.Name,
			RowCount:     coll.DocumentCount,
			Size:         fmt.Sprintf("%d bytes", coll.TotalSize),
			LastModified: coll.LastModified,
		}

		for _, field := range coll.Fields {
			column := types.ColumnInfo{
				Name:     field.Name,
				DataType: field.Type,
				DataProfile: types.DataProfile{
					Quality: 1.0,
				},
			}
			table.Columns = append(table.Columns, column)
		}

		for _, idx := range coll.Indexes {
			index := types.IndexInfo{
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

func (ma *MongoAnalyzerImpl) generateInsights(collections []types.MongoCollectionInfo, stats *types.MongoDatabaseStats) []types.DatabaseInsight {
	var insights []types.DatabaseInsight

	for _, coll := range collections {
		if coll.DocumentCount > 1000000 {
			insight := types.DatabaseInsight{
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
			insight := types.DatabaseInsight{
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

func (ma *MongoAnalyzerImpl) generateSummary(collections []types.MongoCollectionInfo, stats *types.MongoDatabaseStats) types.DatabaseSummary {
	var totalRows int64
	var totalColumns int

	for _, coll := range collections {
		totalRows += coll.DocumentCount
		totalColumns += len(coll.Fields)
	}

	healthScore := ma.calculateHealthScore(collections)
	complexityScore := ma.calculateComplexityScore(collections)

	return types.DatabaseSummary{
		TotalTables:     len(collections),
		TotalColumns:    totalColumns,
		TotalRows:       totalRows,
		TotalSize:       fmt.Sprintf("%d bytes", stats.TotalSize),
		HealthScore:     healthScore,
		ComplexityScore: complexityScore,
	}
}

func (ma *MongoAnalyzerImpl) generateRecommendations(insights []types.DatabaseInsight) []string {
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

func (ma *MongoAnalyzerImpl) calculateHealthScore(collections []types.MongoCollectionInfo) float64 {
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

func (ma *MongoAnalyzerImpl) calculateComplexityScore(collections []types.MongoCollectionInfo) float64 {
	complexity := float64(len(collections)) * 0.1

	for _, coll := range collections {
		complexity += float64(len(coll.Fields)) * 0.05
		complexity += float64(len(coll.Indexes)) * 0.1
	}

	return complexity
}
