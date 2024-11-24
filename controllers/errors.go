package controllers

import (
	"encoding/json"
	"infinirewards/models"
	"net/http"
)

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, message string, code string, details interface{}, statusCode int) {
	errResp := models.ErrorResponse{
		Message: message,
		Code:    code,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errResp)
}

// Common error codes
const (
	ValidationError     = "VALIDATION_ERROR"
	AuthenticationError = "AUTHENTICATION_ERROR"
	AuthorizationError  = "AUTHORIZATION_ERROR"
	NotFoundError       = "NOT_FOUND"
	ConflictError       = "CONFLICT"
	RateLimitError      = "RATE_LIMIT_EXCEEDED"
	InternalServerError = "INTERNAL_ERROR"
)
