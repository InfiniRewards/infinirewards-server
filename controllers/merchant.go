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
	"time"
)

// Factory-related handlers

// CreateMerchantHandler godoc
//	@Summary		Create new merchant
//	@Description	Create a new merchant account with initial points contract
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.CreateMerchantRequest	true	"Merchant creation request"
//	@Success		201		{object}	models.CreateMerchantResponse	"Merchant created successfully"
//	@Failure		400		{object}	models.ErrorResponse			"Invalid request format or validation failed"
//	@Failure		409		{object}	models.ErrorResponse			"Merchant already exists"
//	@Failure		500		{object}	models.ErrorResponse			"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "publicKey": "0x1234...",
//	  "phoneNumber": "+60123456789",
//	  "name": "My Store",
//	  "symbol": "PTS",
//	  "decimals": 18
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "merchantAddress": "0x1234...",
//	  "pointsAddress": "0x5678..."
//	}
//
//	@Example		{json} Error Response (Already Exists):
//
//	{
//	  "message": "Merchant already exists",
//	  "code": "CONFLICT",
//	  "details": {
//	    "phoneNumber": "+60123456789"
//	  }
//	}
//
//	@Router			/merchant [post]
func CreateMerchantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("CreateMerchantHandler called", "method", r.Method)

	var createReq models.CreateMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user := &models.User{}
	if err := user.GetUser(ctx, userID); err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	txHash, merchantAddress, pointsAddress, err := infinirewards.CreateMerchant(
		user.PublicKey,
		user.PhoneNumber,
		createReq.Name,
		createReq.Symbol,
		uint64(createReq.Decimals),
	)
	if err != nil {
		logs.Logger.Error("CreateMerchantHandler create error", "error", err)
		http.Error(w, fmt.Sprintf("Failed to create merchant: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = infinirewards.FundAccount(merchantAddress)
	if err != nil {
		logs.Logger.Error("CreateMerchantHandler fund error", "error", err)
		http.Error(w, "Failed to fund merchant", http.StatusInternalServerError)
		return
	}

	merchant := &models.Merchant{
		Address:   merchantAddress,
		Name:      createReq.Name,
		Symbol:    createReq.Symbol,
		Decimals:  createReq.Decimals,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := merchant.CreateMerchant(ctx, user); err != nil {
		logs.Logger.Error("CreateMerchantHandler create error", "error", err)
		http.Error(w, "Failed to create merchant", http.StatusInternalServerError)
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
//	@Summary		Create collectible contract
//	@Description	Create a new collectible contract for a merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreateCollectibleRequest		true	"Collectible creation request"
//	@Success		201		{object}	models.CreateCollectibleResponse	"Created collectible details"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse				"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse				"Not authorized to create collectibles"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "name": "Special Edition",           // Name of the collectible contract
//	  "description": "Limited collectibles" // Description of the collection
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "address": "0x1234..."
//	}
//
//	@Example		{json} Error Response (Invalid Request):
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
//	@Example		{json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "user is not a registered merchant"
//	  }
//	}
//
//	@Router			/merchant/collectibles [post]
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

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, user.ID); err != nil {
		logs.Logger.Error("CreateCollectibleHandler failed to get merchant", "error", err)
		http.Error(w, "Failed to get merchant", http.StatusBadRequest)
		return
	}

	// Use merchant's address for account
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
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
//	@Summary		Create points contract
//	@Description	Create an additional points contract for a merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		models.CreatePointsContractRequest	true	"Points contract creation request"
//	@Success		201		{object}	models.CreatePointsContractResponse	"Created points contract details"
//	@Failure		400		{object}	models.ErrorResponse				"Invalid request format or validation failed"
//	@Failure		401		{object}	models.ErrorResponse				"Missing or invalid authentication token"
//	@Failure		403		{object}	models.ErrorResponse				"Not authorized to create points contracts"
//	@Failure		500		{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Request Body:
//
//	{
//	  "name": "Premium Points",          // Name of the points token
//	  "symbol": "PPT",                   // Token symbol (3-4 characters)
//	  "description": "Premium rewards",   // Description of the points
//	  "decimals": "18"                   // Decimal places for the token
//	}
//
//	@Example		{json} Success Response:
//
//	{
//	  "transactionHash": "0x9abc...",
//	  "address": "0x1234..."
//	}
//
//	@Example		{json} Error Response (Invalid Symbol):
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
//	@Router			/merchant/points [post]
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

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		logs.Logger.Error("CreatePointsContractHandler failed to get merchant", "error", err)
		http.Error(w, "Failed to get merchant", http.StatusBadRequest)
		return
	}

	// Use merchant's address for account
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
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

// GetPointsContractsHandler godoc
//	@Summary		Get merchant's points contracts
//	@Description	Get all points contracts associated with a merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.GetPointsContractsResponse	"List of points contracts"
//	@Failure		401	{object}	models.ErrorResponse				"Missing or invalid authentication token"
//	@Failure		403	{object}	models.ErrorResponse				"Not a merchant account"
//	@Failure		500	{object}	models.ErrorResponse				"Internal server error"
//	@Example		{json} Success Response:
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
//	@Example		{json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "account is not a merchant"
//	  }
//	}
//
//	@Router			/merchant/points-contracts [get]
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
		http.Error(w, "Failed to get user", http.StatusBadRequest)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		logs.Logger.Error("GetPointsContractsHandler failed to get merchant", "error", err)
		http.Error(w, "Failed to get merchant", http.StatusBadRequest)
		return
	}

	// Use merchant's address for account
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
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
//	@Summary		Get merchant's collectible contracts
//	@Description	Get all collectible contracts associated with a merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.GetCollectibleContractsResponse	"List of collectible contracts"
//	@Failure		401	{object}	models.ErrorResponse					"Missing or invalid authentication token"
//	@Failure		403	{object}	models.ErrorResponse					"Not a merchant account"
//	@Failure		500	{object}	models.ErrorResponse					"Internal server error"
//	@Example		{json} Success Response:
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
//	@Example		{json} Error Response (Not Merchant):
//
//	{
//	  "message": "Not authorized",
//	  "code": "FORBIDDEN",
//	  "details": {
//	    "reason": "account is not a merchant"
//	  }
//	}
//
//	@Router			/merchant/collectible-contracts [get]
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

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		logs.Logger.Error("GetCollectibleContractsHandler failed to get merchant", "error", err)
		http.Error(w, "Failed to get merchant", http.StatusBadRequest)
		return
	}

	// Use merchant's address for account
	account, err := infinirewards.GetAccount(user.PrivateKey, user.PublicKey, merchant.Address)
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

// GetMerchantHandler godoc
//	@Summary		Get merchant details
//	@Description	Get details of a merchant
//	@Tags			merchants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	models.Merchant			"Merchant details"
//	@Failure		401	{object}	models.ErrorResponse	"Missing or invalid authentication token"
//	@Failure		403	{object}	models.ErrorResponse	"Not a merchant account"
//	@Failure		500	{object}	models.ErrorResponse	"Internal server error"
//	@Router			/merchant [get]
func GetMerchantHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logs.Logger.Info("GetMerchantHandler called", "method", r.Method)

	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	merchant := &models.Merchant{}
	if err := merchant.GetMerchant(ctx, userID); err != nil {
		logs.Logger.Error("GetMerchantHandler failed to get merchant", "error", err)
		http.Error(w, "Failed to get merchant", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(merchant)
}
