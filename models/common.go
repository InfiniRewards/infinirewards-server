package models

import (
	"github.com/invopop/jsonschema"
)

type MessageResponse struct {
	Message string `json:"message" example:"message"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func SchemaFor(t any) string {
	schema := jsonschema.Reflect(t)
	data, _ := schema.MarshalJSON()
	return string(data)
}
