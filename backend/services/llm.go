package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LLM struct {
	apiKey string
	model  string
}

func NewLLM(apiKey string) *LLM {
	return &LLM{
		apiKey: apiKey,
		model:  "mistralai/Mistral-7B-Instruct-v0.2",
	}
}

func (l *LLM) GenerateResponse(query string) (string, error) {
	url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", l.model)

	payload := map[string]string{
		"inputs": fmt.Sprintf(`Instructions: Answer the following question concisely. Do not repeat the question or include any prefixes like 'Answer:' or 'Context:' in your response.

Question: %s`, query),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling data: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", l.apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Clean up the response
	response := result[0]["generated_text"].(string)
	return extractAnswer(response), nil
}

func extractAnswer(response string) string {
	// Split by common delimiters
	delimiters := []string{
		"Answer:",
		"Response:",
		"\n",
		"Question:",
		"Context:",
	}

	cleanedResponse := response
	for _, delimiter := range delimiters {
		parts := strings.Split(cleanedResponse, delimiter)
		if len(parts) > 1 {
			// Always take the last meaningful part
			for i := len(parts) - 1; i >= 0; i-- {
				if trimmed := strings.TrimSpace(parts[i]); trimmed != "" {
					cleanedResponse = trimmed
					break
				}
			}
		}
	}

	// Remove any remaining prefixes
	prefixes := []string{
		"The context is about",
		"Based on the context,",
		"According to the context,",
		"The answer is",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(cleanedResponse), strings.ToLower(prefix)) {
			cleanedResponse = cleanedResponse[len(prefix):]
		}
	}

	return strings.TrimSpace(cleanedResponse)
}
