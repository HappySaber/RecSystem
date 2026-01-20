package airecommendation

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type AIClient interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("HF_API_KEY")
	return &OpenAIClient{apiKey: apiKey}, nil
}

func (c *OpenAIClient) Complete(
	ctx context.Context,
	prompt string,
) (string, error) {

	url := "https://api-inference.huggingface.co/models/gpt2"

	bodyBytes, _ := json.Marshal(map[string]string{"inputs": prompt})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	return respData[0]["generated_text"].(string), nil
}
