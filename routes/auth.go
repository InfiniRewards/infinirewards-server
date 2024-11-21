package routes

import (
	"encoding/json"
	"infinirewards/controllers"
	"infinirewards/jwt"
	"infinirewards/middleware"
	"net/http"
)

//	@title			InfiniRewards API
//	@version		1.0
//	@description	InfiniRewards service API for managing rewards and collectibles
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.infinirewards.io/support
//	@contact.email	support@infinirewards.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/
//	@schemes	http https

func SetAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/request-otp", controllers.RequestOTPHandler)
	mux.HandleFunc("POST /auth/authenticate", controllers.AuthenticateHandler)
	mux.HandleFunc("POST /auth/refresh-token", middleware.AuthMiddleware(controllers.RefreshTokenHandler))

	// Add JWKS endpoint for future use
	mux.HandleFunc("GET /.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		jwks, err := jwt.GetJWKS()
		if err != nil {
			http.Error(w, "Failed to get JWKS", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	})
}
