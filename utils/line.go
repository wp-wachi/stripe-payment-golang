package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// LINEMessage represents the request body for LINE API
type LINEMessage struct {
	To       string `json:"to,omitempty"`
	Messages []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"messages"`
}

// SendMessageToLINE sends a message to a LINE user or group
func SendMessageToLINE(message string) error {
	lineToken := os.Getenv("LINE_ACCESS_TOKEN")
	if lineToken == "" {
		log.Println("LINE_ACCESS_TOKEN is not set")
		return nil
	}

	// Create request body
	body := LINEMessage{
		Messages: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			{Type: "text", Text: message},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}

	// Send request to LINE API
	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/broadcast", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lineToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to LINE:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send message to LINE: %v\n", resp.Status)
	}

	return nil
}