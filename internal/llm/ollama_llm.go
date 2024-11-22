package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ollama struct {
	Endpoint string
	Model    string
}

func NewOllama(endpoint string, model string) *Ollama {
	return &Ollama{
		Endpoint: endpoint,
		Model:    model,
	}
}

// Generate sends a prompt to the Ollama endpoint and returns the response
func (o *Ollama) Generate(prompt string) (string, error) {
	// Create the request payload
	payload := map[string]interface{}{
		"model": o.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the POST request
	resp, err := http.Post(o.Endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned error: %s", string(body))
	}

	// Unmarshal the response into a map to extract content
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract and return the content from the response
	if content, ok := response["content"].(string); ok {
		return content, nil
	}

	return "", fmt.Errorf("unexpected response format: %s", string(body))

}

func (o *Ollama) GetModel() string {
	return o.Model
}
