package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/loginOAuth/logger"
	"github.com/loginOAuth/model"
	"github.com/loginOAuth/util"
	"github.com/sirupsen/logrus"
)

func Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var UserInfo model.AppUser

	if err := json.NewDecoder(r.Body).Decode(&UserInfo); err != nil {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"body.request": r.Body,
		}).Error("Invalid request body")
		return
	}

	if _, err := util.FindByPhoneNumber(UserInfo.PhoneNumber); err != nil && err.Error() != "user not found" {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"mongo.error": err.Error(),
		}).Error("Internal Server Error/Mongo Error")
		return
	}

	if existUser, _ := util.FindByPhoneNumber(UserInfo.PhoneNumber); existUser != nil {
		util.GetfailureJsonResponse(w, http.StatusConflict)
		logger.Logger.WithFields(logrus.Fields{
			"mongo.error.conflict": "conflict",
		}).Warn("User already exists with", existUser)
		return
	}

	userID, err := util.CreateUser(UserInfo)

	if err != nil {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		logger.Logger.WithFields(logrus.Fields{
			"mongo.error": err.Error(),
		}).Error(err.Error())
		return
	}

	UserInfo.ID = userID
	var otp string = util.GenerateOTP()

	if err := util.StoreOTP(UserInfo.PhoneNumber, otp); err != nil {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		return
	}

	if err := util.SendOTP(UserInfo.PhoneNumber, otp); err != nil {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
	} else {
		util.GetSuccessJsonResponse(w, http.StatusCreated, UserInfo)
		logger.Logger.WithFields(logrus.Fields{
			"userID": UserInfo.ID,
			"phone":  UserInfo.PhoneNumber,
			"name":   UserInfo.Name,
			"gender": UserInfo.UserGender,
		}).Info("User registered successfully")
	}
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.AuthInfo

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"body.request": r.Body,
		}).Error("Invalid request body")
		return
	}

	otpFromRedis := util.FetchOTP(authBody.PhoneNumber)

	if len(otpFromRedis) == 0 {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		return
	}

	if otpFromRedis != authBody.Otp {
		logger.Logger.WithFields(logrus.Fields{
			"otp.redis":        otpFromRedis,
			"request.body.otp": authBody.Otp,
		}).Error("Otp does not match for phone number, ", authBody.PhoneNumber)
		util.GetfailureJsonResponse(w, http.StatusUnauthorized)
		return
	}

	token := util.GenerateToken(authBody.PhoneNumber)

	if len(token) == 0 {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		return
	}

	util.GetSuccessJsonResponse(w, http.StatusOK, token)
	logger.Logger.WithFields(logrus.Fields{
		"otp": otpFromRedis,
	}).Info("Otp Verified for phone number ", authBody.PhoneNumber)
}

func ResendOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authBody model.AuthInfo

	if err := json.NewDecoder(r.Body).Decode(&authBody); err != nil {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"body.request": r.Body,
		}).Error("Invalid request body")
		return
	}

	otpFromRedis := util.FetchOTP(authBody.PhoneNumber)

	if len(otpFromRedis) == 0 {
		var otp string = util.GenerateOTP()
		if err := util.StoreOTP(authBody.PhoneNumber, otp); err != nil {
			util.GetfailureJsonResponse(w, http.StatusInternalServerError)
			return
		}
		if err := util.SendOTP(authBody.PhoneNumber, otp); err != nil {
			util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		}
	} else {
		if err := util.SendOTP(authBody.PhoneNumber, otpFromRedis); err != nil {
			util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"body.request": r.Body,
		}).Error("Invalid request body")
		return
	}

	if _, err := util.FindByPhoneNumber(request.PhoneNumber); err != nil && err.Error() != "user not found" {
		util.GetfailureJsonResponse(w, http.StatusBadRequest)
		logger.Logger.WithFields(logrus.Fields{
			"mongo.error": err.Error(),
		}).Error("Internal Server Error/Mongo Error")
		return
	}

	existUser, _ := util.FindByPhoneNumber(request.PhoneNumber)
	if existUser == nil {
		util.GetfailureJsonResponse(w, http.StatusUnauthorized)
		logger.Logger.WithFields(logrus.Fields{
			"mongo.error.unauthorised": "Unauthorised",
		}).Warn("User not found with phone number ", request.PhoneNumber)
		return
	}

	var otp string = util.GenerateOTP()

	if err := util.StoreOTP(request.PhoneNumber, otp); err != nil {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		return
	}

	if err := util.SendOTP(request.PhoneNumber, otp); err != nil {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
	} else {
		util.GetSuccessJsonResponse(w, http.StatusOK, existUser)
		logger.Logger.WithFields(logrus.Fields{
			"phone": request.PhoneNumber,
		}).Info("User Logged In successfully")
	}
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the refresh token from the header with bearer
	refreshToken := r.Header.Get("Authorization")

	if len(refreshToken) == 0 {
		util.GetfailureJsonResponse(w, http.StatusInternalServerError)
		logger.Logger.WithFields(logrus.Fields{
			"request.body": r.Body,
		}).Error("Refresh Token is missing")
		return
	}

	// Remove the "Bearer " prefix from the token
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	refreshToken = refreshToken[1 : len(refreshToken)-1]
	// validating the refresh token
	claims, err := util.ParseToken(refreshToken)

	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate a new access token using the claims from the refresh token
	accessToken := util.GenerateAccessToken(claims.PhoneNumber)
	if len(accessToken) == 0 {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	response := map[string]string{ //have to check this
		"access_token": accessToken,
	}
	json.NewEncoder(w).Encode(response)
}
