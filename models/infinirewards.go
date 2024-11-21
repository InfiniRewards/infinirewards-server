package models

// MintCollectibleRequest represents a request to mint collectible tokens
type MintCollectibleRequest struct {
	// CollectibleAddress is the contract address of the collectible
	// example: 0x1234567890abcdef1234567890abcdef12345678
	CollectibleAddress string `json:"collectibleAddress" validate:"required,eth_addr"`

	// To is the recipient address
	// example: 0x9876543210abcdef1234567890abcdef12345678
	To string `json:"to" validate:"required,eth_addr"`

	// TokenId is the ID of the token to mint
	// example: 1
	TokenId string `json:"tokenId" validate:"required,numeric"`

	// Amount is the number of tokens to mint
	// example: 1
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

// MintCollectibleResponse represents the response from minting collectible tokens
type MintCollectibleResponse struct {
	// TransactionHash is the hash of the mint transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

// CreateCollectibleRequest represents a request to create a new collectible contract
type CreateCollectibleRequest struct {
	// Name of the collectible contract
	// example: Special Edition
	Name string `json:"name" validate:"required,min=1,max=100"`

	// Description of the collectible collection
	// example: Limited edition collectibles
	Description string `json:"description" validate:"required,min=1,max=500"`
}

// CreateCollectibleResponse represents the response from creating a collectible contract
type CreateCollectibleResponse struct {
	// TransactionHash is the hash of the creation transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`

	// Address is the deployed contract address
	// example: 0x1234567890abcdef1234567890abcdef12345678
	Address string `json:"address"`
}

// GetCollectibleBalanceResponse represents the response for a balance query
type GetCollectibleBalanceResponse struct {
	// Balance is the number of tokens owned
	// example: 100
	Balance string `json:"balance"`
}

// GetTokenDataResponse represents the metadata for a collectible token
type GetTokenDataResponse struct {
	// PointsContract is the address of the points contract used for purchases
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract"`

	// Price is the cost in points to purchase the token
	// example: 100
	Price string `json:"price"`

	// Expiry is the Unix timestamp when the token expires
	// example: 1735689600
	Expiry int64 `json:"expiry"`

	// Description is the token's metadata description
	// example: Limited edition collectible
	Description string `json:"description"`
}

// RedeemCollectibleRequest represents a request to redeem a collectible
type RedeemCollectibleRequest struct {
	// CollectibleAddress is set from the URL path
	CollectibleAddress string `json:"-"`

	// User is the address redeeming the collectible
	// example: 0x1234567890abcdef1234567890abcdef12345678
	User string `json:"user" validate:"required,eth_addr"`

	// TokenId is the ID of the token to redeem
	// example: 1
	TokenId string `json:"tokenId" validate:"required,numeric"`

	// Amount is the number of tokens to redeem
	// example: 1
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

// RedeemCollectibleResponse represents the response from redeeming a collectible
type RedeemCollectibleResponse struct {
	// TransactionHash is the hash of the redemption transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

// Points-related models

// CreatePointsContractRequest represents a request to create a points contract
type CreatePointsContractRequest struct {
	// Name of the points token
	// example: Premium Points
	Name string `json:"name" validate:"required,min=1,max=100"`

	// Symbol for the points token (3-4 characters)
	// example: PPT
	Symbol string `json:"symbol" validate:"required,len=3|len=4,uppercase"`

	// Description of the points system
	// example: Premium tier loyalty points
	Description string `json:"description" validate:"required,min=1,max=500"`

	// Decimals specifies the number of decimal places
	// example: 18
	Decimals string `json:"decimals" validate:"required,numeric"`
}

// CreatePointsContractResponse represents the response from creating a points contract
type CreatePointsContractResponse struct {
	// TransactionHash is the hash of the creation transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`

	// Address is the deployed contract address
	// example: 0x1234567890abcdef1234567890abcdef12345678
	Address string `json:"address"`
}

// MintPointsRequest represents a request to mint points tokens
type MintPointsRequest struct {
	// PointsContract is the address of the points contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract" validate:"required,eth_addr"`

	// Recipient is the address receiving the points
	// example: 0x9876543210abcdef1234567890abcdef12345678
	Recipient string `json:"recipient" validate:"required,eth_addr"`

	// Amount of points to mint (in smallest unit)
	// example: 100
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

// MintPointsResponse represents the response from minting points
type MintPointsResponse struct {
	// TransactionHash is the hash of the mint transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

// BurnPointsRequest represents a request to burn points tokens
type BurnPointsRequest struct {
	// PointsContract is the address of the points contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract" validate:"required,eth_addr"`

	// Amount of points to burn
	// example: 50
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

// BurnPointsResponse represents the response from burning points
type BurnPointsResponse struct {
	// TransactionHash is the hash of the burn transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

// TransferPointsRequest represents a request to transfer points
type TransferPointsRequest struct {
	// PointsContract is the address of the points contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract" validate:"required,eth_addr"`

	// To is the recipient's address
	// example: 0x9876543210abcdef1234567890abcdef12345678
	To string `json:"to" validate:"required,eth_addr"`

	// Amount of points to transfer
	// example: 25
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

// TransferPointsResponse represents the response from transferring points
type TransferPointsResponse struct {
	// TransactionHash is the hash of the transfer transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

// Contract listing responses

// PointsContractInfo represents details of a points contract
type PointsContractInfo struct {
	// Address of the contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	Address string `json:"address"`

	// Name of the points token
	// example: Store Points
	Name string `json:"name"`

	// Symbol of the points token
	// example: SP
	Symbol string `json:"symbol"`

	// Description of the points token
	// example: Loyalty points for Store XYZ
	Description string `json:"description"`

	// Decimals places for the token
	// example: 18
	Decimals uint8 `json:"decimals"`

	// TotalSupply of the points token
	// example: 1000000
	TotalSupply uint64 `json:"totalSupply"`
}

// GetPointsContractsResponse represents the response listing points contracts
type GetPointsContractsResponse struct {
	// Contracts is the list of points contracts
	Contracts []PointsContractInfo `json:"contracts"`
}

// CollectibleContractInfo represents details of a collectible contract
type CollectibleContractInfo struct {
	// Address of the contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	Address string `json:"address"`

	// Name of the collectible collection
	// example: Special Edition
	Name string `json:"name"`

	// Description of the collection
	// example: Limited collectibles
	Description string `json:"description"`

	// PointsContract is the address of the points contract used for purchases
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract"`

	// TotalSupply across all token types
	// example: 100
	TotalSupply string `json:"totalSupply"`

	// TokenTypes lists all token IDs in the collection
	// example: ["1","2","3"]
	TokenTypes []string `json:"tokenTypes"`

	// TokenPrices lists prices for each token
	// example: ["100","200","300"]
	TokenPrices []string `json:"tokenPrices"`

	// TokenExpiries lists expiry timestamps for each token
	// example: [1735689600,1735689600,1735689600]
	TokenExpiries []uint64 `json:"tokenExpiries"`

	// TokenDescriptions lists descriptions for each token
	// example: ["Gold","Silver","Bronze"]
	TokenDescriptions []string `json:"tokenDescriptions"`
}

// GetCollectibleContractsResponse represents the response listing collectible contracts
type GetCollectibleContractsResponse struct {
	// Contracts is the list of collectible contracts
	Contracts []CollectibleContractInfo `json:"contracts"`
}

// Validation methods

func (r *MintCollectibleRequest) Validate() error {
	if r.CollectibleAddress == "" {
		return &ValidationError{
			Field:   "collectibleAddress",
			Message: "collectible address is required",
		}
	}
	if r.To == "" {
		return &ValidationError{
			Field:   "to",
			Message: "recipient address is required",
		}
	}
	if r.TokenId == "" {
		return &ValidationError{
			Field:   "tokenId",
			Message: "token ID is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

func (r *CreateCollectibleRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}
	if len(r.Name) > 100 {
		return &ValidationError{
			Field:   "name",
			Message: "name must be less than 100 characters",
		}
	}
	if r.Description == "" {
		return &ValidationError{
			Field:   "description",
			Message: "description is required",
		}
	}
	if len(r.Description) > 500 {
		return &ValidationError{
			Field:   "description",
			Message: "description must be less than 500 characters",
		}
	}
	return nil
}

func (r *CreatePointsContractRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}
	if len(r.Name) > 100 {
		return &ValidationError{
			Field:   "name",
			Message: "name must be less than 100 characters",
		}
	}
	if r.Symbol == "" {
		return &ValidationError{
			Field:   "symbol",
			Message: "symbol is required",
		}
	}
	if len(r.Symbol) < 3 || len(r.Symbol) > 4 {
		return &ValidationError{
			Field:   "symbol",
			Message: "symbol must be 3-4 characters",
		}
	}
	if r.Description == "" {
		return &ValidationError{
			Field:   "description",
			Message: "description is required",
		}
	}
	if len(r.Description) > 500 {
		return &ValidationError{
			Field:   "description",
			Message: "description must be less than 500 characters",
		}
	}
	if r.Decimals == "" {
		return &ValidationError{
			Field:   "decimals",
			Message: "decimals is required",
		}
	}
	return nil
}

func (r *MintPointsRequest) Validate() error {
	if r.PointsContract == "" {
		return &ValidationError{
			Field:   "pointsContract",
			Message: "points contract address is required",
		}
	}
	if r.Recipient == "" {
		return &ValidationError{
			Field:   "recipient",
			Message: "recipient address is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

func (r *BurnPointsRequest) Validate() error {
	if r.PointsContract == "" {
		return &ValidationError{
			Field:   "pointsContract",
			Message: "points contract address is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

func (r *TransferPointsRequest) Validate() error {
	if r.PointsContract == "" {
		return &ValidationError{
			Field:   "pointsContract",
			Message: "points contract address is required",
		}
	}
	if r.To == "" {
		return &ValidationError{
			Field:   "to",
			Message: "recipient address is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

func (r *RedeemCollectibleRequest) Validate() error {
	if r.CollectibleAddress == "" {
		return &ValidationError{
			Field:   "collectibleAddress",
			Message: "collectible address is required",
		}
	}
	if r.User == "" {
		return &ValidationError{
			Field:   "user",
			Message: "user address is required",
		}
	}
	if r.TokenId == "" {
		return &ValidationError{
			Field:   "tokenId",
			Message: "token ID is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

// Additional models needed based on linter errors

type GetCollectibleURIResponse struct {
	// URI of the collectible metadata
	// example: https://example.com/metadata/1
	URI string `json:"uri"`
}

type SetTokenDataRequest struct {
	// CollectibleAddress is set from the URL path
	CollectibleAddress string `json:"-"`

	// TokenId is set from the URL path
	TokenId string `json:"-"`

	// PointsContract is the address of the points contract used for purchases
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract" validate:"required,eth_addr"`

	// Price is the cost in points to purchase the token
	// example: 100
	Price string `json:"price" validate:"required,numeric,gt=0"`

	// Expiry is the Unix timestamp when the token expires
	// example: 1735689600
	Expiry uint64 `json:"expiry" validate:"required,gt=0"`

	// Description is the token's metadata description
	// example: Limited edition collectible
	Description string `json:"description" validate:"required,min=1,max=500"`
}

type SetTokenDataResponse struct {
	// TransactionHash is the hash of the set data transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

type GetCollectibleDetailsResponse struct {
	// Description of the collectible collection
	// example: Special Edition Collectibles
	Description string `json:"description"`

	// PointsContract is the address of the points contract
	// example: 0x1234567890abcdef1234567890abcdef12345678
	PointsContract string `json:"pointsContract"`

	// TokenIDs lists all token IDs in the collection
	// example: ["1","2","3"]
	TokenIDs []string `json:"tokenIDs"`

	// TokenPrices lists prices for each token
	// example: ["100","200","300"]
	TokenPrices []string `json:"tokenPrices"`

	// TokenExpiries lists expiry timestamps for each token
	// example: [1735689600,1735689600,1735689600]
	TokenExpiries []uint64 `json:"tokenExpiries"`

	// TokenDescriptions lists descriptions for each token
	// example: ["Gold","Silver","Bronze"]
	TokenDescriptions []string `json:"tokenDescriptions"`
}

type IsCollectibleValidResponse struct {
	// IsValid indicates if the collectible is still valid
	// example: true
	IsValid bool `json:"isValid"`
}

type PurchaseCollectibleRequest struct {
	// CollectibleAddress is set from the URL path
	CollectibleAddress string `json:"-"`

	// User is the address purchasing the collectible
	// example: 0x1234567890abcdef1234567890abcdef12345678
	User string `json:"user" validate:"required,eth_addr"`

	// TokenId is the ID of the token to purchase
	// example: 1
	TokenId string `json:"tokenId" validate:"required,numeric"`

	// Amount is the number of tokens to purchase
	// example: 1
	Amount string `json:"amount" validate:"required,numeric,gt=0"`
}

type PurchaseCollectibleResponse struct {
	// TransactionHash is the hash of the purchase transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`
}

func (r *SetTokenDataRequest) Validate() error {
	if r.CollectibleAddress == "" {
		return &ValidationError{
			Field:   "collectibleAddress",
			Message: "collectible address is required",
		}
	}
	if r.TokenId == "" {
		return &ValidationError{
			Field:   "tokenId",
			Message: "token ID is required",
		}
	}
	if r.PointsContract == "" {
		return &ValidationError{
			Field:   "pointsContract",
			Message: "points contract address is required",
		}
	}
	if r.Price == "" {
		return &ValidationError{
			Field:   "price",
			Message: "price is required",
		}
	}
	if r.Expiry == 0 {
		return &ValidationError{
			Field:   "expiry",
			Message: "expiry timestamp is required",
		}
	}
	if r.Description == "" {
		return &ValidationError{
			Field:   "description",
			Message: "description is required",
		}
	}
	return nil
}

func (r *PurchaseCollectibleRequest) Validate() error {
	if r.CollectibleAddress == "" {
		return &ValidationError{
			Field:   "collectibleAddress",
			Message: "collectible address is required",
		}
	}
	if r.User == "" {
		return &ValidationError{
			Field:   "user",
			Message: "user address is required",
		}
	}
	if r.TokenId == "" {
		return &ValidationError{
			Field:   "tokenId",
			Message: "token ID is required",
		}
	}
	if r.Amount == "" {
		return &ValidationError{
			Field:   "amount",
			Message: "amount is required",
		}
	}
	return nil
}

// GetPointsBalanceResponse represents the response for a points balance query
type GetPointsBalanceResponse struct {
	// Balance is the number of points owned
	// example: 1000
	Balance string `json:"balance"`
}
