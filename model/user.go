package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppUser struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	PhoneNumber string             `json:"phn" bson:"phn"`
	UserGender  string             `json:"gender"`
}

type AuthInfo struct {
	Otp         string `json:"otp"`
	PhoneNumber string `json:"phn"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phn"`
}
