package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/intelligent-algorithm/pkg/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnectorImpl struct {
	client           *mongo.Client
	connectionString string
	databaseName     string
}

func NewMongoConnector(connectionString, databaseName string) interfaces.MongoConnector {
	return &MongoConnectorImpl{
		connectionString: connectionString,
		databaseName:     databaseName,
	}
}

func (mc *MongoConnectorImpl) Connect(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(mc.connectionString)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(5 * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	mc.client = client
	return nil
}

func (mc *MongoConnectorImpl) Close(ctx context.Context) error {
	if mc.client == nil {
		return nil
	}
	err := mc.client.Disconnect(ctx)
	mc.client = nil
	return err
}

func (mc *MongoConnectorImpl) Ping(ctx context.Context) error {
	if mc.client == nil {
		return fmt.Errorf("MongoDB connection is not established")
	}
	return mc.client.Ping(ctx, nil)
}

func (mc *MongoConnectorImpl) GetClient() *mongo.Client {
	return mc.client
}

func (mc *MongoConnectorImpl) GetDatabase(name string) *mongo.Database {
	if mc.client == nil {
		return nil
	}
	if name == "" {
		name = mc.databaseName
	}
	return mc.client.Database(name)
}

func (mc *MongoConnectorImpl) GetDatabaseName() string {
	return mc.databaseName
}

func (mc *MongoConnectorImpl) GetConnectionString() string {
	return mc.connectionString
}
