package routes

import (
	"infinirewards/controllers"
	"infinirewards/middleware"
	"net/http"
)

// @title InfiniRewards API
// @version 1.0
// @description InfiniRewards service API for managing rewards and collectibles
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.infinirewards.io/support
// @contact.email support@infinirewards.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

func SetAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/request-otp", controllers.RequestOTPHandler)

	// @Summary Authenticate
	// @Description Authenticate user using OTP or API key
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param request body models.AuthenticateRequest true "Authentication Request"
	// @Success 200 {object} models.AuthenticateResponse
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /auth/authenticate [post]
	mux.HandleFunc("POST /auth/authenticate", controllers.AuthenticateHandler)

	// @Summary Refresh Token
	// @Description Refresh an existing authentication token
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body models.RefreshTokenRequest true "Token Refresh Request"
	// @Success 200 {object} models.RefreshTokenResponse
	// @Failure 400 {string} string "Bad Request"
	// @Failure 401 {string} string "Unauthorized"
	// @Failure 500 {string} string "Internal Server Error"
	// @Router /auth/refresh-token [post]
	mux.HandleFunc("POST /auth/refresh-token", middleware.AuthMiddleware(controllers.RefreshTokenHandler))
}
