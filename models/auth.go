package models

import "time"

type Token struct {
	ID          string `json:"id"`
	User        string `json:"user"`
	AccessToken string `json:"accessToken"`
	// RefreshToken       string    `json:"refreshToken"` (use ID as refresh token)
	AccessTokenExpiry  time.Time `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry"`
	Device             string    `json:"device"`
	Service            string    `json:"service"`
	CreatedAt          time.Time `json:"createdAt"`
}

type PhoneNumberVerification struct {
	ID        string    `json:"id"`
	Secret    string    `json:"secret"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

type RequestOTPRequest struct {
	// PhoneNumber must be in E.164 format (e.g. +60123456789)
	// example: +60123456789
	PhoneNumber string `json:"phoneNumber" validate:"required,e164"`
}

type RequestOTPResponse struct {
	Message string `json:"message" example:"message"`
	ID      string `json:"id" example:"phoneNumberVerification:01HNAJ6640M9JRRJFQSZZVE3HH"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	Token Token  `json:"token"`
	Creds string `json:"creds"`
}

type AuthenticateRequest struct {
	// ID is either a verification ID (for OTP) or API key ID
	// example: 01HNAJ6640M9JRRJFQSZZVE3HH
	ID string `json:"id" validate:"required"`

	// Token is either an OTP code or API key secret
	// example: 123456
	Token string `json:"token" validate:"required"`

	// Method must be either "otp" or "secret"
	// example: otp
	Method string `json:"method" validate:"required,oneof=otp secret"`

	// Signature is a unique device identifier
	// example: device_signature_123
	Signature string `json:"signature" validate:"required"`

	// Device is a unique device identifier
	// example: device_123
	Device string `json:"device" validate:"required"`
}

type AuthenticateResponse struct {
	Token Token `json:"token"`
}

func (r *RequestOTPRequest) Validate() error {
	if r.PhoneNumber == "" {
		return &ValidationError{
			Field:   "phoneNumber",
			Message: "phone number is required",
		}
	}
	// Add phone number format validation if needed
	return nil
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return &ValidationError{
			Field:   "refreshToken",
			Message: "refresh token is required",
		}
	}
	return nil
}

func (r *AuthenticateRequest) Validate() error {
	if r.ID == "" {
		return &ValidationError{
			Field:   "id",
			Message: "id is required",
		}
	}
	if r.Token == "" {
		return &ValidationError{
			Field:   "token",
			Message: "token is required",
		}
	}
	if r.Method == "" {
		return &ValidationError{
			Field:   "method",
			Message: "method is required",
		}
	}
	if r.Signature == "" {
		return &ValidationError{
			Field:   "signature",
			Message: "signature is required",
		}
	}
	return nil
}
