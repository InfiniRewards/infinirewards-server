package tests

import (
	"bytes"
	"encoding/json"
	"infinirewards/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserManagement(t *testing.T) {
	router := setupTest(t)

	t.Run("Update User", func(t *testing.T) {
		testUser := createTestUserWithAuth(t, router)

		updateReq := models.UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		reqBody, _ := json.Marshal(updateReq)

		req := httptest.NewRequest("PUT", "/user", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testUser.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updatedUser models.User
		err := json.Unmarshal(w.Body.Bytes(), &updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, updateReq.Name, updatedUser.Name)
		assert.Equal(t, updateReq.Email, updatedUser.Email)
	})

	t.Run("API Key Management", func(t *testing.T) {
		testUser := createTestUserWithAuth(t, router)

		// Create API Key
		createKeyReq := models.CreateAPIKeyRequest{
			Name: "Test API Key",
		}
		reqBody, _ := json.Marshal(createKeyReq)

		req := httptest.NewRequest("POST", "/user/api-keys", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		addAuthHeader(req, testUser.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var apiKey models.APIKey
		err := json.Unmarshal(w.Body.Bytes(), &apiKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, apiKey.ID)
		assert.NotEmpty(t, apiKey.Secret)

		// List API Keys
		req = httptest.NewRequest("GET", "/user/api-keys", nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var apiKeys []*models.APIKey
		err = json.Unmarshal(w.Body.Bytes(), &apiKeys)
		assert.NoError(t, err)
		assert.NotEmpty(t, apiKeys)

		// Delete API Key
		req = httptest.NewRequest("DELETE", "/user/api-keys/"+apiKey.ID, nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify API key was deleted by listing again
		req = httptest.NewRequest("GET", "/user/api-keys", nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &apiKeys)
		assert.NoError(t, err)
		assert.Empty(t, apiKeys)
	})

	t.Run("Delete User", func(t *testing.T) {
		testUser := createTestUserWithAuth(t, router)

		req := httptest.NewRequest("DELETE", "/user", nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify user was deleted by trying to get user details
		req = httptest.NewRequest("GET", "/user", nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get User Details", func(t *testing.T) {
		testUser := createTestUserWithAuth(t, router)

		req := httptest.NewRequest("GET", "/user", nil)
		addAuthHeader(req, testUser.Token.AccessToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var user models.User
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, testUser.User.ID, user.ID)
		assert.Equal(t, testUser.User.PhoneNumber, user.PhoneNumber)
		assert.Empty(t, user.PrivateKey, "Private key should not be returned")
	})
}
