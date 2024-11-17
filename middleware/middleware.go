package middleware

import (
	"context"
	"fmt"
	"infinirewards/logs"
	"infinirewards/models"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type CustomClaims struct {
	jwt.Claims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

var (
	signingKey = []byte("your-secret-key") // In production, this should be properly secured
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		token, err := jwt.ParseSigned(tokenString)
		if err != nil {
			logs.Logger.Error("failed to parse token",
				slog.String("handler", "AuthMiddleware"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		var claims CustomClaims
		if err := token.Claims(signingKey, &claims); err != nil {
			logs.Logger.Error("failed to verify token",
				slog.String("handler", "AuthMiddleware"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "invalid token signature", http.StatusUnauthorized)
			return
		}

		// Validate expiry
		if err := claims.ValidateWithLeeway(jwt.Expected{
			Time: time.Now(),
		}, time.Minute); err != nil {
			logs.Logger.Error("token validation failed",
				slog.String("handler", "AuthMiddleware"),
				slog.String("error", err.Error()),
			)
			http.Error(w, "token expired or invalid", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(user models.User) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       signingKey,
	}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", fmt.Errorf("failed to create signer: %w", err)
	}

	cl := CustomClaims{
		Claims: jwt.Claims{
			Subject:   user.ID,
			Issuer:    "infinirewards",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        user.ID,
		},
		UserID: user.ID,
		Role:   "user", // You can set different roles as needed
	}

	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return raw, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseSigned(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims := &CustomClaims{}
	if err := token.Claims(signingKey, claims); err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if err := claims.ValidateWithLeeway(jwt.Expected{
		Time: time.Now(),
	}, time.Minute); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return claims, nil
}

// Helper function to extract user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// Helper function to extract role from context
func GetRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return "", fmt.Errorf("role not found in context")
	}
	return role, nil
}
