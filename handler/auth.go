package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Gajusaran/Look-OAuthService/model"
	"github.com/Gajusaran/Look-OAuthService/schema"
	"github.com/Gajusaran/Look-OAuthService/util"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var otpRequestBody model.AppUser
	// Will change in future
	json.NewDecoder(r.Body).Decode(&otpRequestBody)
	userID, err := util.CreateUser(otpRequestBody)

	// Will discuss this
	w.WriteHeader(http.StatusCreated)

	if err != nil {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{Success: false, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(schema.UserCreatedSuccessResponse{Success: true, Payload: userID, Message: "User created successfully"})
		util.SendOTP(otpRequestBody.PhoneNumber)
	}
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var verifyOtpBody model.UserOtp

	if err := json.NewDecoder(r.Body).Decode(&verifyOtpBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	otpFromRedis, _ := util.FetchOTP(verifyOtpBody.PhoneNumber) // will handle error here also

	if otpFromRedis != verifyOtpBody.Otp {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{
			Success: false,
			Message: "OTP does not match.",
		})
		return
	}

	json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{
		Success: true,
		Message: "OTP matched successfully, user logged in.",
	})
}
