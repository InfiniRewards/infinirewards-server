package tests

import (
	"infinirewards/infinirewards"
	"testing"

	"github.com/joho/godotenv"
)

// StarkNet setup functions
func setupTestStarkNet(t *testing.T) {
	testLogger.Println("Setting up StarkNet test configuration")

	setupStarkNetEnv()

	// Initialize StarkNet module
	if err := infinirewards.ConnectStarknet(); err != nil {
		testLogger.Printf("Failed to initialize StarkNet module: %v", err)
		t.Skipf("StarkNet initialization failed: %v. Skipping tests that require StarkNet.", err)
		return
	}
}

func setupStarkNetEnv() {
	godotenv.Load(".env")
}
