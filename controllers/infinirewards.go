package controllers

import (
	"encoding/json"
	"fmt"
	"infinirewards/infinirewards"
	"infinirewards/logs"
	"infinirewards/middleware"
	"infinirewards/models"
	"math/big"
	"net/http"
	"strings"
)

// GetCollectibleBalanceHandler godoc
//
//	@Summary		Get collectible balance
//	@Metadata	Get balance of collectible tokens for a specific token ID
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			address	path		string									true	"Contract address"	minlength(42)	maxlength(42)	format(hex)
//	@Param			tokenId	path		integer									true	"Token ID"			minimum(0)
//	@Success		200		{object}	models.GetCollectibleBalanceResponse	"Balance retrieved successfully"
//	@Failure		400		{object}	models.ErrorResponse					"Invalid request parameters"
//	@Failure		401		{object}	models.ErrorResponse					"Missing or invalid authentication token"
//	@Failure		404		{object}	models.ErrorResponse					"Contract or token not found"
//	@Failure		500		{object}	models.ErrorResponse					"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "balance": "100"
//	}
//
//	@Example		{json} Error Response (Invalid Address):
//
//	{
//	  "message": "Invalid contract address",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "address",
//	    "reason": "must be a valid hex address"
//	  }
//	}
//
//	@Example		{json} Error Response (Not Found):
//
//	{
//	  "message": "Contract or token not found",
//	  "code": "NOT_FOUND",
//	  "details": {
//	    "address": "0x1234...",
//	    "tokenId": "1"
//	  }
//	}
//
//	@Router			/collectibles/{address}/balance/{tokenId} [get]
func GetCollectibleBalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleBalanceHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]
	tokenIdStr := parts[4]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": tokenIdStr,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	balance, err := infinirewards.BalanceOf(ctx, user.AccountAddress, address, tokenId)
	if err != nil {
		WriteError(w, "Failed to get balance", InternalServerError, map[string]string{
			"reason": "Failed to retrieve balance from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.GetCollectibleBalanceResponse{
		Balance: balance.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleURIHandler godoc
//
//	@Summary		Get collectible URI
//	@Metadata	Get the URI for a specific collectible token's metadata
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string								true	"Contract address"	format(hex)
//	@Param			tokenId	path		integer								true	"Token ID"			minimum(0)
//	@Success		200		{object}	models.GetCollectibleURIResponse	"URI retrieved successfully"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request parameters"
//	@Failure		404		{object}	models.ErrorResponse				"Token not found"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "uri": "https://example.com/metadata/1"
//	}
//
//	@Example		{json} Error Response (Invalid Parameters):
//
//	{
//	  "message": "Invalid parameters",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "address": "invalid hex format",
//	    "tokenId": "must be a non-negative integer"
//	  }
//	}
//
//	@Router			/collectibles/{address}/uri/{tokenId} [get]
func GetCollectibleURIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleURIHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]
	tokenIdStr := parts[4]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": tokenIdStr,
		}, http.StatusBadRequest)
		return
	}

	uri, err := infinirewards.URI(ctx, address, tokenId)
	if err != nil {
		WriteError(w, "Failed to get URI", InternalServerError, map[string]string{
			"reason": "Failed to retrieve URI from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.GetCollectibleURIResponse{
		URI: uri,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// MintCollectibleHandler godoc
//
//	@Summary		Mint collectible tokens
//	@Metadata	Mint new collectible tokens for a specified recipient
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.MintCollectibleRequest	true	"Mint Request"
//	@Success		200		{object}	models.MintCollectibleResponse	"Mint successful"
//	@Failure		400		{object}	models.ErrorResponse			"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse			"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse			"Insufficient permissions to mint"
//	@Failure		500		{object}	models.ErrorResponse			"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "collectibleAddress": "0x1234...",
//	  "to": "0x5678...",
//	  "tokenId": "1",
//	  "amount": "1"
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "amount",
//	    "reason": "must be a positive number"
//	  }
//	}
//
//	@Example		{json} Error Response (Unauthorized):
//
//	{
//	  "message": "Missing or invalid authentication token",
//	  "code": "UNAUTHORIZED",
//	  "details": {
//	    "reason": "token expired"
//	  }
//	}
//
//	@Example		{json} Error Response (Forbidden):
//
//	{
//	  "message": "Insufficient permissions",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "not contract owner"
//	  }
//	}
//
//	@Example		{json} Error Response (Server Error):
//
//	{
//	  "message": "Failed to mint collectible",
//	  "code": "INTERNAL_ERROR",
//	  "details": {
//	    "reason": "blockchain transaction failed",
//	    "txHash": "0xdef..."
//	  }
//	}
//
//	@Router			/merchant/collectibles/mint [post]
func MintCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("MintCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	var mintReq models.MintCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&mintReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	tokenId, ok := new(big.Int).SetString(mintReq.TokenId, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": mintReq.TokenId,
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(mintReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": mintReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		WriteError(w, "Not authorized", AuthorizationError, map[string]string{
			"reason": "User is not a merchant",
		}, http.StatusForbidden)
		return
	}

	// Get account using user's private key and address from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	// Mint collectible
	txHash, err := infinirewards.MintCollectible(account, mintReq.CollectibleAddress, mintReq.To, tokenId, amount)
	if err != nil {
		WriteError(w, "Failed to mint collectible", InternalServerError, map[string]string{
			"reason": "Failed to mint collectible on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.MintCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SetTokenDataHandler godoc
//
//	@Summary		Set token data
//	@Metadata	Set metadata for a collectible token
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			address	path		string						true	"Contract address"	format(hex)
//	@Param			tokenId	path		integer						true	"Token ID"			minimum(0)
//	@Param			request	body		models.SetTokenDataRequest	true	"Token data"
//	@Success		200		{object}	models.SetTokenDataResponse	"Token data updated successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse		"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse		"Not authorized to set token data"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",
//	  "price": "100",
//	  "expiry": 1735689600,
//	  "description": "Limited edition collectible"
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "price",
//	    "reason": "must be a positive number"
//	  }
//	}
//
//	@Router			/collectibles/{address}/token-data/{tokenId} [put]
func SetTokenDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("SetTokenDataHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]
	tokenIdStr := parts[4]

	var setReq models.SetTokenDataRequest
	if err := json.NewDecoder(r.Body).Decode(&setReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	// Validate request and check if path parameters match body
	setReq.CollectibleAddress = address
	setReq.TokenId = tokenIdStr
	if err := setReq.Validate(); err != nil {
		WriteError(w, "Validation failed", ValidationError, map[string]string{
			"reason": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	tokenId, ok := new(big.Int).SetString(setReq.TokenId, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": setReq.TokenId,
		}, http.StatusBadRequest)
		return
	}

	price, ok := new(big.Int).SetString(setReq.Price, 0)
	if !ok {
		WriteError(w, "Invalid price format", ValidationError, map[string]string{
			"reason": "Price must be a valid number",
			"price":  setReq.Price,
		}, http.StatusBadRequest)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, user.ID); err != nil {
		WriteError(w, "Failed to get merchant", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.SetTokenData(
		account,
		setReq.CollectibleAddress,
		tokenId,
		setReq.PointsContract,
		price,
		setReq.Expiry,
		setReq.Metadata,
	)
	if err != nil {
		WriteError(w, "Failed to set token data", InternalServerError, map[string]string{
			"reason": "Failed to set token data on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.SetTokenDataResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Points-related handlers

// MintPointsHandler godoc
//
//	@Summary		Mint points tokens
//	@Metadata	Mint new points tokens for a specified recipient
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.MintPointsRequest	true	"Mint Request"
//	@Success		200		{object}	models.MintPointsResponse	"Points minted successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse		"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse		"Not authorized to mint points"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "recipient": "0x5678...",       // Recipient address
//	  "amount": "100"                 // Amount to mint (in smallest unit)
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "amount",
//	    "reason": "must be a positive number"
//	  }
//	}
//
//	@Example		{json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "only merchants can mint points"
//	  }
//	}
//
//	@Router			/points/mint [post]
func MintPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("MintPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	var mintReq models.MintPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&mintReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(mintReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount format", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": mintReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		WriteError(w, "Not authorized", AuthorizationError, map[string]string{
			"reason": "User is not a merchant",
		}, http.StatusForbidden)
		return
	}

	// Use merchant's address for account
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.MintPoints(
		account,
		mintReq.PointsContract,
		mintReq.Recipient,
		amount,
	)
	if err != nil {
		WriteError(w, "Failed to mint points", InternalServerError, map[string]string{
			"reason": "Failed to mint points on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.MintPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BurnPointsHandler godoc
//
//	@Summary		Burn points tokens
//	@Metadata	Burn points tokens from a merchant's account
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.BurnPointsRequest	true	"Burn Request"
//	@Success		200		{object}	models.BurnPointsResponse	"Points burned successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse		"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse		"Not authorized to burn points"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "amount": "50"                  // Amount to burn
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Insufficient Balance):
//
//	{
//	  "message": "Insufficient balance",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "available": "40",
//	    "requested": "50"
//	  }
//	}
//
//	@Router			/points/burn [post]
func BurnPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("BurnPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	var burnReq models.BurnPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&burnReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(burnReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount format", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": burnReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, user.AccountAddress)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.BurnPoints(
		account,
		burnReq.PointsContract,
		amount,
	)
	if err != nil {
		WriteError(w, "Failed to burn points", InternalServerError, map[string]string{
			"reason": "Failed to burn points on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.BurnPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPointsBalanceHandler godoc
//
//	@Summary		Get points balance
//	@Metadata	Get the points balance for a specific account
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			address	path		string							true	"Points contract address"
//	@Success		200		{object}	models.GetPointsBalanceResponse	"Balance retrieved"
//	@Failure		400		{object}	models.ErrorResponse			"Invalid request format"
//	@Failure		401		{object}	models.ErrorResponse			"Unauthorized"
//	@Failure		500		{object}	models.ErrorResponse			"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "balance": "150"
//	}
//
//	@Router			/points/{address}/balance [get]
func GetPointsBalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetPointsBalanceHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, user.AccountAddress)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	name, symbol, description, decimals, _, err := infinirewards.GetPointsContractDetails(ctx, address)
	if err != nil {
		WriteError(w, "Failed to get points contract details", InternalServerError, map[string]string{
			"reason": "Failed to retrieve points contract details from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	balance, err := infinirewards.GetBalance(ctx, account, address)
	if err != nil {
		WriteError(w, "Failed to get balance", InternalServerError, map[string]string{
			"reason": "Failed to retrieve balance from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.GetPointsBalanceResponse{
		Balance:  balance.String(),
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
		Metadata: description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// TransferPointsHandler godoc
//
//	@Summary		Transfer points between accounts
//	@Metadata	Transfer points from one account to another
//	@Tags			points
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.TransferPointsRequest	true	"Transfer Request"
//	@Success		200		{object}	models.TransferPointsResponse	"Points transferred successfully"
//	@Failure		400		{object}	models.ErrorResponse			"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse			"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse			"Not authorized to transfer points"
//	@Failure		500		{object}	models.ErrorResponse			"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "from": "0x5678...",           // Sender address
//	  "to": "0x9abc...",             // Recipient address
//	  "amount": "25"                 // Amount to transfer
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Insufficient Balance):
//
//	{
//	  "message": "Insufficient balance",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "available": "20",
//	    "requested": "25"
//	  }
//	}
//
//	@Router			/points/transfer [post]
func TransferPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("TransferPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	var transferReq models.TransferPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(transferReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount format", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": transferReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, user.AccountAddress)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.TransferPoints(
		ctx,
		account,
		transferReq.PointsContract,
		transferReq.To,
		amount,
	)
	if err != nil {
		WriteError(w, "Failed to transfer points", InternalServerError, map[string]string{
			"reason": "Failed to transfer points on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.TransferPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTokenDataHandler godoc
//
//	@Summary		Get token data
//	@Metadata	Get metadata for a collectible token
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string						true	"Contract address"	format(hex)
//	@Param			tokenId	path		integer						true	"Token ID"			minimum(0)
//	@Success		200		{object}	models.GetTokenDataResponse	"Token data retrieved successfully"
//	@Failure		400		{object}	models.ErrorResponse		"Invalid request parameters"
//	@Failure		404		{object}	models.ErrorResponse		"Token data not found"
//	@Failure		500		{object}	models.ErrorResponse		"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "pointsContract": "0x1234...",
//	  "price": "100",
//	  "expiry": 1735689600,
//	  "description": "Limited edition collectible"
//	}
//
//	@Example		{json} Error Response (Invalid Parameters):
//
//	{
//	  "message": "Invalid parameters",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "address": "invalid hex format",
//	    "tokenId": "must be a non-negative integer"
//	  }
//	}
//
//	@Example		{json} Error Response (Not Found):
//
//	{
//	  "message": "Token data not found",
//	  "code": "NOT_FOUND",
//	  "details": {
//	    "address": "0x1234...",
//	    "tokenId": "1"
//	  }
//	}
//
//	@Router			/collectibles/{address}/token-data/{tokenId} [get]
func GetTokenDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetTokenDataHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]
	tokenIdStr := parts[4]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": tokenIdStr,
		}, http.StatusBadRequest)
		return
	}

	pointsContract, price, expiry, description, err := infinirewards.GetTokenData(
		ctx,
		address,
		tokenId,
	)
	if err != nil {
		WriteError(w, "Failed to get token data", InternalServerError, map[string]string{
			"reason": "Failed to retrieve token data from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.GetTokenDataResponse{
		PointsContract: pointsContract,
		Price:          price.String(),
		Expiry:         int64(expiry),
		Metadata:       description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RedeemCollectibleHandler godoc
//
//	@Summary		Redeem collectible
//	@Metadata	Redeem a collectible token
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			address	path		string								true	"Contract address"	format(hex)
//	@Param			request	body		models.RedeemCollectibleRequest		true	"Redemption details"
//	@Success		200		{object}	models.RedeemCollectibleResponse	"Redemption successful"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse				"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse				"Not authorized to redeem"
//	@Failure		404		{object}	models.ErrorResponse				"Collectible not found"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "user": "0x5678...",    // User address redeeming the collectible
//	  "tokenId": "1",         // Token ID to redeem
//	  "amount": "1"           // Amount to redeem
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "amount",
//	    "reason": "must be a positive integer"
//	  }
//	}
//
//	@Example		{json} Error Response (Expired):
//
//	{
//	  "message": "Collectible expired",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "expiry": "2024-01-01T00:00:00Z",
//	    "current": "2024-02-01T00:00:00Z"
//	  }
//	}
//
//	@Router			/collectibles/{address}/redeem [post]
func RedeemCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("RedeemCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]

	var redeemReq models.RedeemCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&redeemReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	// Set address from URL path
	redeemReq.CollectibleAddress = address

	tokenId, ok := new(big.Int).SetString(redeemReq.TokenId, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": redeemReq.TokenId,
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(redeemReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount format", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": redeemReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		WriteError(w, "Not authorized", AuthorizationError, map[string]string{
			"reason": "User is not a merchant",
		}, http.StatusForbidden)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.Redeem(
		ctx,
		account,
		redeemReq.CollectibleAddress,
		redeemReq.User,
		tokenId,
		amount,
	)
	if err != nil {
		WriteError(w, "Failed to redeem collectible", InternalServerError, map[string]string{
			"reason": "Failed to redeem collectible on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.RedeemCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleDetailsHandler godoc
//
//	@Summary		Get collectible details
//	@Metadata	Get detailed information about a collectible contract
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string									true	"Contract address"	format(hex)
//	@Success		200		{object}	models.GetCollectibleDetailsResponse	"Collectible details retrieved successfully"
//	@Failure		400		{object}	models.ErrorResponse					"Invalid contract address"
//	@Failure		404		{object}	models.ErrorResponse					"Contract not found"
//	@Failure		500		{object}	models.ErrorResponse					"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "description": "Special Edition Collectibles",
//	  "pointsContract": "0x1234...",
//	  "tokenIDs": ["1", "2", "3"],
//	  "tokenPrices": ["100", "200", "300"],
//	  "tokenExpiries": [1735689600, 1735689600, 1735689600],
//	  "tokenDescriptions": ["Gold", "Silver", "Bronze"]
//	}
//
//	@Example		{json} Error Response (Invalid Address):
//
//	{
//	  "message": "Invalid contract address",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "address",
//	    "reason": "must be a valid hex address"
//	  }
//	}
//
//	@Router			/collectibles/{address} [get]
func GetCollectibleDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleDetailsHandler called", "method", r.Method)

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]

	name, description, pointsContract, tokenIDs, tokenPrices, tokenExpiries, tokenDescriptions, tokenSupplies, err := infinirewards.GetDetails(
		ctx,
		address,
	)
	if err != nil {
		WriteError(w, "Failed to get collectible details", InternalServerError, map[string]string{
			"reason": "Failed to retrieve collectible details from blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Convert big.Int arrays to string arrays
	tokenIDStrings := make([]string, len(tokenIDs))
	tokenPriceStrings := make([]string, len(tokenPrices))
	tokenBalanceStrings := make([]string, len(tokenIDs))
	tokenSupplyStrings := make([]string, len(tokenSupplies))
	for i, id := range tokenIDs {
		if id != nil {
			tokenIDStrings[i] = id.String()
			balance, err := infinirewards.BalanceOf(ctx, user.AccountAddress, address, id)
			if err != nil {
				WriteError(w, "Failed to get token balance", InternalServerError, map[string]string{
					"reason": "Failed to get token balance from blockchain",
					"error":  err.Error(),
				}, http.StatusInternalServerError)
				return
			}
			tokenBalanceStrings[i] = balance.String()
			tokenSupplyStrings[i] = fmt.Sprintf("%d", tokenSupplies[i])
		}
	}
	for i, price := range tokenPrices {
		if price != nil {
			tokenPriceStrings[i] = price.String()
		}
	}

	resp := models.GetCollectibleDetailsResponse{
		Name:              name,
		Address:           address,
		Metadata:          description,
		PointsContract:    pointsContract,
		TokenIDs:          tokenIDStrings,
		TokenPrices:       tokenPriceStrings,
		TokenExpiries:     tokenExpiries,
		TokenDescriptions: tokenDescriptions,
		TokenBalances:     tokenBalanceStrings,
		TokenSupplies:     tokenSupplyStrings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// IsCollectibleValidHandler godoc
//
//	@Summary		Check collectible validity
//	@Metadata	Check if a collectible token is valid (not expired)
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string								true	"Contract address"	format(hex)
//	@Param			tokenId	path		integer								true	"Token ID"			minimum(0)
//	@Success		200		{object}	models.IsCollectibleValidResponse	"Validity status retrieved"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request parameters"
//	@Failure		404		{object}	models.ErrorResponse				"Token not found"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Success Response:
//
//	{
//	  "isValid": true
//	}
//
//	@Example		{json} Error Response (Not Found):
//
//	{
//	  "message": "Token not found",
//	  "code": "NOT_FOUND",
//	  "details": {
//	    "address": "0x1234...",
//	    "tokenId": "1"
//	  }
//	}
//
//	@Router			/collectibles/{address}/valid/{tokenId} [get]
func IsCollectibleValidHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("IsCollectibleValidHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]
	tokenIdStr := parts[4]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": tokenIdStr,
		}, http.StatusBadRequest)
		return
	}

	isValid, err := infinirewards.IsValid(ctx, address, tokenId)
	if err != nil {
		WriteError(w, "Failed to check collectible validity", InternalServerError, map[string]string{
			"reason": "Failed to check collectible validity on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.IsCollectibleValidResponse{
		IsValid: isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PurchaseCollectibleHandler godoc
//
//	@Summary		Purchase collectible
//	@Metadata	Purchase a collectible token using points
//	@Tags			collectibles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			address	path		string								true	"Contract address"	format(hex)
//	@Param			request	body		models.PurchaseCollectibleRequest	true	"Purchase details"
//	@Success		200		{object}	models.PurchaseCollectibleResponse	"Purchase successful"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse				"Missing or invalid authentication token"
//	@Failure		402		{object}	models.ErrorResponse				"Insufficient points balance"
//	@Failure		404		{object}	models.ErrorResponse				"Collectible not found"
//	@Failure		410		{object}	models.ErrorResponse				"Collectible expired"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "user": "0x5678...",
//	  "tokenId": "1",
//	  "amount": "1"
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
//	@Example		{json} Error Response (Insufficient Points):
//
//	{
//	  "message": "Insufficient points balance",
//	  "code": "PAYMENT_REQUIRED",
//	  "details": {
//	    "required": "100",
//	    "available": "50"
//	  }
//	}
//
//	@Example		{json} Error Response (Expired):
//
//	{
//	  "message": "Collectible expired",
//	  "code": "GONE",
//	  "details": {
//	    "expiry": "2024-01-01T00:00:00Z",
//	    "current": "2024-02-01T00:00:00Z"
//	  }
//	}
//
//	@Router			/collectibles/{address}/purchase [post]
func PurchaseCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("PurchaseCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		WriteError(w, "Unauthorized", AuthenticationError, map[string]string{
			"reason": "Missing or invalid authentication token",
		}, http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		WriteError(w, "Invalid URL format", ValidationError, map[string]string{
			"reason": "Missing required path parameters",
		}, http.StatusBadRequest)
		return
	}
	address := parts[2]

	var purchaseReq models.PurchaseCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&purchaseReq); err != nil {
		WriteError(w, "Invalid request format", ValidationError, map[string]string{
			"reason": "Unable to parse JSON request",
		}, http.StatusBadRequest)
		return
	}

	// Set address from URL path
	purchaseReq.CollectibleAddress = address

	tokenId, ok := new(big.Int).SetString(purchaseReq.TokenId, 0)
	if !ok {
		WriteError(w, "Invalid token ID", ValidationError, map[string]string{
			"reason":  "Token ID must be a valid number",
			"tokenId": purchaseReq.TokenId,
		}, http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(purchaseReq.Amount, 0)
	if !ok {
		WriteError(w, "Invalid amount format", ValidationError, map[string]string{
			"reason": "Amount must be a valid number",
			"amount": purchaseReq.Amount,
		}, http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		WriteError(w, "Failed to get user", InternalServerError, map[string]string{
			"reason": "Database operation failed",
		}, http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, user.AccountAddress)
	if err != nil {
		WriteError(w, "Failed to get account", InternalServerError, map[string]string{
			"reason": "Failed to get blockchain account",
		}, http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.Purchase(
		ctx,
		account,
		purchaseReq.CollectibleAddress,
		purchaseReq.User,
		tokenId,
		amount,
	)
	if err != nil {
		WriteError(w, "Failed to purchase collectible", InternalServerError, map[string]string{
			"reason": "Failed to purchase collectible on blockchain",
			"error":  err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	resp := models.PurchaseCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
