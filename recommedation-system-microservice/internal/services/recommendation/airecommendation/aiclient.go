package airecommendation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type AIClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("HF_API_KEY")
	if apiKey == "" {
		return nil, errors.New("HF_API_KEY is not set")
	}
	return &OpenAIClient{apiKey: apiKey}, nil
}

func (c *OpenAIClient) Complete(
	ctx context.Context,
	prompt string,
) (string, error) {

	url := "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.2"

	bodyBytes, _ := json.Marshal(map[string]string{"inputs": prompt})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	text, ok := respData[0]["generated_text"].(string)
	if !ok {
		return "", ErrInvalidAIResponse
	}
	return text, nil
}
