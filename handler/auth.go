package handler

import (
	"encoding/json"
	"fmt"
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

	if _, err := util.FindByPhoneNumber(UserInfo.PhoneNumber); err != nil && err.Error() != "user not found" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if existUser, _ := util.FindByPhoneNumber(UserInfo.PhoneNumber); existUser != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    "User already exists",
			StatusCode: http.StatusConflict,
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

	UserInfo.ID = userID

	go func() {
		var otp string = util.GenerateOTP()
		if err := util.SendOTP(UserInfo.PhoneNumber, otp); err != nil {
			log.Printf("Error sending OTP: %v %+v", err, UserInfo)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(schema.FailureResponse{
				Success:    false,
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			})
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(schema.SuccessResponse{
				Success:    true,
				Payload:    userID,
				Message:    "User created successfully",
				StatusCode: http.StatusCreated,
			})
			go util.StoreOTP(UserInfo.PhoneNumber, otp)
		}
	}()
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.AuthInfo

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	otpFromRedis, err := util.FetchOTP(authBody.PhoneNumber)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	if otpFromRedis != authBody.Otp {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    "OTP does not match.",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	token, err := util.GenerateToken(authBody.PhoneNumber)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schema.SuccessResponse{
		Success:    true,
		Payload:    token,
		Message:    "User verified successfully",
		StatusCode: http.StatusOK,
	})
}

func ResendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.AuthInfo

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(schema.FailureResponse{
			Success:    false,
			Message:    "Invalid request body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	otpFromRedis, err := util.FetchOTP(authBody.PhoneNumber)

	if err != nil {
		go func() {
			var otp string = util.GenerateOTP()
			if err := util.SendOTP(authBody.PhoneNumber, otp); err != nil {
				log.Printf("Error in resending OTP: %v %+v", err, authBody)
			} else {
				util.StoreOTP(authBody.PhoneNumber, otp)
			}
		}()
	} else {
		go func() {
			if err := util.SendOTP(authBody.PhoneNumber, otpFromRedis); err != nil {
				log.Printf("Error in resending OTP: %v %+v", err, authBody)
			} else {
				util.StoreOTP(authBody.PhoneNumber, otpFromRedis)
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
			util.StoreOTP(request.PhoneNumber, otp)
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
	refreshToken = refreshToken[1 : len(refreshToken)-1]
	// validating the refresh token
	claims, err := util.ParseToken(refreshToken)
	fmt.Println(err)
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
