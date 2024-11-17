package models

import (
	"context"
	"encoding/json"
	"fmt"
	"infinirewards/nats"
	"infinirewards/utils"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
)

const KV_BUCKET = "apikeys"

type APIKey struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Secret    string    `json:"secret"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type DeleteAPIKeyRequest struct {
	ID string `json:"id"`
}

// CreateAPIKey creates a new API key for a user
func CreateAPIKey(ctx context.Context, userID string, name string) (*APIKey, error) {
	secret, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	apiKey := &APIKey{
		ID:        ulid.Make().String(),
		UserID:    userID,
		Name:      name,
		Secret:    secret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the API key using hierarchical subject
	subject := fmt.Sprintf("apikeys.%s.%s", userID, apiKey.ID)
	data, err := json.Marshal(apiKey)
	if err != nil {
		return nil, err
	}

	if err := nats.PutKV(ctx, KV_BUCKET, subject, data); err != nil {
		return nil, err
	}

	return apiKey, nil
}

// ListAPIKeys lists all API keys for a user
func ListAPIKeys(ctx context.Context, userID string) ([]*APIKey, error) {
	prefix := fmt.Sprintf("apikeys.%s.", userID)
	transformFunc := func(entry jetstream.KeyValueEntry, apiKey *APIKey) {
		apiKey.Secret = "" // Remove secret before returning
	}

	apiKeys, err := nats.GetKVValues[APIKey](ctx, KV_BUCKET, prefix, transformFunc)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// DeleteAPIKey deletes an API key
func DeleteAPIKey(ctx context.Context, userID string, keyID string) error {
	subject := fmt.Sprintf("apikeys.%s.%s", userID, keyID)

	// Get the key first to verify ownership
	entry, err := nats.GetKV(ctx, KV_BUCKET, subject)
	if err != nil {
		return fmt.Errorf("API key not found")
	}

	var apiKey APIKey
	if err := json.Unmarshal(entry.Value(), &apiKey); err != nil {
		return err
	}

	if apiKey.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return nats.RemoveKV(ctx, KV_BUCKET, subject)
}

// ValidateAPIKey validates an API key and secret
func ValidateAPIKey(ctx context.Context, userID, keyID, secret string) (*APIKey, error) {
	subject := fmt.Sprintf("apikeys.%s.%s", userID, keyID)
	entry, err := nats.GetKV(ctx, KV_BUCKET, subject)
	if err != nil {
		return nil, fmt.Errorf("API key not found: %w", err)
	}

	var apiKey APIKey
	if err := json.Unmarshal(entry.Value(), &apiKey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API key: %w", err)
	}

	if apiKey.Secret != secret {
		return nil, fmt.Errorf("invalid API key secret")
	}

	return &apiKey, nil
}

func (r *CreateAPIKeyRequest) Validate() error {
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
	return nil
}

func (r *DeleteAPIKeyRequest) Validate() error {
	if r.ID == "" {
		return &ValidationError{
			Field:   "id",
			Message: "id is required",
		}
	}
	return nil
}
