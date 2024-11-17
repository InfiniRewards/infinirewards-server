package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	// Test phone numbers that will receive OTP via Discord webhook
	testPhoneNumbers = map[string]bool{
		"+60143382537": true,
		"+60123456789": true, // Add more test numbers as needed
	}
	// Discord webhook URL for test OTPs
	discordWebhookURL string
	isDebugMode       bool
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discordWebhookURL = os.Getenv("MOCK_OTP_DISCORD_URL")
	isDebugMode = os.Getenv("DEBUG_MODE") == "true"
}

// SendPhoneOTP handles all OTP sending logic with proper fallbacks
func SendPhoneOTP(phoneNumber string, code string) error {
	// In debug mode or for test numbers, send to Discord
	if isDebugMode || testPhoneNumbers[phoneNumber] {
		return sendTestOTPToDiscord(phoneNumber, code)
	}

	// Production flow: Try WhatsApp first, fallback to SMS
	isWhatsAppUser, err := IsWhatsAppUser(phoneNumber)
	if err != nil {
		log.Printf("Error checking WhatsApp status: %v", err)
		// Fallback to SMS on WhatsApp check error
		return SendSMS(phoneNumber, code)
	}

	if isWhatsAppUser {
		err = SendWhatsAppOTP(phoneNumber, code)
		if err == nil {
			return nil // Successfully sent via WhatsApp
		}
		log.Printf("WhatsApp OTP failed, falling back to SMS: %v", err)
	}

	// Fallback or default to SMS
	return SendSMS(phoneNumber, code)
}

// sendTestOTPToDiscord sends OTP to Discord webhook for test phone numbers
func sendTestOTPToDiscord(phoneNumber, code string) error {
	message := struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}{
		Username: fmt.Sprintf("OTP Bot (%s)", phoneNumber),
		Content: fmt.Sprintf("```\nTest OTP for %s\nCode: %s\nValid for: 15 minutes\n```",
			phoneNumber, code),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling discord message: %w", err)
	}

	req, err := http.NewRequest("POST", discordWebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending to discord: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook error: status %d", resp.StatusCode)
	}

	return nil
}
