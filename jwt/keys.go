package jwt

import (
	"crypto/ed25519"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type KeyPair struct {
	ID        string
	Private   ed25519.PrivateKey
	Public    ed25519.PublicKey
	CreatedAt time.Time
}

var (
	activeKeys     = make(map[string]*KeyPair)
	keyMutex       sync.RWMutex
	privateKeyDir  = ".private"
	publicKeyDir   = ".public"
	keyRotationAge = 24 * time.Hour
)

// SetKeyDirectories allows setting test directories
func SetKeyDirectories(privateDir, publicDir string) {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	privateKeyDir = privateDir
	publicKeyDir = publicDir
}

// InitKeys initializes the JWT keys and starts key rotation
func InitKeys() error {
	keyMutex.Lock()
	defer keyMutex.Unlock()
	fmt.Printf("Initializing JWT keys")

	// Clear existing keys
	activeKeys = make(map[string]*KeyPair)

	// Create directories if they don't exist
	if err := os.MkdirAll(privateKeyDir, 0700); err != nil {
		return fmt.Errorf("failed to create private key directory: %w", err)
	}
	if err := os.MkdirAll(publicKeyDir, 0755); err != nil {
		return fmt.Errorf("failed to create public key directory: %w", err)
	}

	fmt.Printf("Private key directory: %s\n", privateKeyDir)
	fmt.Printf("Public key directory: %s\n", publicKeyDir)

	// Load existing keys or generate initial keys
	if err := loadOrGenerateKeys(); err != nil {
		return err
	}

	fmt.Printf("Loaded %d keys\n", len(activeKeys))

	// Start key rotation
	go rotateKeys()

	return nil
}

// loadOrGenerateKeys loads existing keys or generates initial keys if none exist
func loadOrGenerateKeys() error {

	files, err := os.ReadDir(privateKeyDir)
	if err != nil {
		return fmt.Errorf("failed to read private key directory: %w", err)
	}

	if len(files) == 0 {
		// Generate initial keys
		for i := 0; i < 3; i++ {
			if err := generateNewKeyPair(); err != nil {
				return fmt.Errorf("failed to generate initial key pair: %w", err)
			}
		}
		return nil
	}

	// Load existing keys
	for _, file := range files {
		id := file.Name()
		info, err := file.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", id, err)
		}

		privPath := filepath.Join(privateKeyDir, id)
		privBytes, err := os.ReadFile(privPath)
		if err != nil {
			return fmt.Errorf("failed to read private key %s: %w", id, err)
		}

		pubPath := filepath.Join(publicKeyDir, id+".pub")
		pubBytes, err := os.ReadFile(pubPath)
		if err != nil {
			return fmt.Errorf("failed to read public key %s: %w", id, err)
		}

		pair := &KeyPair{
			ID:        id,
			Private:   ed25519.PrivateKey(privBytes),
			Public:    ed25519.PublicKey(pubBytes),
			CreatedAt: info.ModTime(),
		}

		activeKeys[id] = pair
	}

	return nil
}

// rotateKeys handles key rotation every 24 hours
func rotateKeys() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := performKeyRotation(); err != nil {
			slog.Error("Failed to rotate keys",
				slog.String("error", err.Error()),
			)
		}
	}
}

// performKeyRotation performs the actual key rotation
func performKeyRotation() error {
	keyMutex.Lock()
	defer keyMutex.Unlock()

	// Check for expired keys
	now := time.Now()
	var expiredKeys []string
	activeCount := 0

	for id, key := range activeKeys {
		if now.Sub(key.CreatedAt) >= keyRotationAge {
			expiredKeys = append(expiredKeys, id)
		} else {
			activeCount++
		}
	}

	// Generate new keys if needed
	if activeCount < 3 {
		for i := 0; i < 3-activeCount; i++ {
			if err := generateNewKeyPair(); err != nil {
				return fmt.Errorf("failed to generate new key pair: %w", err)
			}
		}
	}

	// Remove expired keys
	for _, id := range expiredKeys {
		delete(activeKeys, id)
		os.Remove(filepath.Join(privateKeyDir, id))
		os.Remove(filepath.Join(publicKeyDir, id+".pub"))
	}

	return nil
}

// GetKeyByID returns a key pair by ID
func GetKeyByID(kid string) (*KeyPair, bool) {
	key, ok := activeKeys[kid]
	return key, ok
}

// GetJWKS returns the JWKS for public key verification
func GetJWKS() (*jose.JSONWebKeySet, error) {
	var keys []jose.JSONWebKey
	for _, pair := range activeKeys {
		if time.Since(pair.CreatedAt) < keyRotationAge {
			key := jose.JSONWebKey{
				Key:       pair.Public,
				KeyID:     pair.ID,
				Algorithm: string(jose.EdDSA),
				Use:       "sig",
			}
			keys = append(keys, key)
		}
	}
	return &jose.JSONWebKeySet{Keys: keys}, nil
}

// CreateToken creates a new JWT token
func CreateToken(userID string) (string, error) {
	// Get random key
	key := GetRandomActiveKey()
	if key == nil {
		return "", fmt.Errorf("no valid signing keys available")
	}

	// Create signer
	sig, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.EdDSA,
			Key:       key.Private,
		},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", key.ID),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %w", err)
	}

	// Create claims
	cl := jwt.Claims{
		Subject:   userID,
		Issuer:    "infinirewards",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Expiry:    jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  jwt.Audience{"infinirewards"},
	}

	// Sign token
	return jwt.Signed(sig).Claims(cl).CompactSerialize()
}

// VerifyToken verifies a JWT token
func VerifyToken(tokenString string) (*jwt.Claims, error) {
	tok, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Get key ID from header
	if len(tok.Headers) == 0 {
		return nil, fmt.Errorf("no headers in token")
	}
	kid := tok.Headers[0].KeyID

	// Get key
	key, ok := GetKeyByID(kid)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", kid)
	}

	// Verify claims
	var claims jwt.Claims
	if err := tok.Claims(key.Public, &claims); err != nil {
		return nil, fmt.Errorf("failed to verify claims: %w", err)
	}

	return &claims, nil
}

// GetRandomActiveKey returns a random active key pair
func GetRandomActiveKey() *KeyPair {
	var validKeys []*KeyPair
	for _, key := range activeKeys {
		if time.Since(key.CreatedAt) < keyRotationAge {
			validKeys = append(validKeys, key)
		}
	}

	if len(validKeys) == 0 {
		return nil
	}

	return validKeys[rand.Intn(len(validKeys))]
}

// generateNewKeyPair generates a new key pair
func generateNewKeyPair() error {
	pub, priv, err := ed25519.GenerateKey(crand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate EdDSA key pair: %w", err)
	}

	id := generateKeyID()
	pair := &KeyPair{
		ID:        id,
		Private:   priv,
		Public:    pub,
		CreatedAt: time.Now(),
	}

	// Save private key
	privPath := filepath.Join(privateKeyDir, id)
	if err := os.WriteFile(privPath, priv, 0600); err != nil {
		return fmt.Errorf("failed to save private key: %w", err)
	}

	// Save public key
	pubPath := filepath.Join(publicKeyDir, id+".pub")
	if err := os.WriteFile(pubPath, pub, 0644); err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	activeKeys[id] = pair
	return nil
}

// generateKeyID generates a new key ID
func generateKeyID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
