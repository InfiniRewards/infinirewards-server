package controllers

import (
	"encoding/json"
	"fmt"
	"infinirewards/infinirewards"
	"infinirewards/logs"
	"infinirewards/middleware"
	"infinirewards/models"
	"net/http"
	"strings"
	"time"

	"github.com/NethermindEth/starknet.go/account"
)

// UserGetUserHandler godoc
//	@Summary		Get user details
//	@Description	Get authenticated user details
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.User				"User details retrieved successfully"
//	@Failure		401	{object}	models.ErrorResponse	"Authentication error"
//	@Failure		404	{object}	models.ErrorResponse	"User not found"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Example		{json} Success Response:
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
//	@Router			/user [get]
func UserGetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userGetUserHandler called", "method", r.Method)

	// Get user ID from context instead of URL
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	err = user.GetUser(ctx, userID)
	if err != nil {
		logs.Logger.Error("userGetUserHandler Failed to get user", "error", err, "userId", userID)
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	// Remove sensitive information before sending the response
	user.PrivateKey = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UserCreateUserHandler godoc
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateUserRequest	true	"User Creation Request"
//	@Success		201		{object}	models.User					"User created successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format"
//	@Failure		409		{object}	models.ErrorResponse		"User already exists"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "phoneNumber": "+60123456789",
//	  "email": "user@example.com",
//	  "name": "John Doe",
//	  "avatar": "https://example.com/avatar.jpg"
//	}
//
//	@Example		{json} Success Response:
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
//	@Example		{json} Error Response (Invalid Request):
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
//	@Router			/user [post]
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

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	if user.PrivateKey != "" {
		http.Error(w, "user already created", http.StatusConflict)
		return
	}

	// Generate a new private key using starknet-go
	_, publicKey, privateKey := account.GetRandomKeys()
	_, addr, err := infinirewards.CreateUser(publicKey.String(), user.PhoneNumber)
	if err != nil {
		logs.Logger.Error("userCreateUserHandler Failed to generate keys", "error", err)
		http.Error(w, "Failed to generate keys", http.StatusInternalServerError)
		return
	}

	_, err = infinirewards.FundAccount(addr)
	if err != nil {
		logs.Logger.Error("userCreateUserHandler Failed to fund account", "error", err)
		http.Error(w, "Failed to fund account", http.StatusInternalServerError)
		return
	}

	user.Name = createUserRequest.Name
	user.Email = createUserRequest.Email
	user.Avatar = createUserRequest.Avatar

	user.PrivateKey = privateKey.String()
	user.PublicKey = publicKey.String()
	user.AccountAddress = addr
	user.UpdatedAt = time.Now()

	err = user.UpdateUser(ctx)
	if err != nil {
		logs.Logger.Error("userCreateUserHandler Failed to create user", "error", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Remove sensitive information before sending the response
	user.PrivateKey = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UserUpdateUserHandler godoc
//	@Summary		Update user
//	@Description	Update authenticated user details
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.UpdateUserRequest	true	"User Update Request"
//	@Success		200		{object}	models.User					"User updated successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format"
//	@Failure		401		{object}	models.ErrorResponse		"Unauthorized access"
//	@Failure		404		{object}	models.ErrorResponse		"User not found"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "email": "updated@example.com",
//	  "name": "John Updated",
//	  "avatar": "https://example.com/new-avatar.jpg"
//	}
//
//	@Example		{json} Success Response:
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
//	@Router			/user [put]
func UserUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userUpdateUserHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var updateUserRequest models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := updateUserRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	// Remove sensitive information before sending response
	user.PrivateKey = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UserDeleteUserHandler godoc
//	@Summary		Delete user
//	@Description	Delete authenticated user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.MessageResponse	"User deleted successfully"
//	@Failure		401	{object}	models.ErrorResponse	"Unauthorized access"
//	@Failure		404	{object}	models.ErrorResponse	"User not found"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Router			/user [delete]
func UserDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userDeleteUserHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	if err := user.DeleteUser(ctx); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MessageResponse{
		Message: "User deleted successfully",
	})
}

// UserCreateAPIKeyHandler godoc
//	@Summary		Create API key
//	@Description	Create a new API key for authenticated user
//	@Tags			api-keys
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreateAPIKeyRequest	true	"API Key Creation Request"
//	@Success		201		{object}	models.APIKey				"API key created successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format"
//	@Failure		401		{object}	models.ErrorResponse		"Unauthorized access"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "name": "My API Key"
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "id": "key_01HNA...",
//	  "name": "My API Key",
//	  "key": "sk_01HNA...",
//	  "createdAt": "2024-01-01T00:00:00Z"
//	}
//
//	@Router			/user/api-keys [post]
func UserCreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userCreateAPIKeyHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
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
//	@Summary		List API keys
//	@Description	List all API keys for authenticated user
//	@Tags			api-keys
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.APIKey			"List of API keys"
//	@Failure		401	{object}	models.ErrorResponse	"Unauthorized access"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Example		{json} Success Response:
//
//	[
//	  {
//	    "id": "key_01HNA...",
//	    "name": "My API Key",
//	    "createdAt": "2024-01-01T00:00:00Z"
//	  }
//	]
//
//	@Router			/user/api-keys [get]
func UserListAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userListAPIKeysHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
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
//	@Summary		Delete API key
//	@Description	Delete an API key for authenticated user
//	@Tags			api-keys
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			keyId	path		string					true	"API Key ID"
//	@Success		200		{object}	models.MessageResponse	"API key deleted successfully"
//	@Failure		400		{object}	models.ErrorResponse	"Invalid request format"
//	@Failure		401		{object}	models.ErrorResponse	"Unauthorized access"
//	@Failure		500		{object}	models.ErrorResponse	"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "message": "API key deleted successfully"
//	}
//
//	@Example		{json} Error Response (Invalid Key):
//
//	{
//	  "message": "Invalid API key",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "reason": "API key not found or does not belong to user"
//	  }
//	}
//
//	@Router			/user/api-keys/{keyId} [delete]
func UserDeleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("userDeleteAPIKeyHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract keyId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid URL format", http.StatusBadRequest)
		return
	}
	keyID := parts[len(parts)-1]

	// Verify the API key belongs to the user
	keyUserID := strings.Split(keyID, ".")[0]
	if keyUserID != userID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := models.DeleteAPIKey(ctx, keyID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete API key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MessageResponse{
		Message: "API key deleted successfully",
	})
}
