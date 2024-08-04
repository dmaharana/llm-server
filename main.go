package main

import (
	"fmt"
	"gollm/appenv"
	"log"
	"net/http"
)

const (
	genApi  = "/api/generate"
	chatApi = "/api/chat"
	tagApi  = "/api/tags"
	model   = "llama3:latest"
)

var (
	url     = "http://localhost:11434"
	webPort = "8011"
)

type Config struct {
}

func main() {
	// withUserInput()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := Config{}

	appEnv, err := appenv.ReadConfig()
	if err == nil {
		url = appEnv.LlmUrl
		webPort = appEnv.AppPort
	}

	// start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting service on port %s", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
