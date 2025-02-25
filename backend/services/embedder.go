package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Embedder struct {
	apiKey string
	model  string
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func NewEmbedder(apiKey string) *Embedder {
	return &Embedder{
		apiKey: apiKey,
		model:  "sentence-transformers/all-MiniLM-L6-v2",
	}
}

type HuggingFaceEmbedder struct {
	apiKey string
	model  string
}

type HuggingFaceResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func NewHuggingFaceEmbedder(apiKey string) *HuggingFaceEmbedder {
	return &HuggingFaceEmbedder{
		apiKey: apiKey,
		model:  "sentence-transformers/all-MiniLM-L6-v2",
	}
}

func (e *HuggingFaceEmbedder) GetEmbeddings(text string) ([]float32, error) {
	url := "https://api-inference.huggingface.co/pipeline/feature-extraction/sentence-transformers/all-MiniLM-L6-v2"
	jsonData, err := json.Marshal(map[string]string{
		"inputs": text,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var rawNumbers []float64
	if err := json.NewDecoder(resp.Body).Decode(&rawNumbers); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	result := make([]float32, len(rawNumbers))
	for i, num := range rawNumbers {
		result[i] = float32(num)
	}

	return result, nil
}
