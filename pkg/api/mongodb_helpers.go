package api

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/cherry-pick/pkg/intelligence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Server) getMongoCollectionData(service *intelligence.Service, collectionName string, page, limit int) ([]map[string]interface{}, int64, error) {
	mongoService := service.GetMongoService()
	if mongoService == nil {
		return nil, 0, fmt.Errorf("not a MongoDB service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connector := mongoService.GetConnector()
	if connector == nil {
		return nil, 0, fmt.Errorf("MongoDB connector not available")
	}

	client := connector.GetClient()
	if client == nil {
		return nil, 0, fmt.Errorf("MongoDB client not available")
	}

	database := connector.GetDatabase("")
	collection := database.Collection(collectionName)

	totalCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	skip := int64((page - 1) * limit)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var documents []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, fmt.Errorf("failed to decode document: %w", err)
		}

		if id, ok := doc["_id"]; ok {
			doc["_id"] = fmt.Sprintf("%v", id)
		}

		documents = append(documents, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	return documents, totalCount, nil
}

func (s *Server) getMongoCollectionStats(service *intelligence.Service, collectionName string) (*CollectionStatsResponse, error) {
	mongoService := service.GetMongoService()
	if mongoService == nil {
		return nil, fmt.Errorf("not a MongoDB service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connector := mongoService.GetConnector()
	if connector == nil {
		return nil, fmt.Errorf("MongoDB connector not available")
	}

	database := connector.GetDatabase("")
	collection := database.Collection(collectionName)

	docCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	sampleSize := int64(100)
	if docCount < sampleSize {
		sampleSize = docCount
	}

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(sampleSize))
	if err != nil {
		return nil, fmt.Errorf("failed to sample documents: %w", err)
	}
	defer cursor.Close(ctx)

	fieldStats := make(map[string]*FieldStats)
	var documents []map[string]interface{}

	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		documents = append(documents, doc)

		for fieldName, value := range doc {
			if fieldName == "_id" {
				continue
			}

			if _, exists := fieldStats[fieldName]; !exists {
				fieldStats[fieldName] = &FieldStats{
					Name:         fieldName,
					Count:        0,
					UniqueCount:  0,
					NullCount:    0,
					SampleValues: make([]string, 0),
				}
			}

			stat := fieldStats[fieldName]
			stat.Count++

			if value == nil {
				stat.NullCount++
			} else {
				stat.Type = getValueType(value)

				if len(stat.SampleValues) < 5 {
					stat.SampleValues = append(stat.SampleValues, fmt.Sprintf("%v", value))
				}
			}
		}
	}

	for fieldName, stat := range fieldStats {
		uniqueValues := make(map[string]bool)
		for _, doc := range documents {
			if value, exists := doc[fieldName]; exists && value != nil {
				uniqueValues[fmt.Sprintf("%v", value)] = true
			}
		}
		stat.UniqueCount = int64(len(uniqueValues))
	}

	var fields []FieldStats
	for _, stat := range fieldStats {
		fields = append(fields, *stat)
	}

	indexes, err := s.getMongoIndexes(collection, ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get indexes: %v\n", err)
		indexes = []IndexStats{}
	}

	response := &CollectionStatsResponse{
		CollectionName: collectionName,
		DocumentCount:  docCount,
		Fields:         fields,
		Indexes:        indexes,
	}

	return response, nil
}

func (s *Server) searchMongoCollection(service *intelligence.Service, collectionName, query string) ([]map[string]interface{}, error) {
	mongoService := service.GetMongoService()
	if mongoService == nil {
		return nil, fmt.Errorf("not a MongoDB service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connector := mongoService.GetConnector()
	database := connector.GetDatabase("")
	collection := database.Collection(collectionName)

	filter := bson.M{
		"$or": []bson.M{
			{"$text": bson.M{"$search": query}},
			// Fallback to regex search on string fields
		},
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetLimit(50))
	if err != nil {
		filter = bson.M{}
		cursor, err = collection.Find(ctx, filter, options.Find().SetLimit(20))
		if err != nil {
			return nil, fmt.Errorf("failed to search collection: %w", err)
		}
	}
	defer cursor.Close(ctx)

	var documents []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		if id, ok := doc["_id"]; ok {
			doc["_id"] = fmt.Sprintf("%v", id)
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

func (s *Server) getMongoIndexes(collection *mongo.Collection, ctx context.Context) ([]IndexStats, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []IndexStats
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			continue
		}

		indexStat := IndexStats{
			Name:     fmt.Sprintf("%v", index["name"]),
			IsUnique: false,
		}

		if key, ok := index["key"].(bson.M); ok {
			for fieldName := range key {
				indexStat.Keys = append(indexStat.Keys, fieldName)
			}
		}

		if unique, ok := index["unique"].(bool); ok {
			indexStat.IsUnique = unique
		}

		indexes = append(indexes, indexStat)
	}

	return indexes, nil
}

func getValueType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "number"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	case time.Time:
		return "date"
	case []interface{}:
		return "array"
	case map[string]interface{}, bson.M:
		return "object"
	default:
		t := reflect.TypeOf(v)
		if t != nil {
			return t.String()
		}
		return "unknown"
	}
}
