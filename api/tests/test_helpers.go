package tests

import (
	"fmt"
	"infinirewards/infinirewards"
	"infinirewards/jwt"
	"infinirewards/logs"
	"infinirewards/nats"
	"infinirewards/routes"
	"infinirewards/utils"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

var (
	testLogger *log.Logger
	natsServer *natsserver.Server
)

type StarkNetAccount struct {
	Address    string
	PrivateKey string
}

// Router setup
func setupTestRouter() *http.ServeMux {
	// Create new router
	mux := http.NewServeMux()

	// Set up all routes
	routes.SetAuthRoutes(mux)
	routes.SetUserRoutes(mux)
	routes.SetMerchantRoutes(mux)
	routes.SetInfiniRewardsRoutes(mux)

	return mux
}

func initializeServices() error {
	// Set required environment variables
	os.Setenv("WHATSAPP_TOKEN", "test_whatsapp_token")
	os.Setenv("WHATSAPP_PHONE_NUMBER_ID", "test_phone_number_id")
	os.Setenv("MACROKIOSK_USERNAME", "test_username")
	os.Setenv("MACROKIOSK_PASSWORD", "test_password")
	os.Setenv("MACROKIOSK_SENDER_ID", "test_sender")
	os.Setenv("JWT_SECRET", "test_jwt_secret")

	// Initialize services
	if err := utils.InitWhatsApp(); err != nil {
		return fmt.Errorf("failed to initialize WhatsApp: %v", err)
	}

	if err := utils.InitMacroKiosk(); err != nil {
		return fmt.Errorf("failed to initialize MacroKiosk: %v", err)
	}

	return nil
}

// Add this helper function
func addAuthHeader(req *http.Request, token string) {
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// setupTest sets up the test environment
func setupTest(t *testing.T) *http.ServeMux {
	// Initialize test logger first
	testLogger = log.New(os.Stdout, "[TEST] ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	testLogger.Printf("Setting up test environment for: %s", t.Name())

	logs.InitHandler("")

	// Initialize logs handler
	if err := setupLogsHandler(); err != nil {
		t.Fatalf("Failed to setup logs handler: %v", err)
	}
	testLogger.Printf("✅ Logs handler initialized")

	// Initialize JWT signing keys
	if err := setupTestJWTKeys(); err != nil {
		t.Fatalf("Failed to setup JWT keys: %v", err)
	}
	testLogger.Printf("✅ JWT keys initialized")

	// Initialize required services
	if err := initializeServices(); err != nil {
		testLogger.Printf("❌ Service initialization failed: %v", err)
		t.Fatalf("Failed to initialize services: %v", err)
	}
	testLogger.Printf("✅ Services initialized")

	// Setup NATS
	if err := setupNATSServer(); err != nil {
		testLogger.Printf("❌ NATS setup failed: %v", err)
		t.Fatalf("Failed to setup NATS: %v", err)
	}
	testLogger.Printf("✅ NATS setup completed")

	// Setup StarkNet
	setupTestStarkNet(t)
	if infinirewards.Client != nil {
		testLogger.Printf("✅ StarkNet client initialized")
	} else {
		testLogger.Printf("⚠️ StarkNet client not initialized")
	}

	// Create and return router
	router := setupTestRouter()
	if router == nil {
		testLogger.Printf("❌ Router creation failed")
		t.Fatal("Failed to create test router")
	}
	testLogger.Printf("✅ Router created")

	return router
}

// Add this function to setup JWT keys for testing
func setupTestJWTKeys() error {
	// Create test directories
	testPrivateKeyDir := ".private_test"
	testPublicKeyDir := ".public_test"

	// Clean up any existing test directories
	os.RemoveAll(testPrivateKeyDir)
	os.RemoveAll(testPublicKeyDir)

	// Set test directories for JWT package
	jwt.SetKeyDirectories(testPrivateKeyDir, testPublicKeyDir)

	// Initialize JWT keys
	if err := jwt.InitKeys(); err != nil {
		return fmt.Errorf("failed to initialize JWT keys: %w", err)
	}

	return nil
}

// Add this function to setup logs handler
func setupLogsHandler() error {
	// Set up slog handler with test configuration
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	// Create a test log file
	logFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		return fmt.Errorf("failed to create test log file: %w", err)
	}

	// Initialize the handler with both file and console output
	handler := slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Initialize the logs package
	logs.InitHandler("")

	return nil
}

// Update cleanup function to handle log files
func cleanup() error {
	testLogger.Println("Starting cleanup process")

	// Clean up NATS
	if nats.NC != nil {
		testLogger.Println("Cleaning up NATS KV buckets")
		if err := nats.CleanupKVBuckets(); err != nil {
			testLogger.Printf("Error cleaning up KV buckets: %v", err)
		}

		testLogger.Println("Closing NATS connection")
		nats.Close()
	}

	if natsServer != nil {
		testLogger.Println("Shutting down NATS server")
		natsServer.Shutdown()
		natsServer.WaitForShutdown()
	}

	// Clean up JWT test directories
	testLogger.Println("Cleaning up JWT test directories")
	os.RemoveAll(".private_test")
	os.RemoveAll(".public_test")

	// Clean up test log files
	testLogger.Println("Cleaning up test log files")
	files, err := filepath.Glob("test-*.log")
	if err == nil {
		for _, f := range files {
			os.Remove(f)
		}
	}

	testLogger.Println("Cleanup completed")
	return nil
}
