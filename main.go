package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	url    = "http://localhost:11434"
	genApi = "/api/generate"
	tagApi = "/api/tags"
	model  = "llama3"
	// model = "qwen27b:5"
	// model = "gemma1:q8"
)

const webPort = "8011"

type Config struct {
}

func main() {
	// withUserInput()
        log.SetFlags(log.LstdFlags | log.Lshortfile)
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
