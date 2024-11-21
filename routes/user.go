package routes

import (
	"infinirewards/controllers"
	"infinirewards/middleware"
	"net/http"
)

//	@title			InfiniRewards API
//	@version		1.0
//	@description	InfiniRewards service API for managing rewards and collectibles

func SetUserRoutes(mux *http.ServeMux) {
	//	@Summary		Get User
	//	@Description	Get authenticated user details
	//	@Tags			user
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{object}	models.User
	//	@Failure		400	{string}	string	"Bad Request"
	//	@Failure		401	{string}	string	"Unauthorized"
	//	@Failure		404	{string}	string	"User not found"
	//	@Failure		500	{string}	string	"Internal Server Error"
	//	@Router			/user [get]
	mux.HandleFunc("GET /user", middleware.AuthMiddleware(controllers.UserGetUserHandler))

	//	@Summary		Create User
	//	@Description	Create a new user
	//	@Tags			user
	//	@Accept			json
	//	@Produce		json
	//	@Param			request	body		models.CreateUserRequest	true	"User Creation Request"
	//	@Success		201		{object}	models.User
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/user [post]
	mux.HandleFunc("POST /user", middleware.AuthMiddleware(controllers.UserCreateUserHandler))

	//	@Summary		Update User
	//	@Description	Update authenticated user details
	//	@Tags			user
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.UpdateUserRequest	true	"User Update Request"
	//	@Success		200		{object}	models.User
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		404		{string}	string	"User not found"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/user [put]
	mux.HandleFunc("PUT /user", middleware.AuthMiddleware(controllers.UserUpdateUserHandler))

	//	@Summary		Delete User
	//	@Description	Delete authenticated user
	//	@Tags			user
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{object}	models.User
	//	@Failure		401	{string}	string	"Unauthorized"
	//	@Failure		404	{string}	string	"User not found"
	//	@Failure		500	{string}	string	"Internal Server Error"
	//	@Router			/user [delete]
	mux.HandleFunc("DELETE /user", middleware.AuthMiddleware(controllers.UserDeleteUserHandler))

	// API Key routes
	//	@Summary		Create API Key
	//	@Description	Create a new API key for authenticated user
	//	@Tags			api-keys
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.CreateAPIKeyRequest	true	"API Key Creation Request"
	//	@Success		201		{object}	models.APIKey
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/user/api-keys [post]
	mux.HandleFunc("POST /user/api-keys", middleware.AuthMiddleware(controllers.UserCreateAPIKeyHandler))

	//	@Summary		List API Keys
	//	@Description	List all API keys for authenticated user
	//	@Tags			api-keys
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{array}		models.APIKey
	//	@Failure		401	{string}	string	"Unauthorized"
	//	@Failure		500	{string}	string	"Internal Server Error"
	//	@Router			/user/api-keys [get]
	mux.HandleFunc("GET /user/api-keys", middleware.AuthMiddleware(controllers.UserListAPIKeysHandler))

	//	@Summary		Delete API Key
	//	@Description	Delete an API key
	//	@Tags			api-keys
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			keyId	path		string	true	"API Key ID"
	//	@Success		200		{object}	models.MessageResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/user/api-keys/{keyId} [delete]
	mux.HandleFunc("DELETE /user/api-keys/{keyId}", middleware.AuthMiddleware(controllers.UserDeleteAPIKeyHandler))
}
