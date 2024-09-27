package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jchackert/jchterm/internal/config"
	"github.com/jchackert/jchterm/internal/logger"
)

func executeAsk(args []string) (string, error) {
	logger.Log("executeAsk called with args: %v", args)
	if len(args) == 0 {
		logger.Log("executeAsk: No arguments provided")
		return "", fmt.Errorf("Usage: ask <your question>")
	}
	question := strings.Join(args, " ")
	logger.Log("executeAsk: Calling askClaude with question: %s", question)
	return askClaude(question)
}

func askClaude(question string) (string, error) {
	logger.Log("askClaude called with question: %s", question)
	payload := map[string]interface{}{
		"model":      "claude-3-opus-20240229",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "user", "content": question},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Log("askClaude: Error preparing request: %v", err)
		return "", fmt.Errorf("Error preparing request: %v", err)
	}

	logger.Log("askClaude: Sending request to API URL: %s", config.ApiURL)
	req, err := http.NewRequest("POST", config.ApiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Log("askClaude: Error creating request: %v", err)
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", config.GetApiKey())
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log("askClaude: Error sending request: %v", err)
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	logger.Log("askClaude: Received response with status: %s", resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log("askClaude: Error reading response: %v", err)
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	logger.Log("askClaude: Response body: %s", string(body))

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Log("askClaude: Error parsing response: %v", err)
		return "", fmt.Errorf("Error parsing response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := "Unknown error"
		if errObj, ok := result["error"].(map[string]interface{}); ok {
			if msg, ok := errObj["message"].(string); ok {
				errorMsg = msg
			}
		}
		logger.Log("askClaude: API error: %s", errorMsg)
		return "", fmt.Errorf("API error: %s", errorMsg)
	}

	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		logger.Log("askClaude: Unexpected response format")
		return "", fmt.Errorf("Error: Unexpected response format")
	}

	firstContent, ok := content[0].(map[string]interface{})
	if !ok {
		logger.Log("askClaude: Unexpected content format")
		return "", fmt.Errorf("Error: Unexpected content format")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		logger.Log("askClaude: Missing text in content")
		return "", fmt.Errorf("Error: Missing text in content")
	}

	logger.Log("askClaude: Returning response: %s", text)
	return fmt.Sprintf("Claude says: %s", text), nil
}
