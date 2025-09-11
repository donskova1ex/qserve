package main

import (
	"context"
	"fmt"
	"log"
	"qserve/internal/config"
	"qserve/internal/database"
)

func main() {
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
	fmt.Println("Database connected successfully")
	if err := dbManager.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}
