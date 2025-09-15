package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"qserve/internal/config"
	"qserve/internal/database"
	"qserve/internal/handler"
	"qserve/internal/middleware"
)

func main() {
	logger := middleware.NewLogger()
	fmt.Println("🔧 Welcome to qserve setup!")
	fmt.Println("============================")

	cfg, err := config.RunNewSetupWizard()
	if err != nil {
		log.Fatalf("Setup failed: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	dbManager := database.NewConnectionManager(cfg)
	defer dbManager.Close()

	ctx := context.Background()

	if err := dbManager.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connected successfully")

	if err := dbManager.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	validator := database.NewQueryValidator()
	queryHandler := handler.NewQueryHandler(dbManager, validator)

	corsMiddleware := middleware.CorsMiddleware
	loggerMiddleware := middleware.NewLoggingMiddleware(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /query", queryHandler.HandleQuery)
	mux.HandleFunc("GET /health", queryHandler.HandleHealthCheck)

	handler := loggerMiddleware.LoggerMiddleware(corsMiddleware(mux))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	logger.Info("Server started", "port", cfg.Port)
	logger.Info("Available endpoints:",
		"POST /query", "Execute SQL queries",
		"GET /health", "Health check",
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	logger.Info("Server stopped")
}
