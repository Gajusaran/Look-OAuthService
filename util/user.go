package util

import (
	"context"
	"log"

	"github.com/nainanisumit/loginOAuth/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createUser(otpRequestBody model.AppUser) (primitive.ObjectID, error) {

	userCreated, err := collection.InsertOne(context.TODO(), otpRequestBody)

	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	// Will replace this from user payload
	return userCreated.InsertedID.(primitive.ObjectID), nil
}
