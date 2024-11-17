package controllers

import (
	"encoding/json"
	"fmt"
	"infinirewards/infinirewards"
	"infinirewards/logs"
	"infinirewards/middleware"
	"infinirewards/models"
	"infinirewards/nats"
	"net/http"
	"strings"
	"time"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/oklog/ulid/v2"
)

// UserGetUserHandler godoc
// @Summary Get user details
// @Description Get detailed information about a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID" format(ulid) example(user:01HNA...)
// @Success 200 {object} models.User "User details retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 401 {object} models.ErrorResponse "Authentication error"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "id": "user:01HNA...",
//	  "phoneNumber": "+60123456789",
//	  "email": "user@example.com",
//	  "name": "John Doe",
//	  "avatar": "https://example.com/avatar.jpg",
//	  "accountAddress": "0x1234...",
//	  "createdAt": "2024-01-01T00:00:00Z",
//	  "updatedAt": "2024-01-01T00:00:00Z"
//	}
//
// @Example {json} Error Response (Invalid ID):
//
//	{
//	  "message": "Invalid user ID format",
//	  "code": "INVALID_REQUEST",
//	  "details": {
//	    "field": "id",
//	    "reason": "must be a valid ULID starting with 'user:'"
//	  }
//	}
//
// @Example {json} Error Response (Not Found):
//
//	{
//	  "message": "User not found",
//	  "code": "NOT_FOUND",
//	  "details": {
//	    "id": "user:01HNA..."
//	  }
//	}
//
// @Example {json} Error Response (Unauthorized):
//
//	{
//	  "message": "Authentication failed",
//	  "code": "UNAUTHORIZED",
//	  "details": {
//	    "reason": "invalid or expired token"
//	  }
//	}
//
// @Example {json} Error Response (Server Error):
//
//	{
//	  "message": "Internal server error",
//	  "code": "INTERNAL_ERROR",
//	  "details": {
//	    "reason": "database connection error"
//	  }
//	}
//
// @Router /users/{id} [get]
func UserGetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userGetUserHandler called", "method", r.Method)

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	userID = strings.TrimSuffix(userID, "/")

	user := &models.User{}
	err := user.GetUser(ctx, userID)
	if err != nil {
		logs.Logger.Error("userGetUserHandler Failed to get user", "error", err, "userId", userID)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Remove sensitive information before sending the response
	user.PrivateKey = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UserCreateUserHandler godoc
// @Summary Create new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User creation request"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 409 {object} models.ErrorResponse "User already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "phoneNumber": "+60123456789",
//	  "email": "user@example.com",
//	  "name": "John Doe",
//	  "avatar": "https://example.com/avatar.jpg"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "id": "user:01HNA...",
//	  "phoneNumber": "+60123456789",
//	  "email": "user@example.com",
//	  "name": "John Doe",
//	  "avatar": "https://example.com/avatar.jpg",
//	  "accountAddress": "0x1234...",
//	  "createdAt": "2024-01-01T00:00:00Z",
//	  "updatedAt": "2024-01-01T00:00:00Z"
//	}
//
// @Example {json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "phoneNumber": "must be in E.164 format",
//	    "email": "must be a valid email address"
//	  }
//	}
//
// @Example {json} Error Response (User Exists):
//
//	{
//	  "message": "User already exists",
//	  "code": "CONFLICT",
//	  "details": {
//	    "phoneNumber": "+60123456789"
//	  }
//	}
//
// @Example {json} Error Response (Server Error):
//
//	{
//	  "message": "Failed to create user",
//	  "code": "INTERNAL_ERROR",
//	  "details": {
//	    "reason": "blockchain transaction failed"
//	  }
//	}
//
// @Router /users [post]
func UserCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userCreateUserHandler called", "method", r.Method)

	var createUserRequest models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := createUserRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a new private key using starknet-go
	_, merchantPubKey, merchantPrivKey := account.GetRandomKeys()
	_, addr, err := infinirewards.CreateUser(merchantPubKey.String(), createUserRequest.PhoneNumber)
	if err != nil {
		logs.Logger.Error("userCreateUserHandler Failed to generate keys", "error", err)
		http.Error(w, "Failed to generate keys", http.StatusInternalServerError)
		return
	}

	userPublicKey, userSeed := nats.GenerateUserKey()

	user := &models.User{
		ID:             "user:" + ulid.Make().String(),
		PhoneNumber:    createUserRequest.PhoneNumber,
		Email:          createUserRequest.Email,
		Name:           createUserRequest.Name,
		Avatar:         createUserRequest.Avatar,
		PrivateKey:     merchantPrivKey.String(),
		AccountAddress: addr,
		NatsPublicKey:  userPublicKey,
		NKey:           string(userSeed),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = user.CreateUser(ctx)
	if err != nil {
		logs.Logger.Error("userCreateUserHandler Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Remove sensitive information before sending the response
	user.PrivateKey = ""
	user.NKey = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UserUpdateUserHandler godoc
// @Summary Update user details
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "User update request"
// @Success 200 {object} models.User "Updated user details"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "email": "updated@example.com",
//	  "name": "John Updated",
//	  "avatar": "https://example.com/new-avatar.jpg"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "id": "user:01HNA...",
//	  "phoneNumber": "+60123456789",
//	  "email": "updated@example.com",
//	  "name": "John Updated",
//	  "avatar": "https://example.com/new-avatar.jpg",
//	  "accountAddress": "0x1234...",
//	  "createdAt": "2024-01-01T00:00:00Z",
//	  "updatedAt": "2024-01-01T00:00:00Z"
//	}
//
// @Router /users/{id} [put]
func UserUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userUpdateUserHandler called", "method", r.Method)

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	userID = strings.TrimSuffix(userID, "/")

	var updateUserRequest models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := updateUserRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify user has permission to update this user
	authUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil || authUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Update user fields
	user.Name = updateUserRequest.Name
	user.Email = updateUserRequest.Email
	user.Avatar = updateUserRequest.Avatar
	user.UpdatedAt = time.Now()

	if err := user.UpdateUser(ctx); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UserDeleteUserHandler godoc
// @Summary Delete user
// @Description Delete an existing user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.MessageResponse "User deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "message": "User deleted successfully"
//	}
//
// @Router /users/{id} [delete]
func UserDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userDeleteUserHandler called", "method", r.Method)

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	userID = strings.TrimSuffix(userID, "/")

	// Verify user has permission to delete this user
	authUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil || authUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if err := user.DeleteUser(ctx); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UserCreateAPIKeyHandler godoc
// @Summary Create API key
// @Description Create a new API key for a user
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body models.CreateAPIKeyRequest true "API key creation request"
// @Success 201 {object} models.APIKey "Created API key details"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "name": "My API Key"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "id": "key_01HNA...",
//	  "name": "My API Key",
//	  "key": "sk_01HNA...",
//	  "createdAt": "2024-01-01T00:00:00Z"
//	}
//
// @Router /users/{id}/api-keys [post]
func UserCreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userCreateAPIKeyHandler called", "method", r.Method)

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	userID = strings.Split(userID, "/")[0]

	// Verify user has permission
	authUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil || authUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var createRequest models.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := createRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiKey, err := models.CreateAPIKey(ctx, userID, createRequest.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create API key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)
}

// UserListAPIKeysHandler godoc
// @Summary List API keys
// @Description List all API keys for a user
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {array} models.APIKey "List of API keys"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	[
//	  {
//	    "id": "key_01HNA...",
//	    "name": "My API Key",
//	    "createdAt": "2024-01-01T00:00:00Z"
//	  }
//	]
//
// @Router /users/{id}/api-keys [get]
func UserListAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userListAPIKeysHandler called", "method", r.Method)

	// Get user ID from URL path
	userID := strings.TrimPrefix(r.URL.Path, "/users/")
	userID = strings.Split(userID, "/")[0]

	// Verify user has permission
	authUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil || authUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	apiKeys, err := models.ListAPIKeys(ctx, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list API keys: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// UserDeleteAPIKeyHandler godoc
// @Summary Delete API key
// @Description Delete an API key
// @Tags api-keys
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param keyId path string true "API Key ID"
// @Success 200 {object} models.MessageResponse "API key deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized access"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "message": "API key deleted successfully"
//	}
//
// @Router /users/{id}/api-keys/{keyId} [delete]
func UserDeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userDeleteAPIKeyHandler called", "method", r.Method)

	// Get user ID and key ID from URL path
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/users/"), "/")
	if len(parts) != 4 { // Expected: [userID, "api-keys", keyID, ""]
		http.Error(w, "invalid URL format", http.StatusBadRequest)
		return
	}
	userID := parts[0]
	keyID := parts[2]

	// Verify user has permission
	authUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil || authUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := models.DeleteAPIKey(ctx, userID, keyID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete API key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MessageResponse{
		Message: "API key deleted successfully",
	})
}
