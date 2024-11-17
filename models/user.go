package models

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"infinirewards/nats"
	"time"
)

type CreateUserRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
}

type UpdateUserRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
}

// User represents a user in the InfiniRewards system.
type User struct {
	ID             string    `json:"id"`
	PhoneNumber    string    `json:"phoneNumber"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	NKey           string    `json:"nKey"`
	NatsPublicKey  string    `json:"natsPublicKey"`
	PrivateKey     string    `json:"privateKey"`     // StarkNet private key
	AccountAddress string    `json:"accountAddress"` // StarkNet account address
}

const (
	usersBucket = "users"
)

// HashPhoneNumber generates a SHA-256 hash of the phone number
func HashPhoneNumber(phoneNumber string) string {
	hash := sha256.Sum256([]byte(phoneNumber))
	return fmt.Sprintf("0x%x", hash)
}

// CreateUser creates a new user in NATS KV Store
func (u *User) CreateUser(ctx context.Context) error {
	// Generate ID from phone number hash
	u.ID = HashPhoneNumber(u.PhoneNumber)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// Store user data
	userData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	if err := nats.PutKV(ctx, usersBucket, u.ID, userData); err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID from NATS KV Store
func (u *User) GetUser(ctx context.Context, id string) error {
	userKV, err := nats.GetKV(ctx, usersBucket, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return json.Unmarshal(userKV.Value(), u)
}

// GetUserFromPhoneNumber retrieves a user by phone number from NATS KV Store
func (u *User) GetUserFromPhoneNumber(ctx context.Context, phoneNumber string) error {
	// Generate ID from phone number hash
	id := HashPhoneNumber(phoneNumber)

	// Get user data directly using the hashed phone number as ID
	userKV, err := nats.GetKV(ctx, usersBucket, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return json.Unmarshal(userKV.Value(), u)
}

// UpdateUser updates an existing user in NATS KV Store
func (u *User) UpdateUser(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	userData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	if err := nats.PutKV(ctx, usersBucket, u.ID, userData); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from NATS KV Store
func (u *User) DeleteUser(ctx context.Context) error {

	// Delete user data
	if err := nats.RemoveKV(ctx, usersBucket, u.ID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *CreateUserRequest) Validate() error {
	if r.PhoneNumber == "" {
		return &ValidationError{
			Field:   "phoneNumber",
			Message: "phone number is required",
		}
	}
	if r.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name is required",
		}
	}
	// Email is optional but should be valid if provided
	if r.Email != "" {
		// Add email format validation if needed
	}
	return nil
}

func (r *UpdateUserRequest) Validate() error {
	// For update requests, we'll allow partial updates
	// but validate format of provided fields
	if r.PhoneNumber != "" {
		// Add phone number format validation if needed
	}
	if r.Email != "" {
		// Add email format validation if needed
	}
	if r.Name != "" && len(r.Name) > 100 {
		return &ValidationError{
			Field:   "name",
			Message: "name must be less than 100 characters",
		}
	}
	return nil
}
