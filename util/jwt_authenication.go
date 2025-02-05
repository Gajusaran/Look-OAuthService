package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PhoneNumber string `json:"phn"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(phonenumber string) (string, error) {
	claims := &Claims{
		PhoneNumber: phonenumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("could not create access token: %v", err)
	}

	return signedToken, nil
}

func GenerateRefreshToken(phonenumber string) (string, error) {
	claims := &Claims{
		PhoneNumber: phonenumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token expires in 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the refresh token with the secret key
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("could not create refresh token: %v", err)
	}

	return signedToken, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC and return the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and return claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateToken(phoneNumber string) (map[string]string, error) {
	// JWT logic , generate jwt and send
	accessToken, err := GenerateAccessToken(phoneNumber)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(phoneNumber)
	if err != nil {
		return nil, err
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return response, nil
}
