package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func replyToMessage(number string, message string) {
	messages := []Message{
		{Role: "system", Content: config.SystemPrompt},
		{Role: "user", Content: message},
	}
	response, err := callOllamaChatAPI(config.Model, messages)
	if err != nil {
		fmt.Printf("Error talking to Ollama: %v\n", err)
		return
	}
	sendSms(number, response)
}

func callOllamaChatAPI(model string, messages []Message) (string, error) {
	url := config.OllamaBase + "/api/chat"

	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	// Marshal request into JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request JSON: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the streamed response line by line
	scanner := bufio.NewScanner(resp.Body)
	var fullResponse string

	for scanner.Scan() {
		line := scanner.Text()

		// Parse each line as a StreamChunk
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			return "", fmt.Errorf("failed to parse chunk: %w", err)
		}

		fullResponse += chunk.Message.Content

		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading response stream: %w", err)
	}

	return fullResponse, nil
}
