package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppUser struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	PhoneNumber string             `json:"phn"`
}

type UserOtp struct {
	Otp         string `json:"otp"`
	PhoneNumber string `json:"phn"`
}

type LoggedUser struct {
	PhoneNumber string `json:"phn"`
}
