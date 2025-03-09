package util

import (
	"context"
	"time"

	"github.com/loginOAuth/database"
	"github.com/loginOAuth/logger"
	"github.com/sirupsen/logrus"
)

func StoreOTP(phoneNumber string, otp string) error {
	if err := database.Rdb.Set(context.Background(), phoneNumber, otp, 5*time.Minute).Err(); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"redis.storeotp.error": err.Error(),
			"otp":                  otp,
		}).Error("Otp not stored for phone number ", phoneNumber)
		return err
	} else {
		logger.Logger.WithFields(logrus.Fields{
			"otp": otp,
		}).Info("Otp stored for phone number ", phoneNumber)
		return nil
	}
}

func FetchOTP(phoneNumber string) string {
	otp, err := database.Rdb.Get(context.Background(), phoneNumber).Result()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"redis.fetchotp.error": err.Error(),
		}).Error("Otp not fetched for phone number ", phoneNumber)
	} else {
		logger.Logger.WithFields(logrus.Fields{
			"otp": otp,
		}).Info("Otp fetched for phone number ", phoneNumber)
	}
	return otp
}
