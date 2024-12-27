package authy

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/resend/resend-go/v2"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// OTPRecord stores the OTP information
type OTPRecord struct {
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// In-memory store for OTPs (in production, use Redis or similar)
var otpStore = make(map[string]OTPRecord)

// GenerateOTP creates a 4-digit OTP
func GenerateOTP() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	otp := ""
	for i := 0; i < 6; i++ {
		otp += fmt.Sprintf("%d", int(bytes[i])%10)
	}
	return otp
}

// AuthenticateViaEmail generates and sends OTP via email
func AuthenticateViaEmail(email string) error {
	// Generate 4-digit OTP
	otp := GenerateOTP()

	// Store OTP with expiration (15 minutes)
	otpStore[email] = OTPRecord{
		Email:     email,
		OTP:       otp,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Initialize Resend client
	apiKey := "re_3Qm4iMF8_5Y5iuFJN8iosJKimcMTdUBNa"

	client := resend.NewClient(apiKey)

	// Create email HTML content
	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>Your Verification Code</h2>
			<p>Your OTP for Krishi Bazar app is: <strong>%s</strong></p>
			<p>This code will expire in 15 minutes.</p>
		</div>
	`, otp)

	params := &resend.SendEmailRequest{
		From:    "Krishi Bazar <onboarding@resend.dev>",
		To:      []string{email},
		Html:    htmlContent,
		Subject: "Your Verification Code - Krishi Bazar",
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// VerifyOTP checks if the provided OTP is valid
func VerifyOTP(email, providedOTP string) error {
	record, exists := otpStore[email]
	if !exists {
		return fmt.Errorf("no OTP found for this email")
	}

	if time.Now().After(record.ExpiresAt) {
		delete(otpStore, email)
		return fmt.Errorf("OTP has expired")
	}

	if record.OTP != providedOTP {
		return fmt.Errorf("invalid OTP")
	}

	// Clean up used OTP
	delete(otpStore, email)
	return nil
}











//  AUTH VIA PHONE NUMBER -> TWILLIO

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