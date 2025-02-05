package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

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
		return err.Errorf("Failed to send OTP: %v", err)
	}

	return nil
}

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := r.Intn(10000)
	return strconv.Itoa(otp)
}
