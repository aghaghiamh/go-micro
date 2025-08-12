package adaptor

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URI 	 string
	Username string
	Password string
}

func ConnectToMongo(conf MongoConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(conf.URI)
	clientOptions.SetAuth(options.Credential{
		Username: conf.Username,
		Password: conf.Password,
	})

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Err connecting Mongodb: ", err)
		return nil, err
	}
	return client, nil
}

func Disconnect(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}