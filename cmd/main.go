package main

import (
	"log"

	_ "DiscordBotAgent/docs"
)

// @title Discord Bot Agent API
// @version 1.0
// @description Control plane for Discord Bot Agent modules and configuration.
// @host localhost:8001
// @BasePath /
func main() {
	application, err := New()
	if err != nil {
		log.Fatalf("initialization failed: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
