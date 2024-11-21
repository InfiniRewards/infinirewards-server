package main

import (
	"context"
	_ "infinirewards/docs" // This line is necessary for swagger
	"infinirewards/infinirewards"
	"infinirewards/jwt"
	"infinirewards/logs"
	"infinirewards/routes"
	"infinirewards/utils"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			InfiniRewards API
//	@version		1.0
//	@description	API for InfiniRewards - A Web3 Loyalty and Rewards Platform

//	@contact.name	API Support
//	@contact.url	http://www.infinirewards.io/support
//	@contact.email	support@infinirewards.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/
//	@schemes	http https

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer <token>**

//	@x-extension-openapi	{"example": "value on a json format"}

func main() {
	// Initialize logger
	logs.InitHandler("")

	logs.Logger.Info("starting server",
		slog.String("handler", "main"),
	)

	// Initialize JWT keys
	if err := jwt.InitKeys(); err != nil {
		logs.Logger.Error("failed to initialize JWT keys",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Initialize Starknet connection
	if err := infinirewards.ConnectStarknet(); err != nil {
		logs.Logger.Error("failed to connect to Starknet",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Initialize WhatsApp
	if err := utils.InitWhatsApp(); err != nil {
		logs.Logger.Error("failed to initialize WhatsApp",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	// Initialize MacroKiosk
	if err := utils.InitMacroKiosk(); err != nil {
		logs.Logger.Error("failed to initialize MacroKiosk",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create new ServeMux
	mux := http.NewServeMux()

	// Register routes
	routes.SetAuthRoutes(mux)
	routes.SetUserRoutes(mux)
	routes.SetMerchantRoutes(mux)
	routes.SetInfiniRewardsRoutes(mux)

	// Only serve Swagger docs in development/staging environments
	if os.Getenv("ENV") != "production" {
		mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))
	}

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		logs.Logger.Info("server started",
			slog.String("handler", "main"),
			slog.String("addr", server.Addr),
		)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logs.Logger.Error("server error",
				slog.String("handler", "main"),
				slog.String("error", err.Error()),
			)
		}
	}()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	logs.Logger.Info("shutting down server",
		slog.String("handler", "main"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logs.Logger.Error("error stopping server",
			slog.String("handler", "main"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logs.Logger.Info("server stopped gracefully",
		slog.String("handler", "main"),
	)
}
