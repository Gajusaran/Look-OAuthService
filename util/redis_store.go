package util

import (
	"context"
	"fmt"
	"time"

	"github.com/Gajusaran/Look-OAuthService/database"
)

func StoreOTP(phoneNumber string, otp string) error {
	err := database.Rdb.Set(context.Background(), phoneNumber, otp, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP: %v", err)
	}
	return nil
}

func FetchOTP(phoneNumber string) (string, error) {
	otp, err := database.Rdb.Get(context.Background(), phoneNumber).Result()
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %v", err)
	}
	return otp, nil
}
