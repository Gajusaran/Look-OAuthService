package util

import (
	"encoding/json"
	"net/http"

	"github.com/loginOAuth/schema"
)

func GetfailureJsonResponse(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(schema.FailureResponse{
		Success:    false,
		StatusCode: statusCode,
	})
}

func GetSuccessJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(schema.SuccessResponse{
		Success:    true,
		Payload:    payload,
		StatusCode: statusCode,
	})
}
