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

func SendOTP(phoneNumber string, otp string) {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_PHONE_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	message := fmt.Sprintf("Your OTP is: %s", otp)

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(fromPhone)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params) // will discuss about handle this error, res
	if err != nil {
		log.Fatalf("Failed to send OTP: %v", err)
	}

	StoreOTP(phoneNumber, otp) //error handling
}

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := r.Intn(10000)
	return strconv.Itoa(otp)
}
