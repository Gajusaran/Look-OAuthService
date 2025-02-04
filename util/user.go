package util

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Gajusaran/Look-OAuthService/database"
	"github.com/Gajusaran/Look-OAuthService/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(otpRequestBody model.AppUser) (primitive.ObjectID, error) {
	//Prevents operations from hanging forever if an external system (like MongoDB) is slow or unresponsive. Helps control resource usage and improve system performance.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //Always call the cancel() function after the operation is complete to clean up resources:
	userCreated, err := database.Collecation.InsertOne(ctx, otpRequestBody)

	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	// Will replace this from user payload
	return userCreated.InsertedID.(primitive.ObjectID), nil
}

func FindByPhoneNumber(phoneNumber string) (*model.AppUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.AppUser
	filter := bson.M{"phn": phoneNumber}

	err := database.Collecation.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { // spcail way of handling error if is matched with mongoerr will use coustom error handling here
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
