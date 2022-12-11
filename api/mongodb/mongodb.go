package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
	UUID   string             `bson:"uuid,omitempty"`
}

func Connect() (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(os.Getenv("ME_CONFIG_MONGODB_URL"))
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		fmt.Println("xxx", err.Error())
		return nil, err
	}

	database := client.Database("go-tracing")
	podcastsCollection := database.Collection("podcasts")

	err = podcastsCollection.Drop(context.Background())

	podcast := Podcast{
		Title:  "The Polyglot Developer",
		Author: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
		UUID:   "db417e37-5302-4a9d-8c95-a0eb1ae29ab4",
	}
	_, err = podcastsCollection.InsertOne(context.Background(), podcast)

	//

	return client, nil
}
