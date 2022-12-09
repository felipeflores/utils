package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoPersistence struct {
	logger   *zap.Logger
	client   *mongo.Client
	database string
}

type ConfigMongo struct {
	Host     string
	Database string
	Username string
	Password string
}

func NewMongoPersistence(config ConfigMongo, logger *zap.Logger) (*MongoPersistence, error) {

	credential := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			log.Print(evt.CommandFinishedEvent)
		},
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Host).SetAuth(credential).SetMonitor(cmdMonitor))
	if err != nil {
		return nil, err
	}

	return &MongoPersistence{
		logger:   logger,
		client:   client,
		database: config.Database,
	}, nil
}

func (m *MongoPersistence) InsertOne(ctx context.Context, collection string, t interface{}) {
	c := m.getDatabase().Collection(collection)
	_, err := c.InsertOne(ctx, t)
	if err != nil {
		m.logger.Error(err.Error())
	}
}

func (m *MongoPersistence) getDatabase() *mongo.Database {
	return m.client.Database(m.database)
}

func (m *MongoPersistence) Aggregate(ctx context.Context, collection string, t interface{}, results interface{}) error {
	c := m.getDatabase().Collection(collection)
	cursor, err := c.Aggregate(ctx, t)
	if err != nil {
		return err
	}
	var r []bson.M
	if err = cursor.All(ctx, &r); err != nil {
		return err
	}

	// convert map to json
	jsonString, err := json.Marshal(r)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonString))
	err = json.Unmarshal(jsonString, &results)
	if err != nil {
		return err
	}
	return nil
}
