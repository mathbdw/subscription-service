package main

import (
	"log"

	"github.com/mathbdw/subscription-service/config"
	"github.com/mathbdw/subscription-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.ReadConfigYML("config.yml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run app
	app.RunApp(cfg)
}
