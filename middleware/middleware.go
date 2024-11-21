package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"infinirewards/jwt"
	"infinirewards/logs"
)

// Define a custom type for context keys
type contextKey struct {
	name string
}

// Define the user ID key as a variable
var userIDKey = contextKey{"user_id"}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the JWT token using the keys package
		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			logs.Logger.Error("Failed to verify token",
				slog.String("error", err.Error()),
			)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the user ID from claims to the context using the custom key
		ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
