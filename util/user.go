package util

import (
	"context"
	"log"

	"github.com/Gajusaran/Look-OAuthService/database"
	"github.com/Gajusaran/Look-OAuthService/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(otpRequestBody model.AppUser) (primitive.ObjectID, error) {

	userCreated, err := database.Collecation.InsertOne(context.TODO(), otpRequestBody)

	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	// Will replace this from user payload
	return userCreated.InsertedID.(primitive.ObjectID), nil
}
