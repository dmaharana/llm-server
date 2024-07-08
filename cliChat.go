package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type RequestData struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type RequestPayload struct {
	Model        string         `json:"model"`
	Prompt       string         `json:"prompt"`
	ChatMessages []Conversation `json:"conversation"`
}

type Conversation struct {
	Id       string `json:"id"`
	Query    string `json:"query"`
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
}

type ResponseData struct {
	Response  string `json:"response"`
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

func spinner(delay time.Duration, done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c thinking...", r)
				time.Sleep(delay)
			}
		}
	}
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your question: ")
	question, _ := reader.ReadString('\n')
	return strings.TrimSpace(question)
}

func withUserInput() {
	// Get user input
	question := getUserInput()

	prompt := promptSetup(question, "", "1")

	data := RequestData{
		Model:  model,
		Prompt: prompt,
		// Prompt: "write me a story about python.",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Start the spinner
	done := make(chan bool)
	go spinner(100*time.Millisecond, done)

	startTime := time.Now()

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		done <- true
		fmt.Println("\nError making request:", err)
		return
	}
	defer resp.Body.Close()

	// Stop the spinner
	done <- true
	fmt.Print("\r") // Clear the spinner line

	// Calculate and display the response time
	responseTime := time.Since(startTime)
	fmt.Printf("\rFirst response received in: %v\n", responseTime)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading response:", err)
			return
		}

		var responseData ResponseData
		err = json.Unmarshal(line, &responseData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		fmt.Print(responseData.Response)
	}

	// Calculate and display the response time
	finalResponseTime := time.Since(startTime)
	fmt.Printf("\n\n\rResponse received in: %v\n", finalResponseTime)
}

// depending on the option update the prompt with additional information
func promptSetup(prompt string, context string, option string) string {

	switch option {
	case "1":
		return fmt.Sprintf(`Instruction: let's think step-by-step through this problem: {{%s}}.

output: only a JSON object with two attributes: reasoning and final answer.`, prompt)
	case "2":
		return fmt.Sprintf(`instruction: answer the following question only if you know the answer, the answer appears in the context, or you can make a well-informed guess; otherwise tell me you don't know.
		===
		context: {{%s}} ====== now, answer this question: {{%s}}`, context, prompt)

	default:
		return prompt
	}
}
