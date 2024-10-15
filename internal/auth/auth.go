package authy

import (
	"fmt"
	"os"
	"strings"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

func formatPhoneNumber(number string) string {
	// Remove any non-digit characters
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	number = strings.ReplaceAll(number, "(", "")
	number = strings.ReplaceAll(number, ")", "")

	// If the number doesn't start with '+', assume it's an Indian number and add +91
	if !strings.HasPrefix(number, "+") {
		number = "+91" + number
	}

	return number
}

func Authenticate(number string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	verifySid := os.Getenv("TWILIO_VERIFY_SID")

	// Format the phone number
	formattedNumber := formatPhoneNumber(number)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &verify.CreateVerificationParams{}
	params.SetTo(formattedNumber)
	params.SetChannel("sms")

	res, err := client.VerifyV2.CreateVerification(verifySid, params)
	if err != nil {
		fmt.Printf("Error sending verification: %v\n", err)
		return fmt.Errorf("error sending verification code: %v", err)
	}

	if res.Status != nil {
		fmt.Println("Verification status:", *res.Status)
	}

	return nil
}

func VerifyCode(number string, code string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	verifySid := os.Getenv("TWILIO_VERIFY_SID")

	// Format the phone number
	formattedNumber := formatPhoneNumber(number)

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(formattedNumber)
	params.SetCode(code)

	res, err := client.VerifyV2.CreateVerificationCheck(verifySid, params)
	if err != nil {
		fmt.Printf("Error checking verification: %v\n", err)
		return fmt.Errorf("error checking verification code: %v", err)
	}

	if res.Status != nil && *res.Status == "approved" {
		fmt.Println("Verification successful")
		return nil
	} else {
		fmt.Println("Verification failed")
		return fmt.Errorf("verification failed")
	}
}