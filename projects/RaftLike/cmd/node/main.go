package main

import (
	"RaftLike/internal/app"
	"log"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("error in creating app. %s", err)
	}
	app.Start()
}
