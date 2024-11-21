package tests

import (
	"infinirewards/infinirewards"
	"os"
	"testing"
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
	os.Setenv("NETWORK", "devnet")
	os.Setenv("RPC_PROVIDER_URL_DEVNET", "http://127.0.0.1:5050")
	os.Setenv("MASTER_PRIVATE_KEY_DEVNET", "0x0000000000000000000000000000000071d7bb07b9a64f6f78ac4c816aff4da9")
	os.Setenv("MASTER_ACCOUNT_ADDRESS_DEVNET", "0x64b48806902a367c8598f4f95c305e8c1a1acba5f082d294a43793113115691")
	os.Setenv("INFINI_REWARDS_FACTORY_ADDRESS_DEVNET", "0x7b7da57f1b8cd14ed14b7d7aa68fea3b38d2568936dc92fc70c553f8b5a8ce5")
}
