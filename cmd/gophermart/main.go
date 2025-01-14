package main

import (
	"log"

	"github.com/Makovey/gophermart/internal/app"
)

func main() {
	appl := app.NewApp()
	if err := appl.InitDependencies(); err != nil {
		log.Fatalf("critical dependencies initialized with error: %v", err.Error())
	}
	appl.Run()
}
