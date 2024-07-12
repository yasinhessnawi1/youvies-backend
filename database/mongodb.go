package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var Client *mongo.Client

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not found in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	log.Println("Connected to MongoDB!")
}

func InsertItem(item interface{}, title string, collectionName string) error {
	collection := Client.Database("youvies").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exist, err := IfItemExists(bson.M{"title": title}, collectionName)
	if err != nil {
		return err
	}
	if !exist {
		_, err := collection.InsertOne(ctx, item)
		if err != nil {
			return err
		}
	} else {
		return errors.New("item already exists")

	}
	return nil
}
func DeleteItem(filter interface{}, collectionName string) error {
	collection := Client.Database("youvies").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document found to delete")
	}

	return nil
}
func EditItem(filter interface{}, update interface{}, collectionName string) error {
	collection := Client.Database("youvies").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no document found to update")
	}

	return nil
}
func IfItemExists(filter interface{}, collectionName string) (bool, error) {
	collection := Client.Database("youvies").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func FindItem(filter interface{}, collectionName string, result interface{}) error {
	collection := Client.Database("youvies").Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return err
	}

	return nil
}