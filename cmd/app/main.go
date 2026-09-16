package main

import (
	"log"
	"os"

	"github.com/mathbdw/subscription-service/config"
	"github.com/mathbdw/subscription-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.ReadConfigYML("config.yml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		err := app.RunHealth(cfg)
		if err != nil {
			os.Exit(1)
		}
	} else {
		// Run app
		app.RunApp(cfg)
	}
}
