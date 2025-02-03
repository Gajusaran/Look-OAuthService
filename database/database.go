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

const dbName = "User"
const colName = "UserInfo"

var Collecation *mongo.Collection

func init() {
	fmt.Println("intlize")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connecationString := os.Getenv("MONGODB_CONNECTION_STRING")

	clientOptions := options.Client().ApplyURI(connecationString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	Collecation = client.Database(dbName).Collection(colName)

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
