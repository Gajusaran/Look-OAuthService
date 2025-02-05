package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Gajusaran/Look-OAuthService/database"
)

func StoreOTP(phoneNumber string, otp string) {
	if err := database.Rdb.Set(context.Background(), phoneNumber, otp, 5*time.Minute).Err(); err != nil {
		log.Fatalf("failed to store OTP: %v", err)
	}
}

func FetchOTP(phoneNumber string) (string, error) {
	otp, err := database.Rdb.Get(context.Background(), phoneNumber).Result()
	if err != nil {
		return "", fmt.Errorf("failed to fetch OTP: %v", err)
	}
	return otp, nil
}
