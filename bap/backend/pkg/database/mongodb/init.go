package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Enum collection names
const (
	WorkerProfileCollection = "workerProfile"
	JobCollection                = "job"
)

// MongoClient structure contains all the database collections and the instance of the database
type MongoClient struct {
	Client                  *mongo.Client
	Database                *mongo.Database
	WorkerProfileCollection *mongo.Collection
	JobCollection                *mongo.Collection
}

var (
	Client            *MongoClient
	dbName            = "onest"
	ConnectionTimeout = 20 * time.Second
	backgroundContext = context.Background()
)

// NewMongoClient Initialize initializes database connection
func NewMongoClient() (*MongoClient, error) {
	client, err := connect()
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoClient{
		Database:                database,
		WorkerProfileCollection: database.Collection(WorkerProfileCollection),
		JobCollection:                database.Collection(JobCollection),
		Client:                  client,
	}, nil
}

func connect() (*mongo.Client, error) {
	if config.Config.DbServer == "" || config.Config.DbUser == "" || config.Config.DbPassword == "" {
		return nil, errors.New("invalid db credentials")
	}

	credential := options.Credential{
		Username: config.Config.DbUser,
		Password: config.Config.DbPassword,
	}

	clientOptions := options.Client().ApplyURI(config.Config.DbServer).SetAuth(credential)

	client, err := mongo.Connect(backgroundContext, clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(backgroundContext, ConnectionTimeout)
	defer cancel()
	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
