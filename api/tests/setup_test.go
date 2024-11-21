package tests

import (
	"context"
	"infinirewards/infinirewards"
	"infinirewards/logs"
	"infinirewards/nats"
	"infinirewards/utils"
	"log/slog"
	"os"
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func TestSetup(t *testing.T) {
	t.Log("Starting setup test")

	// Set required environment variables
	os.Setenv("WHATSAPP_TOKEN", "test_whatsapp_token")
	os.Setenv("WHATSAPP_PHONE_NUMBER_ID", "test_phone_number_id")
	os.Setenv("MACROKIOSK_USERNAME", "test_username")
	os.Setenv("MACROKIOSK_PASSWORD", "test_password")
	os.Setenv("MACROKIOSK_SENDER_ID", "test_sender")
	os.Setenv("JWT_SECRET", "test_jwt_secret")

	// Initialize logger
	logs.InitHandler("")

	// Initialize WhatsApp
	if err := utils.InitWhatsApp(); err != nil {
		logs.Logger.Error("failed to initialize WhatsApp",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		t.Fatalf("failed to initialize WhatsApp: %v", err)
	}

	// Initialize MacroKiosk
	if err := utils.InitMacroKiosk(); err != nil {
		logs.Logger.Error("failed to initialize MacroKiosk",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		t.Fatalf("failed to initialize MacroKiosk: %v", err)
	}

	// Create router
	router := setupTestRouter()
	if router == nil {
		t.Fatal("Failed to create test router")
	}
	t.Log("✅ Router created successfully")

	// Test NATS connection
	if nats.NC == nil {
		t.Fatal("NATS connection not initialized")
	}

	// Test NATS connection is alive
	if !nats.NC.IsConnected() {
		t.Fatal("NATS connection is not active")
	}
	t.Log("✅ NATS connection verified")

	// Test JetStream is available
	js, err := jetstream.New(nats.NC)
	if err != nil {
		t.Fatalf("Failed to create JetStream context: %v", err)
	}

	// Test KV store
	_, err = js.KeyValue(context.Background(), "users")
	if err != nil {
		t.Fatalf("Failed to access users KV bucket: %v", err)
	}
	t.Log("✅ JetStream KV store verified")

	// Initialize StarkNet
	setupTestStarkNet(t)
	if infinirewards.Client != nil {
		t.Log("✅ StarkNet client initialized")
	} else {
		t.Log("⚠️ StarkNet client not initialized - some tests may be skipped")
	}

	t.Log("All setup tests passed successfully")
}

// TestMain is already defined in test_helpers.go
