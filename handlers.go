package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func (app *Config) ChatResponse(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	var reqPayload RequestPayload
	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var finalRes ResponseData
	finalRes.Response = reqPayload.Prompt

	log.Printf("Received model: %s", reqPayload.Model)
	log.Printf("Received question: %s", reqPayload.Prompt)

	responses := []string{}

	var llmRequest RequestData
	llmRequest.Model = reqPayload.Model
	llmRequest.Prompt = reqPayload.Prompt

	jsonData, err := json.Marshal(llmRequest)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Printf("Sending request: %s", string(jsonData))

	res, err := http.Post(url+genApi, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	log.Printf("Status code: %d", res.StatusCode)

	if res.StatusCode != 200 {
		app.errorJSON(w, fmt.Errorf("(%d) %s: Not able to process request with the selected model", res.StatusCode, http.StatusText(res.StatusCode)))
		return
	}

	ws := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			app.errorJSON(w, err)
			return
		}

		var llmRes ResponseData
		err = json.Unmarshal(line, &llmRes)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		responses = append(responses, llmRes.Response)
		finalRes.Response = strings.Join(responses, "")
		finalRes.Done = llmRes.Done

		// log.Printf("Response: %s", llmRes.Response)
		app.writeJSON(w, http.StatusOK, llmRes)
		// app.writeJSON(w, http.StatusOK, finalRes)

		// Flush the response
		ws.Flush()
	}

	// app.writeJSON(w, http.StatusOK, finalRes)
}

// handle model list get
func (app *Config) GetModels(w http.ResponseWriter, r *http.Request) {
	jsonresp := jsonResponse{
		Error:   false,
		Message: "success",
	}

	// call the tag list endpoint
	res, err := http.Get(url + tagApi)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()
	var models Models
	err = json.NewDecoder(res.Body).Decode(&models)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jsonresp.Data = models

	app.writeJSON(w, http.StatusOK, jsonresp)
}
