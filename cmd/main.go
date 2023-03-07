package main

import (
	"github.com/andrsj/go-rabbit-image/internal/app"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

func main() {
	var log logger.Logger = logger.NewLogrusLogger("debug")
	log = log.Named("main")

	app, err := app.New(log)
	if err != nil {
		log.Fatal("Error creating App object", logger.M{
			"error": err,
		})
	}

	app.Start()
	app.WaitForShutdown()

	if err := app.Stop(); err != nil {
		log.Fatal("Server shutdown error", logger.M{
			"error": err,
		})
	}

}
