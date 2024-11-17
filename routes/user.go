package routes

import (
	"infinirewards/controllers"
	"infinirewards/middleware"
	"net/http"
)

// @title InfiniRewards API
// @version 1.0
// @description InfiniRewards service API for managing rewards and collectibles

func SetUserRoutes(mux *http.ServeMux) {
	// @Summary Get User
	// @Description Get user details by ID
	// @Tags users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Success 200 {object} models.User
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 404 {string} string "User not found"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id} [get]
	mux.HandleFunc("GET /users/{id}", middleware.AuthMiddleware(controllers.UserGetUserHandler))

	// @Summary Create User
	// @Description Create a new user
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param request body models.CreateUserRequest true "User Creation Request"
	// @Success 201 {object} models.User
	// @Failure 400 {string} string "Bad Request"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users [post]
	mux.HandleFunc("POST /users", controllers.UserCreateUserHandler)

	// @Summary Update User
	// @Description Update user details
	// @Tags users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Param request body models.UpdateUserRequest true "User Update Request"
	// @Success 200 {object} models.User
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 404 {string} string "User not found"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id} [put]
	mux.HandleFunc("PUT /users/{id}", middleware.AuthMiddleware(controllers.UserUpdateUserHandler))

	// @Summary Delete User
	// @Description Delete a user
	// @Tags users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Success 200 {object} models.User
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 404 {string} string "User not found"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id} [delete]
	mux.HandleFunc("DELETE /users/{id}", middleware.AuthMiddleware(controllers.UserDeleteUserHandler))

	// API Key routes
	// @Summary Create API Key
	// @Description Create a new API key for a user
	// @Tags api-keys
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Param request body models.CreateAPIKeyRequest true "API Key Creation Request"
	// @Success 201 {object} models.APIKey
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id}/api-keys [post]
	mux.HandleFunc("POST /users/{id}/api-keys", middleware.AuthMiddleware(controllers.UserCreateAPIKeyHandler))

	// @Summary List API Keys
	// @Description List all API keys for a user
	// @Tags api-keys
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Success 200 {array} models.APIKey
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id}/api-keys [get]
	mux.HandleFunc("GET /users/{id}/api-keys", middleware.AuthMiddleware(controllers.UserListAPIKeysHandler))

	// @Summary Delete API Key
	// @Description Delete an API key
	// @Tags api-keys
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID"
	// @Param keyId path string true "API Key ID"
	// @Success 200 {object} models.MessageResponse
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /users/{id}/api-keys/{keyId} [delete]
	mux.HandleFunc("DELETE /users/{id}/api-keys/{keyId}", middleware.AuthMiddleware(controllers.UserDeleteAPIKeyHandler))
}
