package main

import (
	"log"

	"github.com/andrsj/go-rabbit-image/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Error creating App object: %s", err)
	}

	app.Start()
	app.WaitForShutdown()
	if err := app.Stop(); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}
}
