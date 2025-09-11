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

	fmt.Printf("\nConfiguration:\n%+v\n\n", cfg)
}
