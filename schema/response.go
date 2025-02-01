package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCreatedSuccessResponse struct {
	Success bool               `json:"success"`
	Payload primitive.ObjectID `json:"data"`
	Message string             `json:"message"`
}

type UserCreatedFailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
