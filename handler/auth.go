package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Gajusaran/Look-OAuthService/model"
	"github.com/Gajusaran/Look-OAuthService/schema"
	"github.com/Gajusaran/Look-OAuthService/util"
)

func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var UserInfo model.AppUser

	if err := json.NewDecoder(r.Body).Decode(&UserInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	userID, err := util.CreateUser(UserInfo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schema.SuccessResponse{
		Success:    true,
		Payload:    userID,
		Message:    "User created successfully",
		StatusCode: http.StatusCreated,
	})

	go func() {
		var otp string = util.GenerateOTP()
		if err := util.SendOTP(UserInfo.PhoneNumber, otp); err != nil {
			log.Printf("Error sending OTP: %v %+v", err, UserInfo)
		} else {
			util.StoreOTP(UserInfo.phoneNumber, otp)
		}
	}()

}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.AuthInfo

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	otpFromRedis, err := util.FetchOTP(authBody.PhoneNumber) 

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success: false,
			Message: err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if otpFromRedis != authBody.Otp {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success: false,
			Message: "OTP does not match.",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
    
	token,err:=util.GenerateToken(authBody.phoneNumber)

	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success: false,
			Message: err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schema.SuccessResponse{ 
		Success: true,
		Payload: token,
		Message: "User verified successfully",
        StatusCode: http.StatusOK,
	})
}

func ResendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.UserOtp

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	otpFromRedis, err := util.FetchOTP(authBody.PhoneNumber)

	if err != nil {
		go func() {
			var otp string = util.GenerateOTP()
			if err := util.SendOTP(authBody.PhoneNumber, otp); err != nil {
				log.Printf("Error in resending OTP: %v %+v", err, authBody)
			} else {
				util.StoreOTP(authBody.phoneNumber, otp)
			}
		}()
	} else {
		go func() {
			if err := util.SendOTP(authBody.PhoneNumber, otpFromRedis); err != nil {
				log.Printf("Error in resending OTP: %v %+v", err, authBody)
			} else {
				util.StoreOTP(authBody.phoneNumber, otpFromRedis)
			}
		}()
	}
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
	go func() {
		var otp string = util.GenerateOTP()
		if err := util.SendOTP(request.PhoneNumber, otp); err != nil {
			log.Printf("Error sending OTP: %v %+v", err, request)
		} else {
			util.StoreOTP(request.phoneNumber, otp)
		}
	}()
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
