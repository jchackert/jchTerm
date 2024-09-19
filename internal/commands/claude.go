package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jchackert/jchterm/internal/config"
)

func executeAsk(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("Usage: ask <your question>")
	}
	question := strings.Join(args, " ")
	return askClaude(question)
}

func askClaude(question string) (string, error) {
	payload := map[string]string{
		"prompt": question,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error preparing request: %v", err)
	}

	req, err := http.NewRequest("POST", config.ApiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("Error parsing response: %v", err)
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("Error: Unexpected response format")
	}

	return fmt.Sprintf("Claude says: %s", response), nil
}
