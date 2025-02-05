package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoDatabaseName = "User"
const collectionName = "UserInfo"

var Collection *mongo.Collection

func init() {
	fmt.Println("intlize")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Changing context from TODO to Background because
	// database connection is top level operation in application
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING")))
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(mongoDatabaseName).Collection(collectionName)

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
