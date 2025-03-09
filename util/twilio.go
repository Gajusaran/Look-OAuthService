package util

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/loginOAuth/logger"
	"github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOTP(phoneNumber string, otp string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(fromPhone)
	params.SetBody(fmt.Sprintf("Your OTP is: %s", otp))

	if _, err := client.Api.CreateMessage(params); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"twilio.error": err.Error(),
			"otp":          otp,
		}).Error("Error sending otp to phone number, ", phoneNumber)
		return err
	}else{
		logger.Logger.WithFields(logrus.Fields{
			"otp": otp,
		}).Info("Otp sent to phone number ",phoneNumber)
		return nil
	}

}

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := r.Intn(10000)
	return strconv.Itoa(otp)
}
