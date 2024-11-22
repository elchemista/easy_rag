package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaEmbeddings struct {
	Endpoint string
	Model    string
}

func NewOllamaEmbeddings(endpoint string, model string) *OllamaEmbeddings {
	return &OllamaEmbeddings{
		Endpoint: endpoint,
		Model:    model,
	}
}

// Vectorize generates an embedding for the provided text
func (o *OllamaEmbeddings) Vectorize(text string) ([][]float32, error) {
	// Define the request payload
	payload := map[string]string{
		"model": o.Model,
		"input": text,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/api/embed", o.Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response: %s", body)
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Assuming the response JSON contains an "embedding" field with a float32 array
	var response struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Embeddings, nil
}

func (o *OllamaEmbeddings) GetModel() string {
	return o.Model
}
