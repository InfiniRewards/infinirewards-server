package models

import (
	"context"
	"encoding/json"
	"fmt"
	"infinirewards/nats"
	"time"
)

type CreateMerchantRequest struct {
	// Name of the merchant
	// example: My Store
	Name string `json:"name" validate:"required,min=1,max=100"`

	// Symbol for the points token
	// example: PTS
	Symbol string `json:"symbol" validate:"required,len=3|len=4,uppercase"`

	// Decimals for the points token
	// example: 18
	Decimals uint8 `json:"decimals" validate:"required,gte=0,lte=18"`
}

type CreateMerchantResponse struct {
	// TransactionHash is the hash of the creation transaction
	// example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
	TransactionHash string `json:"transactionHash"`

	// MerchantAddress is the deployed merchant contract address
	// example: 0x1234567890abcdef1234567890abcdef12345678
	MerchantAddress string `json:"merchantAddress"`

	// PointsAddress is the deployed points contract address
	// example: 0x9876543210abcdef1234567890abcdef12345678
	PointsAddress string `json:"pointsAddress"`
}

type Merchant struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Decimals  uint8     `json:"decimals"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	merchantsBucket = "merchants"
)

// CreateMerchant creates a new merchant in NATS KV Store
func (m *Merchant) CreateMerchant(ctx context.Context, user *User) error {
	// Generate ID from phone number hash
	m.ID = HashPhoneNumber(user.PhoneNumber)
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	// Store merchant data
	merchantData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal merchant: %w", err)
	}

	// Check for existing merchant
	existingMerchant := &Merchant{}
	if err := existingMerchant.GetMerchant(ctx, m.ID); err == nil {
		return fmt.Errorf("merchant with phone number %s already exists", user.PhoneNumber)
	}

	if err := nats.PutKV(ctx, merchantsBucket, m.ID, merchantData); err != nil {
		return fmt.Errorf("failed to store merchant: %w", err)
	}

	return nil
}

// GetMerchant retrieves a merchant by ID from NATS KV Store
func (m *Merchant) GetMerchant(ctx context.Context, id string) error {
	merchantKV, err := nats.GetKV(ctx, merchantsBucket, id)
	if err != nil {
		return fmt.Errorf("failed to get merchant: %w", err)
	}

	return json.Unmarshal(merchantKV.Value(), m)
}

// GetMerchantFromPhoneNumber retrieves a merchant by phone number from NATS KV Store
func (m *Merchant) GetMerchantFromPhoneNumber(ctx context.Context, phoneNumber string) error {
	// Generate ID from phone number hash
	id := HashPhoneNumber(phoneNumber)

	// Get merchant data directly using the hashed phone number as ID
	merchantKV, err := nats.GetKV(ctx, merchantsBucket, id)
	if err != nil {
		return fmt.Errorf("failed to get merchant: %w", err)
	}

	return json.Unmarshal(merchantKV.Value(), m)
}

// UpdateMerchant updates an existing merchant in NATS KV Store
func (m *Merchant) UpdateMerchant(ctx context.Context) error {
	m.UpdatedAt = time.Now()
	merchantData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal merchant: %w", err)
	}

	if err := nats.PutKV(ctx, merchantsBucket, m.ID, merchantData); err != nil {
		return fmt.Errorf("failed to update merchant: %w", err)
	}

	return nil
}

// DeleteMerchant deletes a merchant from NATS KV Store
func (m *Merchant) DeleteMerchant(ctx context.Context) error {

	// Delete merchant data
	if err := nats.RemoveKV(ctx, merchantsBucket, m.ID); err != nil {
		return fmt.Errorf("failed to delete merchant: %w", err)
	}

	return nil
}

func (r *CreateMerchantRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name is required",
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
	return nil
}
