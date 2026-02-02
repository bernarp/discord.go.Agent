package main

import (
	"log"
)

func main() {
	application, err := New()
	if err != nil {
		log.Fatalf("initialization failed: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application error: %v", err)
	}
}
