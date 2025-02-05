package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppUser struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	PhoneNumber string             `json:"phn" bson:"phn"`
}

type AuthInfo struct {
	Otp         string `json:"otp"`
	PhoneNumber string `json:"phn"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phn"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
