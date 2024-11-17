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

// MintCollectibleHandler godoc
// @Summary Mint collectible tokens
// @Description Mint new collectible tokens for a specified recipient
// @Tags collectibles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MintCollectibleRequest true "Mint Request"
// @Success 200 {object} models.MintCollectibleResponse "Mint successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Insufficient permissions to mint"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "collectibleAddress": "0x1234...",
//	  "to": "0x5678...",
//	  "tokenId": "1",
//	  "amount": "1"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Invalid Request):
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
// @Example {json} Error Response (Unauthorized):
//
//	{
//	  "message": "Missing or invalid authentication token",
//	  "code": "UNAUTHORIZED",
//	  "details": {
//	    "reason": "token expired"
//	  }
//	}
//
// @Example {json} Error Response (Forbidden):
//
//	{
//	  "message": "Insufficient permissions",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "not contract owner"
//	  }
//	}
//
// @Example {json} Error Response (Server Error):
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
// @Router /collectibles/mint [post]
func MintCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("MintCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		logs.Logger.Error("MintCollectibleHandler failed to get user", "error", err, "userId", userID)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	var mintReq models.MintCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&mintReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := mintReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert string amounts to big.Int
	tokenId, ok := new(big.Int).SetString(mintReq.TokenId, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(mintReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get account using user's private key and address from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("MintCollectibleHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	// Mint collectible
	txHash, err := infinirewards.MintCollectible(account, mintReq.CollectibleAddress, mintReq.To, tokenId, amount)
	if err != nil {
		logs.Logger.Error("MintCollectibleHandler mint error", "error", err)
		http.Error(w, "Failed to mint collectible", http.StatusInternalServerError)
		return
	}

	resp := models.MintCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleBalanceHandler godoc
// @Summary Get collectible balance
// @Description Get balance of collectible tokens for a specific token ID
// @Tags collectibles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Contract address" minlength(42) maxlength(42) format(hex)
// @Param tokenId path integer true "Token ID" minimum(0)
// @Success 200 {object} models.GetCollectibleBalanceResponse "Balance retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 404 {object} models.ErrorResponse "Contract or token not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "balance": "100"
//	}
//
// @Example {json} Error Response (Invalid Address):
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
// @Example {json} Error Response (Not Found):
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
// @Router /collectibles/{address}/balance/{tokenId} [get]
func GetCollectibleBalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleBalanceHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]
	tokenIdStr := parts[5]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		logs.Logger.Error("GetCollectibleBalanceHandler failed to get user", "error", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's private key and address from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("GetCollectibleBalanceHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	balance, err := infinirewards.BalanceOf(ctx, account, address, tokenId)
	if err != nil {
		logs.Logger.Error("GetCollectibleBalanceHandler balance error", "error", err)
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	resp := models.GetCollectibleBalanceResponse{
		Balance: balance.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleURIHandler godoc
// @Summary Get collectible URI
// @Description Get the URI for a specific collectible token's metadata
// @Tags collectibles
// @Accept json
// @Produce json
// @Param address path string true "Contract address" format(hex)
// @Param tokenId path integer true "Token ID" minimum(0)
// @Success 200 {object} models.GetCollectibleURIResponse "URI retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} models.ErrorResponse "Token not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "uri": "https://example.com/metadata/1"
//	}
//
// @Example {json} Error Response (Invalid Parameters):
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
// @Router /collectibles/{address}/uri/{tokenId} [get]
func GetCollectibleURIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleURIHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]
	tokenIdStr := parts[5]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	uri, err := infinirewards.URI(ctx, address, tokenId)
	if err != nil {
		logs.Logger.Error("GetCollectibleURIHandler URI error", "error", err)
		http.Error(w, "Failed to get URI", http.StatusInternalServerError)
		return
	}

	resp := models.GetCollectibleURIResponse{
		URI: uri,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SetTokenDataHandler godoc
// @Summary Set token data
// @Description Set metadata for a collectible token
// @Tags collectibles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Contract address" format(hex)
// @Param tokenId path integer true "Token ID" minimum(0)
// @Param request body models.SetTokenDataRequest true "Token data"
// @Success 200 {object} models.SetTokenDataResponse "Token data updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to set token data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",
//	  "price": "100",
//	  "expiry": 1735689600,
//	  "description": "Limited edition collectible"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Invalid Request):
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
// @Router /collectibles/{address}/token-data/{tokenId} [put]
func SetTokenDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("SetTokenDataHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]
	tokenIdStr := parts[5]

	var setReq models.SetTokenDataRequest
	if err := json.NewDecoder(r.Body).Decode(&setReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate request and check if path parameters match body
	setReq.CollectibleAddress = address
	setReq.TokenId = tokenIdStr
	if err := setReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	tokenId, ok := new(big.Int).SetString(setReq.TokenId, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	price, ok := new(big.Int).SetString(setReq.Price, 0)
	if !ok {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.SetTokenData(
		account,
		setReq.CollectibleAddress,
		tokenId,
		setReq.PointsContract,
		price,
		setReq.Expiry,
		setReq.Description,
	)
	if err != nil {
		logs.Logger.Error("SetTokenDataHandler set error", "error", err)
		http.Error(w, "Failed to set token data", http.StatusInternalServerError)
		return
	}

	resp := models.SetTokenDataResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Factory-related handlers

// CreateMerchantHandler godoc
// @Summary Create new merchant
// @Description Create a new merchant account with initial points contract
// @Tags merchants
// @Accept json
// @Produce json
// @Param request body models.CreateMerchantRequest true "Merchant creation request"
// @Success 201 {object} models.CreateMerchantResponse "Merchant created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 409 {object} models.ErrorResponse "Merchant already exists"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "publicKey": "0x1234...",
//	  "phoneNumber": "+60123456789",
//	  "name": "My Store",
//	  "symbol": "PTS",
//	  "decimals": 18
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "merchantAddress": "0x1234...",
//	  "pointsAddress": "0x5678..."
//	}
//
// @Example {json} Error Response (Already Exists):
//
//	{
//	  "message": "Merchant already exists",
//	  "code": "CONFLICT",
//	  "details": {
//	    "phoneNumber": "+60123456789"
//	  }
//	}
//
// @Router /merchants [post]
func CreateMerchantHandler(w http.ResponseWriter, r *http.Request) {
	logs.Logger.Info("CreateMerchantHandler called", "method", r.Method)

	var createReq models.CreateMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	txHash, merchantAddress, pointsAddress, err := infinirewards.CreateMerchant(
		createReq.PublicKey,
		createReq.PhoneNumber,
		createReq.Name,
		createReq.Symbol,
		uint64(createReq.Decimals),
	)
	if err != nil {
		logs.Logger.Error("CreateMerchantHandler create error", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create merchant: %v", err), http.StatusInternalServerError)
		return
	}

	resp := models.CreateMerchantResponse{
		TransactionHash: txHash,
		MerchantAddress: merchantAddress,
		PointsAddress:   pointsAddress,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// CreateCollectibleHandler godoc
// @Summary Create collectible contract
// @Description Create a new collectible contract for a merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateCollectibleRequest true "Collectible creation request"
// @Success 201 {object} models.CreateCollectibleResponse "Created collectible details"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to create collectibles"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "name": "Special Edition",           // Name of the collectible contract
//	  "description": "Limited collectibles" // Description of the collection
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "address": "0x1234..."
//	}
//
// @Example {json} Error Response (Invalid Request):
//
//	{
//	  "message": "Invalid request format",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "name",
//	    "reason": "name is required and must be between 1 and 100 characters"
//	  }
//	}
//
// @Example {json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "user is not a registered merchant"
//	  }
//	}
//
// @Router /merchants/collectibles [post]
func CreateCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("CreateCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		logs.Logger.Error("CreateCollectibleHandler failed to get user", "error", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	var createReq models.CreateCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("CreateCollectibleHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	txHash, address, err := infinirewards.CreateInfiniRewardsCollectible(
		account,
		createReq.Name,
		createReq.Description,
	)
	if err != nil {
		logs.Logger.Error("CreateCollectibleHandler create error", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create collectible: %v", err), http.StatusInternalServerError)
		return
	}

	resp := models.CreateCollectibleResponse{
		TransactionHash: txHash,
		Address:         address,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// CreatePointsContractHandler godoc
// @Summary Create points contract
// @Description Create an additional points contract for a merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePointsContractRequest true "Points contract creation request"
// @Success 201 {object} models.CreatePointsContractResponse "Created points contract details"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to create points contracts"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "name": "Premium Points",          // Name of the points token
//	  "symbol": "PPT",                   // Token symbol (3-4 characters)
//	  "description": "Premium rewards",   // Description of the points
//	  "decimals": "18"                   // Decimal places for the token
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "address": "0x1234..."
//	}
//
// @Example {json} Error Response (Invalid Symbol):
//
//	{
//	  "message": "Invalid token symbol",
//	  "code": "VALIDATION_ERROR",
//	  "details": {
//	    "field": "symbol",
//	    "reason": "must be 3-4 uppercase characters"
//	  }
//	}
//
// @Router /merchants/points [post]
func CreatePointsContractHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("CreatePointsContractHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		logs.Logger.Error("CreatePointsContractHandler failed to get user", "error", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	var createReq models.CreatePointsContractRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("CreatePointsContractHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	decimals, ok := new(big.Int).SetString(createReq.Decimals, 0)
	if !ok {
		http.Error(w, "Invalid decimals format", http.StatusBadRequest)
		return
	}

	txHash, address, err := infinirewards.CreateAdditionalPointsContract(
		ctx,
		account,
		createReq.Name,
		createReq.Symbol,
		createReq.Description,
		decimals,
	)
	if err != nil {
		logs.Logger.Error("CreatePointsContractHandler create error", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create points contract: %v", err), http.StatusInternalServerError)
		return
	}

	resp := models.CreatePointsContractResponse{
		TransactionHash: txHash,
		Address:         address,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Points-related handlers

// MintPointsHandler godoc
// @Summary Mint points tokens
// @Description Mint new points tokens for a specified recipient
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MintPointsRequest true "Mint Request"
// @Success 200 {object} models.MintPointsResponse "Points minted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to mint points"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "recipient": "0x5678...",       // Recipient address
//	  "amount": "100"                 // Amount to mint (in smallest unit)
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Invalid Request):
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
// @Example {json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "only merchants can mint points"
//	  }
//	}
//
// @Router /points/mint [post]
func MintPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("MintPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var mintReq models.MintPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&mintReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(mintReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("MintPointsHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.MintPoints(
		account,
		mintReq.PointsContract,
		mintReq.Recipient,
		amount,
	)
	if err != nil {
		logs.Logger.Error("MintPointsHandler mint error", "error", err)
		http.Error(w, "Failed to mint points", http.StatusInternalServerError)
		return
	}

	resp := models.MintPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BurnPointsHandler godoc
// @Summary Burn points tokens
// @Description Burn points tokens from a merchant's account
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BurnPointsRequest true "Burn Request"
// @Success 200 {object} models.BurnPointsResponse "Points burned successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to burn points"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "amount": "50"                  // Amount to burn
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Insufficient Balance):
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
// @Router /points/burn [post]
func BurnPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("BurnPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var burnReq models.BurnPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&burnReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(burnReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("BurnPointsHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.BurnPoints(
		account,
		burnReq.PointsContract,
		amount,
	)
	if err != nil {
		logs.Logger.Error("BurnPointsHandler burn error", "error", err)
		http.Error(w, "Failed to burn points", http.StatusInternalServerError)
		return
	}

	resp := models.BurnPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPointsBalanceHandler godoc
// @Summary Get points balance
// @Description Get the points balance for a specific account
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Points contract address"
// @Success 200 {object} models.GetPointsBalanceResponse "Balance retrieved"
// @Failure 400 {object} models.ErrorResponse "Invalid request format"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "balance": "150"
//	}
//
// @Router /points/{address}/balance [get]
func GetPointsBalanceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetPointsBalanceHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("GetPointsBalanceHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	balance, err := infinirewards.GetBalance(ctx, account, address)
	if err != nil {
		logs.Logger.Error("GetPointsBalanceHandler balance error", "error", err)
		http.Error(w, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	resp := models.GetPointsBalanceResponse{
		Balance: balance.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// TransferPointsHandler godoc
// @Summary Transfer points between accounts
// @Description Transfer points from one account to another
// @Tags points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.TransferPointsRequest true "Transfer Request"
// @Success 200 {object} models.TransferPointsResponse "Points transferred successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to transfer points"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "pointsContract": "0x1234...",  // Points contract address
//	  "from": "0x5678...",           // Sender address
//	  "to": "0x9abc...",             // Recipient address
//	  "amount": "25"                 // Amount to transfer
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Insufficient Balance):
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
// @Router /points/transfer [post]
func TransferPointsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("TransferPointsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var transferReq models.TransferPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(transferReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("TransferPointsHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	txHash, err := infinirewards.TransferPoints(
		ctx,
		account,
		transferReq.PointsContract,
		transferReq.From,
		transferReq.To,
		amount,
	)
	if err != nil {
		logs.Logger.Error("TransferPointsHandler transfer error", "error", err)
		http.Error(w, "Failed to transfer points", http.StatusInternalServerError)
		return
	}

	resp := models.TransferPointsResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTokenDataHandler godoc
// @Summary Get token data
// @Description Get metadata for a collectible token
// @Tags collectibles
// @Accept json
// @Produce json
// @Param address path string true "Contract address" format(hex)
// @Param tokenId path integer true "Token ID" minimum(0)
// @Success 200 {object} models.GetTokenDataResponse "Token data retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} models.ErrorResponse "Token data not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "pointsContract": "0x1234...",
//	  "price": "100",
//	  "expiry": 1735689600,
//	  "description": "Limited edition collectible"
//	}
//
// @Example {json} Error Response (Invalid Parameters):
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
// @Example {json} Error Response (Not Found):
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
// @Router /collectibles/{address}/token-data/{tokenId} [get]
func GetTokenDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetTokenDataHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]
	tokenIdStr := parts[5]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	pointsContract, price, expiry, description, err := infinirewards.GetTokenData(
		ctx,
		address,
		tokenId,
	)
	if err != nil {
		logs.Logger.Error("GetTokenDataHandler get error", "error", err)
		http.Error(w, "Failed to get token data", http.StatusInternalServerError)
		return
	}

	resp := models.GetTokenDataResponse{
		PointsContract: pointsContract,
		Price:          price.String(),
		Expiry:         int64(expiry),
		Description:    description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RedeemCollectibleHandler godoc
// @Summary Redeem collectible
// @Description Redeem a collectible token
// @Tags collectibles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Contract address" format(hex)
// @Param request body models.RedeemCollectibleRequest true "Redemption details"
// @Success 200 {object} models.RedeemCollectibleResponse "Redemption successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not authorized to redeem"
// @Failure 404 {object} models.ErrorResponse "Collectible not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "user": "0x5678...",    // User address redeeming the collectible
//	  "tokenId": "1",         // Token ID to redeem
//	  "amount": "1"           // Amount to redeem
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Invalid Request):
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
// @Example {json} Error Response (Expired):
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
// @Router /collectibles/{address}/redeem [post]
func RedeemCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("RedeemCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]

	var redeemReq models.RedeemCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&redeemReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Set address from URL path
	redeemReq.CollectibleAddress = address

	tokenId, ok := new(big.Int).SetString(redeemReq.TokenId, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(redeemReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("RedeemCollectibleHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
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
		logs.Logger.Error("RedeemCollectibleHandler redeem error", "error", err)
		http.Error(w, "Failed to redeem collectible", http.StatusInternalServerError)
		return
	}

	resp := models.RedeemCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleDetailsHandler godoc
// @Summary Get collectible details
// @Description Get detailed information about a collectible contract
// @Tags collectibles
// @Accept json
// @Produce json
// @Param address path string true "Contract address" format(hex)
// @Success 200 {object} models.GetCollectibleDetailsResponse "Collectible details retrieved successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid contract address"
// @Failure 404 {object} models.ErrorResponse "Contract not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
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
// @Example {json} Error Response (Invalid Address):
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
// @Router /collectibles/{address} [get]
func GetCollectibleDetailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleDetailsHandler called", "method", r.Method)

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[2]

	description, pointsContract, tokenIDs, tokenPrices, tokenExpiries, tokenDescriptions, err := infinirewards.GetDetails(
		ctx,
		address,
	)
	if err != nil {
		logs.Logger.Error("GetCollectibleDetailsHandler details error", "error", err)
		http.Error(w, "Failed to get collectible details", http.StatusInternalServerError)
		return
	}

	// Convert big.Int arrays to string arrays
	tokenIDStrings := make([]string, len(tokenIDs))
	tokenPriceStrings := make([]string, len(tokenPrices))
	for i, id := range tokenIDs {
		if id != nil {
			tokenIDStrings[i] = id.String()
		}
	}
	for i, price := range tokenPrices {
		if price != nil {
			tokenPriceStrings[i] = price.String()
		}
	}

	resp := models.GetCollectibleDetailsResponse{
		Description:       description,
		PointsContract:    pointsContract,
		TokenIDs:          tokenIDStrings,
		TokenPrices:       tokenPriceStrings,
		TokenExpiries:     tokenExpiries,
		TokenDescriptions: tokenDescriptions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// IsCollectibleValidHandler godoc
// @Summary Check collectible validity
// @Description Check if a collectible token is valid (not expired)
// @Tags collectibles
// @Accept json
// @Produce json
// @Param address path string true "Contract address" format(hex)
// @Param tokenId path integer true "Token ID" minimum(0)
// @Success 200 {object} models.IsCollectibleValidResponse "Validity status retrieved"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 404 {object} models.ErrorResponse "Token not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "isValid": true
//	}
//
// @Example {json} Error Response (Not Found):
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
// @Router /collectibles/{address}/valid/{tokenId} [get]
func IsCollectibleValidHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("IsCollectibleValidHandler called", "method", r.Method)

	// Extract address and tokenId from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]
	tokenIdStr := parts[5]

	tokenId, ok := new(big.Int).SetString(tokenIdStr, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	isValid, err := infinirewards.IsValid(ctx, address, tokenId)
	if err != nil {
		logs.Logger.Error("IsCollectibleValidHandler validation error", "error", err)
		http.Error(w, "Failed to check collectible validity", http.StatusInternalServerError)
		return
	}

	resp := models.IsCollectibleValidResponse{
		IsValid: isValid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PurchaseCollectibleHandler godoc
// @Summary Purchase collectible
// @Description Purchase a collectible token using points
// @Tags collectibles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param address path string true "Contract address" format(hex)
// @Param request body models.PurchaseCollectibleRequest true "Purchase details"
// @Success 200 {object} models.PurchaseCollectibleResponse "Purchase successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request format or validation failed"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 402 {object} models.ErrorResponse "Insufficient points balance"
// @Failure 404 {object} models.ErrorResponse "Collectible not found"
// @Failure 410 {object} models.ErrorResponse "Collectible expired"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Request Body:
//
//	{
//	  "user": "0x5678...",
//	  "tokenId": "1",
//	  "amount": "1"
//	}
//
// @Example {json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc..."
//	}
//
// @Example {json} Error Response (Insufficient Points):
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
// @Example {json} Error Response (Expired):
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
// @Router /collectibles/{address}/purchase [post]
func PurchaseCollectibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("PurchaseCollectibleHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract address from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	address := parts[3]

	var purchaseReq models.PurchaseCollectibleRequest
	if err := json.NewDecoder(r.Body).Decode(&purchaseReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Set address from URL path
	purchaseReq.CollectibleAddress = address

	tokenId, ok := new(big.Int).SetString(purchaseReq.TokenId, 0)
	if !ok {
		http.Error(w, "Invalid tokenId format", http.StatusBadRequest)
		return
	}

	amount, ok := new(big.Int).SetString(purchaseReq.Amount, 0)
	if !ok {
		http.Error(w, "Invalid amount format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("PurchaseCollectibleHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
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
		logs.Logger.Error("PurchaseCollectibleHandler purchase error", "error", err)
		http.Error(w, "Failed to purchase collectible", http.StatusInternalServerError)
		return
	}

	resp := models.PurchaseCollectibleResponse{
		TransactionHash: txHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Merchant-related handlers

// GetPointsContractsHandler godoc
// @Summary Get merchant's points contracts
// @Description Get all points contracts associated with a merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.GetPointsContractsResponse "List of points contracts"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not a merchant account"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "contracts": [
//	    {
//	      "address": "0x1234...",
//	      "name": "Store Points",
//	      "symbol": "SP",
//	      "decimals": "18",
//	      "totalSupply": "1000000"
//	    }
//	  ]
//	}
//
// @Example {json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "account is not a merchant"
//	  }
//	}
//
// @Router /merchants/points-contracts [get]
func GetPointsContractsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetPointsContractsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("GetPointsContractsHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	contracts, err := infinirewards.GetPointsContracts(ctx, account)
	if err != nil {
		logs.Logger.Error("GetPointsContractsHandler contracts error", "error", err)
		http.Error(w, "Failed to get points contracts", http.StatusInternalServerError)
		return
	}

	// Convert []string to []PointsContractInfo
	contractInfos := make([]models.PointsContractInfo, len(contracts))
	for i, addr := range contracts {
		// Get contract details
		name, symbol, description, decimals, totalSupply, err := infinirewards.GetPointsContractDetails(ctx, addr)
		if err != nil {
			logs.Logger.Error("GetPointsContractsHandler details error", "error", err, "address", addr)
			continue
		}

		contractInfos[i] = models.PointsContractInfo{
			Address:     addr,
			Name:        name,
			Symbol:      symbol,
			Description: description,
			Decimals:    uint8(decimals),
			TotalSupply: uint64(totalSupply),
		}
	}

	resp := models.GetPointsContractsResponse{
		Contracts: contractInfos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectibleContractsHandler godoc
// @Summary Get merchant's collectible contracts
// @Description Get all collectible contracts associated with a merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.GetCollectibleContractsResponse "List of collectible contracts"
// @Failure 401 {object} models.ErrorResponse "Missing or invalid authentication token"
// @Failure 403 {object} models.ErrorResponse "Not a merchant account"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Example {json} Success Response:
//
//	{
//	  "contracts": [
//	    {
//	      "address": "0x1234...",
//	      "name": "Special Edition",
//	      "description": "Limited collectibles",
//	      "totalSupply": "100",
//	      "tokenTypes": ["1","2","3"],
//	      "tokenPrices": ["100","200","300"],
//	      "tokenExpiries": [1735689600,1735689600,1735689600],
//	      "tokenDescriptions": ["Gold","Silver","Bronze"]
//	    }
//	  ]
//	}
//
// @Example {json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "account is not a merchant"
//	  }
//	}
//
// @Router /merchants/collectible-contracts [get]
func GetCollectibleContractsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetCollectibleContractsHandler called", "method", r.Method)

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Use user's credentials from database
	account, err := infinirewards.GetAccount(user.PrivateKey, user.AccountAddress)
	if err != nil {
		logs.Logger.Error("GetCollectibleContractsHandler account error", "error", err)
		http.Error(w, "Failed to get account", http.StatusInternalServerError)
		return
	}

	contracts, err := infinirewards.GetCollectibleContracts(ctx, account)
	if err != nil {
		logs.Logger.Error("GetCollectibleContractsHandler contracts error", "error", err)
		http.Error(w, "Failed to get collectible contracts", http.StatusInternalServerError)
		return
	}

	// Convert []string to []CollectibleContractInfo
	contractInfos := make([]models.CollectibleContractInfo, len(contracts))
	for i, addr := range contracts {
		// Get contract details - GetDetails returns 8 values:
		// description, pointsContract, _, tokenIDs, tokenPrices, tokenExpiries, tokenDescriptions, err
		description, pointsContract, tokenIDs, tokenPrices, tokenExpiries, tokenDescriptions, err := infinirewards.GetDetails(
			ctx,
			addr,
		)
		if err != nil {
			logs.Logger.Error("GetCollectibleContractsHandler details error", "error", err, "address", addr)
			continue
		}

		// Convert tokenIDs and tokenPrices from []*big.Int to []string
		tokenIDStrings := make([]string, len(tokenIDs))
		tokenPriceStrings := make([]string, len(tokenPrices))
		for j, id := range tokenIDs {
			if id != nil {
				tokenIDStrings[j] = id.String()
			}
		}
		for j, price := range tokenPrices {
			if price != nil {
				tokenPriceStrings[j] = price.String()
			}
		}

		// Get total supply by counting valid tokens
		totalSupply := "0"
		if len(tokenIDStrings) > 0 {
			totalSupply = fmt.Sprintf("%d", len(tokenIDStrings))
		}

		contractInfos[i] = models.CollectibleContractInfo{
			Address:           addr,
			Name:              pointsContract, // Using pointsContract as name since GetDetails doesn't return name
			Description:       description,
			PointsContract:    pointsContract,
			TotalSupply:       totalSupply,
			TokenTypes:        tokenIDStrings,
			TokenPrices:       tokenPriceStrings,
			TokenExpiries:     tokenExpiries,
			TokenDescriptions: tokenDescriptions,
		}
	}

	resp := models.GetCollectibleContractsResponse{
		Contracts: contractInfos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
