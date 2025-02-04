package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Gajusaran/Look-OAuthService/model"
	"github.com/Gajusaran/Look-OAuthService/schema"
	"github.com/Gajusaran/Look-OAuthService/util"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var otpRequestBody model.AppUser // need to change the name of there sturct variable in each funcation
	// Will change in future
	json.NewDecoder(r.Body).Decode(&otpRequestBody)
	userID, err := util.CreateUser(otpRequestBody)

	// Will discuss this
	w.WriteHeader(http.StatusCreated)

	if err != nil {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{Success: false, Message: err.Error()})
	} else {
		json.NewEncoder(w).Encode(schema.UserCreatedSuccessResponse{Success: true, Payload: userID, Message: "User created successfully"})
		util.SendOTP(otpRequestBody.PhoneNumber, util.GenerateOTP())
	}
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var verifyOtpBody model.UserOtp

	if err := json.NewDecoder(r.Body).Decode(&verifyOtpBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	otpFromRedis, err := util.FetchOTP(verifyOtpBody.PhoneNumber) // will handle error here also

	if err != nil {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{
			Success: false,
			Message: "Invalid otp please try again", // incase otp expire after time limit, this need to be modify
		})
		return
	}

	if otpFromRedis != verifyOtpBody.Otp {
		json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{
			Success: false,
			Message: "OTP does not match.",
		})
		return
	}

	// JWT logic , generate jwt and send
	accessToken, err := util.GenerateAccessToken(verifyOtpBody.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := util.GenerateRefreshToken(verifyOtpBody.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	// Success message with token
	json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{ // need to make response struct so we will send all things in a type of response
		Success: true,
		Message: "OTP matched successfully, user logged in.",
	})
	json.NewEncoder(w).Encode(response)
}

func ResendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var verifyOtpBody model.UserOtp

	if err := json.NewDecoder(r.Body).Decode(&verifyOtpBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	otpFromRedis, err := util.FetchOTP(verifyOtpBody.PhoneNumber)

	if err != nil {
		util.SendOTP(verifyOtpBody.PhoneNumber, util.GenerateOTP()) //if otp is not there redis or expired, will send new otp
		//error handling for sendotp fun
	} else {
		//otp already exits
		util.SendOTP(verifyOtpBody.PhoneNumber, otpFromRedis) //error handling
	}

	json.NewEncoder(w).Encode(schema.UserCreatedFailureResponse{
		Success: true,
		Message: "OTP sent again successfully.",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || strings.TrimSpace(request.PhoneNumber) == "" {
		http.Error(w, "Invalid request: phone number required", http.StatusBadRequest)
		return
	}
	// if user not exists will return for register
	_, err = util.FindByPhoneNumber(request.PhoneNumber)

	if err != nil {
		http.Error(w, "User not found. Please register first.", http.StatusUnauthorized)
		return
	}

	response := model.LoginResponse{
		Success: true,
		Message: "Login request processed. Proceed with OTP verification.",
	}
	json.NewEncoder(w).Encode(response)
	util.SendOTP(request.PhoneNumber, util.GenerateOTP())
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the refresh token from the header with bearer
	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		http.Error(w, "Refresh token is missing", http.StatusBadRequest)
		return
	}

	// Remove the "Bearer " prefix from the token
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	// validating the refresh token
	claims, err := util.ParseToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token using the claims from the refresh token
	accessToken, err := util.GenerateAccessToken(claims.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	response := map[string]string{ //have to check this
		"access_token": accessToken,
	}
	json.NewEncoder(w).Encode(response)
}
