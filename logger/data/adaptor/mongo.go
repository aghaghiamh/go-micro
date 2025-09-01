package adaptor

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	DB           string
	MongoAddress string
	Username     string
	Password     string
}

func ConnectToMongo(conf MongoConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s/%s?authSource=admin&tls=true", conf.MongoAddress, conf.DB))
	clientOptions.SetAuth(options.Credential{
		Username: conf.Username,
		Password: conf.Password,
	})

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Err connecting Mongodb: ", err)
		return nil, err
	}

	log.Println("✅ Successfully connected to MongoDB!")
	return client, nil
}

func Disconnect(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}
