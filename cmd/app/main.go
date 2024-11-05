package main

import (
	"log"
	"os"

	"studentgit.kata.academy/movie-recommendation-platform/telegram-bot-service/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(a.Run())
}
