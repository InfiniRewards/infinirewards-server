package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"infinirewards/models"

	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	router := setupTest(t)

	t.Run("API Key Flow", func(t *testing.T) {
		testLogger.Printf("Starting API key flow test")

		// Create user with auth
		testUser := createTestUserWithAuth(t, router)
		testLogger.Printf("Created test user and got auth token")

		// Create API key
		createKeyReq := models.CreateAPIKeyRequest{
			Name: "Test API Key",
		}
		reqBody, err := json.Marshal(createKeyReq)
		assert.NoError(t, err, "Failed to marshal create API key request")

		req := httptest.NewRequest("POST", "/user/api-keys", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testUser.Token.AccessToken)

		w := httptest.NewRecorder()
		testLogger.Printf("Creating API key")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Response body: %s", w.Body.String())

		var apiKey models.APIKey
		err = json.Unmarshal(w.Body.Bytes(), &apiKey)
		assert.NoError(t, err, "Failed to unmarshal API key response")
		assert.NotEmpty(t, apiKey.ID, "API key ID should not be empty")
		assert.NotEmpty(t, apiKey.Secret, "API key secret should not be empty")
		testLogger.Printf("API key created successfully")

		// Test authenticate with API key
		authReq := models.AuthenticateRequest{
			ID:        apiKey.ID,
			Method:    "secret",
			Token:     apiKey.Secret,
			Signature: "test_device_signature",
		}
		reqBody, err = json.Marshal(authReq)
		assert.NoError(t, err, "Failed to marshal auth request")

		req = httptest.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		testLogger.Printf("Authenticating with API key")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Response body: %s", w.Body.String())

		var authResp models.AuthenticateResponse
		err = json.Unmarshal(w.Body.Bytes(), &authResp)
		assert.NoError(t, err, "Failed to unmarshal auth response")
		assert.NotEmpty(t, authResp.Token, "Auth token should not be empty")
		testLogger.Printf("API key authentication successful, %v", authResp.Token)

		// Test authenticate with invalid API key
		authReq.Token = "invalid_secret"
		reqBody, err = json.Marshal(authReq)
		assert.NoError(t, err, "Failed to marshal invalid auth request")

		req = httptest.NewRequest("POST", "/auth/authenticate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		testLogger.Printf("Testing invalid API key")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should fail with invalid API key")
		testLogger.Printf("Invalid API key test successful")

		// Test refresh token
		refreshReq := models.RefreshTokenRequest{
			RefreshToken: authResp.Token.ID,
		}
		reqBody, err = json.Marshal(refreshReq)
		assert.NoError(t, err, "Failed to marshal refresh token request")

		req = httptest.NewRequest("POST", "/auth/refresh-token", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, authResp.Token.AccessToken)

		w = httptest.NewRecorder()
		testLogger.Printf("Testing token refresh")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Token refresh should succeed, %v", w.Body.String())

		var refreshResp models.RefreshTokenResponse
		err = json.Unmarshal(w.Body.Bytes(), &refreshResp)
		assert.NoError(t, err, "Failed to unmarshal refresh response")
		assert.NotEmpty(t, refreshResp.Token.AccessToken, "New access token should not be empty")
		assert.NotEqual(t, authResp.Token.AccessToken, refreshResp.Token.AccessToken, "New token should be different")
		testLogger.Printf("Token refresh successful")

		// Verify can access protected endpoint with new token
		req = httptest.NewRequest("GET", "/user", nil)
		addAuthHeader(req, refreshResp.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should be able to access protected endpoint with new token")
		testLogger.Printf("Protected endpoint access with refreshed token successful")
	})
}
