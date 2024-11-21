package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"infinirewards/jwt"
	"infinirewards/logs"
	"infinirewards/middleware"
	"infinirewards/models"
	"infinirewards/nats"
	"infinirewards/utils"
	"log/slog"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// RequestOTPHandler godoc
//	@Summary		Request OTP
//	@Description	Request a one-time password for authentication
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RequestOTPRequest	true	"OTP Request"
//	@Success		200		{object}	models.RequestOTPResponse	"OTP sent successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format or validation failed"
//	@Failure		429		{object}	models.ErrorResponse		"Too many requests"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "phoneNumber": "+60123456789"  // E.164 format required
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "message": "OTP sent successfully",
//	  "id": "01HNAJ6640M9JRRJFQSZZVE3HH"
//	}
//
//	@Example		{json} Error Response (Invalid Phone):
//
//	{
//	  "message": "Invalid phone number format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "phoneNumber",
//	    "reason": "must be in E.164 format (e.g. +60123456789)"
//	  }
//	}
//
//	@Example		{json} Error Response (Rate Limit):
//
//	{
//	  "message": "Too many OTP requests",
//	  "code": "RATE_LIMIT_EXCEEDED",
//	  "details": {
//	    "retryAfter": "60s",
//	    "limit": "3 requests per hour"
//	  }
//	}
//
//	@Example		{json} Error Response (Server Error):
//
//	{
//	  "message": "Failed to send OTP",
//	  "code": "INTERNAL_ERROR",
//	  "details": {
//	    "reason": "SMS service unavailable"
//	  }
//	}
//
//	@Router			/auth/request-otp [post]
func RequestOTPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs.Logger.Info("processing OTP request",
		slog.String("handler", "RequestOTPHandler"),
		slog.String("method", r.Method),
	)

	var requestOTPRequest models.RequestOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&requestOTPRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := requestOTPRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	phoneNumber := requestOTPRequest.PhoneNumber
	user := models.User{}
	err := user.GetUserFromPhoneNumber(ctx, phoneNumber)
	if err != nil {
		// Create new user if not found
		// userPublicKey, userSeed := nats.GenerateUserKey()
		user = models.User{
			PhoneNumber: phoneNumber,
			CreatedAt:   time.Now(),
		}
		if err := user.CreateUser(ctx); err != nil {
			logs.Logger.Error("failed to create user",
				slog.String("handler", "RequestOTPHandler"),
				slog.String("phone_number", phoneNumber),
				slog.String("error", err.Error()),
			)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "InfiniRewards Auth",
		AccountName: phoneNumber,
	})
	if err != nil {
		logs.Logger.Error("failed to generate TOTP key",
			slog.String("handler", "RequestOTPHandler"),
			slog.String("error", err.Error()),
		)
		http.Error(w, "failed to generate OTP", http.StatusInternalServerError)
		return
	}

	phoneNumberVerification := models.PhoneNumberVerification{
		ID:        ulid.Make().String(),
		Secret:    key.Secret(),
		User:      user.ID,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(phoneNumberVerification)
	if err != nil {
		logs.Logger.Error("failed to marshal phone verification",
			slog.String("handler", "RequestOTPHandler"),
			slog.String("verification_id", phoneNumberVerification.ID),
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := nats.PutKV(ctx, "phoneVerification", phoneNumberVerification.ID, data); err != nil {
		logs.Logger.Error("failed to store phone verification",
			slog.String("handler", "RequestOTPHandler"),
			slog.String("error", err.Error()),
		)
		http.Error(w, "failed to store verification", http.StatusInternalServerError)
		return
	}

	passcode, err := totp.GenerateCodeCustom(key.Secret(), time.Now(), totp.ValidateOpts{
		Period:    300,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		logs.Logger.Error("failed to generate OTP code",
			slog.String("handler", "RequestOTPHandler"),
			slog.String("error", err.Error()),
		)
		http.Error(w, "failed to generate OTP", http.StatusInternalServerError)
		return
	}

	if err := utils.SendPhoneOTP(phoneNumber, passcode); err != nil {
		logs.Logger.Error("failed to send OTP",
			slog.String("handler", "RequestOTPHandler"),
			slog.String("phone_number", phoneNumber),
			slog.String("error", err.Error()),
		)
		http.Error(w, "failed to send OTP", http.StatusInternalServerError)
		return
	}

	response := models.RequestOTPResponse{
		Message: "OTP sent successfully",
		ID:      phoneNumberVerification.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AuthenticateHandler godoc
//	@Summary		Authenticate user
//	@Description	Authenticate user using OTP or API key
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.AuthenticateRequest	true	"Authentication Request"
//	@Success		200		{object}	models.AuthenticateResponse	"Authentication successful"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse		"Authentication failed"
//	@Failure		429		{object}	models.ErrorResponse		"Too many attempts"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body (OTP):
//
//	{
//	  "id": "01HNAJ6640M9JRRJFQSZZVE3HH",  // Verification ID from request-otp
//	  "token": "123456",                     // 6-digit OTP code
//	  "method": "otp",                       // Authentication method
//	  "signature": "device_signature"        // Device identifier
//	}
//
//	@Example		{json} Request Body (API Key):
//
//	{
//	  "id": "key_01HNA...",                 // API key ID
//	  "token": "sk_live_...",               // API key secret
//	  "method": "secret",                   // Authentication method
//	  "signature": "device_signature"       // Device identifier
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "token": {
//	    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
//	    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
//	    "expiresAt": "2024-01-01T00:00:00Z"
//	  }
//	}
//
//	@Example		{json} Error Response (Invalid OTP):
//
//	{
//	  "message": "Invalid or expired OTP",
//	  "code": "AUTHENTICATION_FAILED",
//	  "details": {
//	    "reason": "incorrect code or expired verification",
//	    "remainingAttempts": 2
//	  }
//	}
//
//	@Example		{json} Error Response (Invalid API Key):
//
//	{
//	  "message": "Invalid API key",
//	  "code": "AUTHENTICATION_FAILED",
//	  "details": {
//	    "reason": "API key not found or inactive"
//	  }
//	}
//
//	@Example		{json} Error Response (Too Many Attempts):
//
//	{
//	  "message": "Too many authentication attempts",
//	  "code": "RATE_LIMIT_EXCEEDED",
//	  "details": {
//	    "retryAfter": "300s",
//	    "limit": "5 attempts per verification"
//	  }
//	}
//
//	@Router			/auth/authenticate [post]
func AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs.Logger.Info("processing authentication request",
		slog.String("handler", "AuthenticateHandler"),
		slog.String("method", r.Method),
	)

	var authenticateRequest models.AuthenticateRequest
	if err := json.NewDecoder(r.Body).Decode(&authenticateRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := authenticateRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User

	switch authenticateRequest.Method {
	case "otp":
		// Handle OTP authentication
		if err := handleOTPAuthentication(ctx, &user, authenticateRequest); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	case "secret":
		// Handle API key authentication
		if err := handleAPIKeyAuthentication(ctx, &user, authenticateRequest); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, "unsupported authentication method", http.StatusBadRequest)
		return
	}

	// Generate token
	token, err := createJWT(user.ID, authenticateRequest.Device)
	if err != nil {
		http.Error(w, "failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.AuthenticateResponse{
		Token: *token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler godoc
//	@Summary		Refresh token
//	@Description	Refresh an existing authentication token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.RefreshTokenRequest	true	"Token Refresh Request"
//	@Success		200		{object}	models.RefreshTokenResponse	"Token refreshed successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format"
//	@Failure		401		{object}	models.ErrorResponse		"Invalid or expired token"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "accessToken": "current_access_token",
//	  "refreshToken": "refresh_token"
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "token": {
//	    "accessToken": "new_access_token",
//	    "refreshToken": "new_refresh_token",
//	    "expiresAt": "2024-01-01T00:00:00Z"
//	  }
//	}
//
//	@Example		{json} Error Response (Invalid Token):
//
//	{
//	  "message": "Invalid refresh token",
//	  "code": "INVALID_TOKEN",
//	  "details": {
//	    "reason": "token expired or revoked"
//	  }
//	}
//
//	@Example		{json} Error Response (Token Mismatch):
//
//	{
//	  "message": "Token mismatch",
//	  "code": "INVALID_TOKEN",
//	  "details": {
//	    "reason": "access token does not match refresh token"
//	  }
//	}
//
//	@Example		{json} Error Response (Server Error):
//
//	{
//	  "message": "Failed to refresh token",
//	  "code": "INTERNAL_ERROR",
//	  "details": {
//	    "reason": "token storage error"
//	  }
//	}
//
//	@Router			/auth/refresh-token [post]
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs.Logger.Info("processing token refresh request",
		slog.String("handler", "RefreshTokenHandler"),
		slog.String("method", r.Method),
	)

	var refreshTokenRequest models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := refreshTokenRequest.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldToken, err := getToken(ctx, refreshTokenRequest.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if oldToken.User != userID {
		http.Error(w, "user mismatch", http.StatusUnauthorized)
		return
	}

	user := models.User{}
	if err := user.GetUser(ctx, oldToken.User); err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Generate new token
	newToken, err := createJWT(user.ID, oldToken.Device)
	if err != nil {
		http.Error(w, "failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove old token
	if err := removeToken(ctx, oldToken.ID); err != nil {
		logs.Logger.Error("failed to remove old token",
			slog.String("handler", "RefreshTokenHandler"),
			slog.String("old_token_id", oldToken.ID),
			slog.String("error", err.Error()),
		)
	}

	response := models.RefreshTokenResponse{
		Token: *newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func handleOTPAuthentication(ctx context.Context, user *models.User, req models.AuthenticateRequest) error {
	entry, err := nats.GetKV(ctx, "phoneVerification", req.ID)
	if err != nil {
		return err
	}

	var phoneNumberVerification models.PhoneNumberVerification
	if err := json.Unmarshal(entry.Value(), &phoneNumberVerification); err != nil {
		return err
	}

	valid, err := totp.ValidateCustom(
		req.Token,
		phoneNumberVerification.Secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    300,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA512,
		},
	)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("invalid OTP")
	}

	if err := user.GetUser(ctx, phoneNumberVerification.User); err != nil {
		return err
	}

	return nil
}

func handleAPIKeyAuthentication(ctx context.Context, user *models.User, req models.AuthenticateRequest) error {
	apiKey, err := models.ValidateAPIKey(ctx, req.ID, req.Token)
	if err != nil {
		return err
	}

	return user.GetUser(ctx, apiKey.UserID)
}

func getToken(ctx context.Context, tokenID string) (*models.Token, error) {
	data, err := nats.GetKV(ctx, "token", tokenID)
	if err != nil {
		return nil, err
	}

	var userToken models.Token
	if err := json.Unmarshal(data.Value(), &userToken); err != nil {
		return nil, err
	}

	return &userToken, nil
}

func removeToken(ctx context.Context, tokenID string) error {
	return nats.RemoveKV(ctx, "token", tokenID)
}

func createJWT(userID string, device string) (*models.Token, error) {
	// Generate a random token string
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	token, err := jwt.CreateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create token model
	userToken := models.Token{
		ID:                 ulid.Make().String(),
		User:               userID,
		AccessToken:        token,
		AccessTokenExpiry:  time.Now().Add(24 * time.Hour),
		RefreshTokenExpiry: time.Now().Add(30 * 24 * time.Hour),
		Device:             device,
		Service:            "infinirewards",
		CreatedAt:          time.Now(),
	}

	// Store token in NATS KV
	tokenBytes, err := json.Marshal(userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := nats.PutKV(context.Background(), "token", userToken.ID, tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}

	return &userToken, nil
}
