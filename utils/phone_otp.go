package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var debugPhoneNumbers = []string{
	"+60123456789",
	"+60129876543",
	"+60143382537",
}

// SendPhoneOTP sends OTP to the given phone number
func SendPhoneOTP(phoneNumber string, otp string) error {
	// For test phone numbers (+TEST...), just log and do nothing
	if strings.HasPrefix(phoneNumber, "+TEST") {
		slog.Info("Test OTP generated",
			slog.String("phone", phoneNumber),
			slog.String("otp", otp),
		)
		return nil
	}

	message := fmt.Sprintf("Your InfiniRewards verification code is: %s", otp)

	// Check if it's a debug phone number
	for _, debugNumber := range debugPhoneNumbers {
		if phoneNumber == debugNumber {
			return SendDiscordMessage(message)
		}
	}

	// Check if user is WhatsApp user
	isWhatsApp, err := IsWhatsAppUser(phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to check WhatsApp status: %w", err)
	}

	if isWhatsApp {
		return SendWhatsAppOTP(phoneNumber, message)
	}

	// For regular phone numbers, send via MacroKiosk SMS
	return SendSMS(phoneNumber, message)
}

// SendDiscordMessage sends a message to Discord webhook
func SendDiscordMessage(message string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL not set")
	}

	payload := map[string]string{
		"content": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send discord message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord webhook returned status: %d", resp.StatusCode)
	}

	return nil
}
