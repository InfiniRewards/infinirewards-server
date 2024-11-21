package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"infinirewards/models"
	"infinirewards/nats"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

// TestUser represents a test user with their auth token
type TestUser struct {
	User  *models.User
	Token *models.Token
}

// createTestUserWithAuth creates a test user following the proper auth flow
func createTestUserWithAuth(t *testing.T, router *http.ServeMux) *TestUser {
	testLogger.Printf("Starting test user creation flow")

	// Generate test phone number
	phoneNumber := generateTestPhoneNumber()
	testLogger.Printf("Using test phone number: %s", phoneNumber)

	// Step 1: Request OTP
	token := requestOTPAndAuthenticate(t, router, phoneNumber)
	testLogger.Printf("Got initial auth token")

	// Step 2: Create user (deploy account onchain)
	createReq := models.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	reqBody, _ := json.Marshal(createReq)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, token.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create user. Response: %s", w.Body.String())

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err, "Failed to unmarshal user response")
	assert.NotEmpty(t, user.ID, "User ID should not be empty")

	testLogger.Printf("Created user with ID: %s", user.ID)
	return &TestUser{
		User:  &user,
		Token: token,
	}
}

// requestOTPAndAuthenticate handles the OTP request and authentication flow
func requestOTPAndAuthenticate(t *testing.T, router *http.ServeMux, phoneNumber string) *models.Token {
	// Request OTP
	otpReq := models.RequestOTPRequest{
		PhoneNumber: phoneNumber,
	}
	reqBody, _ := json.Marshal(otpReq)

	req := httptest.NewRequest("POST", "/auth/request-otp", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Failed to request OTP. Response: %s", w.Body.String())

	var otpResp models.RequestOTPResponse
	err := json.Unmarshal(w.Body.Bytes(), &otpResp)
	assert.NoError(t, err, "Failed to unmarshal OTP response")
	assert.NotEmpty(t, otpResp.ID, "OTP ID should not be empty")

	// Get verification from NATS
	entry, err := nats.GetKV(context.Background(), "phoneVerification", otpResp.ID)
	assert.NoError(t, err, "Failed to get verification from NATS")

	var verification models.PhoneNumberVerification
	err = json.Unmarshal(entry.Value(), &verification)
	assert.NoError(t, err, "Failed to unmarshal verification")

	// Generate valid OTP using the secret
	otp, err := totp.GenerateCodeCustom(verification.Secret, time.Now(), totp.ValidateOpts{
		Period:    300,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	assert.NoError(t, err, "Failed to generate OTP")

	// Authenticate with valid OTP
	authReq := models.AuthenticateRequest{
		ID:        otpResp.ID,
		Method:    "otp",
		Token:     otp,
		Signature: "test_device_signature",
	}
	reqBody, _ = json.Marshal(authReq)

	req = httptest.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Failed to authenticate. Response: %s", w.Body.String())

	var authResp models.AuthenticateResponse
	err = json.Unmarshal(w.Body.Bytes(), &authResp)
	assert.NoError(t, err, "Failed to unmarshal auth response")
	assert.NotEmpty(t, authResp.Token, "Auth token should not be empty")

	return &authResp.Token
}

// createTestMerchantWithAuth creates a test merchant following the proper auth flow
func createTestMerchantWithAuth(t *testing.T, router *http.ServeMux) *TestUser {
	testLogger.Printf("Creating test merchant")

	// First create a user with auth
	testUser := createTestUserWithAuth(t, router)

	// Create merchant contract
	merchantReq := models.CreateMerchantRequest{
		Name:     testUser.User.Name,
		Symbol:   "TST",
		Decimals: 18,
	}
	reqBody, _ := json.Marshal(merchantReq)

	req := httptest.NewRequest("POST", "/merchant", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, testUser.Token.AccessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create merchant. Response: %s", w.Body.String())

	testLogger.Printf("Created merchant for user ID: %s", testUser.User.ID)
	return testUser
}

// generateTestPhoneNumber generates a phone number for testing
func generateTestPhoneNumber() string {
	return fmt.Sprintf("+TEST%07d", rand.Intn(10000000))
}
