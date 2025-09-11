package main

import (
	"fmt"
	"log"
	"qserve/internal/config"
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
	fmt.Println("Configuration is valid")

	fmt.Printf("\nConfiguration:\n%+v\n\n", cfg)
}
