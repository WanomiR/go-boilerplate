package main

import (
	"log"
	"os"

	_ "github.com/wanomir/go-boilerplate/docs"
	"github.com/wanomir/go-boilerplate/internal/app"
)

// @title GoBoilerplate
// @version 0.0.2
// @description Go-boilerplate API

// @host localhost:8888
// @basePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(a.Run())
}
