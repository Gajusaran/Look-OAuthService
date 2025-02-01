package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nainanisumit/loginOAuth/model"
	"github.com/nainanisumit/loginOAuth/schema"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var otpRequestBody model.AppUser
	// Will change in future
	json.NewDecoder(r.Body).Decode(&otpRequestBody)

	userID, err := util.createUser(otpRequestBody)

	// Will discuss this
	w.WriteHeader(http.StatusCreated)

	if err != nil {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{Success: false, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(schema.UserCreatedSuccessResponse{Success: true, Payload: userID, Message: "User created successfully"})
	}
}
