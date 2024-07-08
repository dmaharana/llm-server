package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	url   = "http://localhost:11434/api/generate"
	model = "llama3"
	// model = "qwen27b:5"
	// model = "gemma1:q8"
)

const webPort = "8011"

type Config struct {
}

func main() {
	// withUserInput()

	app := Config{}
	// start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting service on port %s", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
