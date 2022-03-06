package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"linnoxlewis/tg-bot-eng-dictionary/internal/config"
	"log"
	"time"
)

type Mongodb struct {
	cfg      *config.Config
	client   *mongo.Client
	Database *mongo.Database
}

const dbUri = "mongodb://%s:%s/"

func Init(ctx context.Context, cfg *config.Config) *Mongodb {
	uri := fmt.Sprintf(dbUri, cfg.GetMongoHost(), cfg.GetMongoPort())
	clientOptions := options.Client().ApplyURI(uri).
		SetAuth(options.Credential{
			AuthSource: cfg.GetMongoDatabase(),
			Username:   cfg.GetMongoUser(),
			Password:   cfg.GetMongoPwd(),
		})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Mongodb{
		client:   client,
		Database: client.Database(cfg.GetMongoDatabase()),
		cfg:      cfg,
	}
}

func (m *Mongodb) GetClient() *mongo.Client {
	return m.client
}
func (m *Mongodb) CreateIndex(collectionName string, field string, unique bool) bool {

	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(unique),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.Database.Collection(collectionName)
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
