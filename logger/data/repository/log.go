package repository

import (
	"context"
	"log"
	"log-service/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LogRepo struct {
	client *mongo.Client
}

func New(c *mongo.Client) *LogRepo {
	return &LogRepo{
		client: c,
	}
}

func (r *LogRepo) Insert(dlog domain.LogEntry) error {
	collection := r.client.Database("logs").Collection("logs")

	_, inErr := collection.InsertOne(context.TODO(), LogEntry{
		Name:      dlog.Name,
		Data:      dlog.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if inErr != nil {
		log.Println("Error insering into logs: ", inErr)
		return inErr
	}

	return nil
}

func (r *LogRepo) All() ([]*domain.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, fErr := collection.Find(context.TODO(), bson.D{{}}, opts)
	if fErr != nil {
		log.Println("Error finding all docs: ", fErr)
		return nil, fErr
	}
	defer cursor.Close(ctx)

	var logs []*domain.LogEntry

	for cursor.Next(ctx) {
		var logItem LogEntry
		if dErr := cursor.Decode(&logItem); dErr != nil {
			log.Println("Error decoding logEntry into slice: ", dErr)
			return nil, dErr
		}
		logs = append(logs, &domain.LogEntry{
			ID:   logItem.ID,
			Name: logItem.Name,
			Data: logItem.Data,
		})
	}

	return logs, nil
}
