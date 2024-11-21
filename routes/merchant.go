package routes

import (
	"infinirewards/controllers"
	"infinirewards/middleware"
	"net/http"
)

func SetMerchantRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /merchant", middleware.AuthMiddleware(controllers.GetMerchantHandler))
}
