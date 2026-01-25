package airecommendation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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

	url := "https://router.huggingface.co/v1/chat/completions"

	reqBody := map[string]any{
		"model": "meta-llama/Meta-Llama-3-8B-Instruct",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
		"max_tokens":  256,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"hf error: status=%d body=%s",
			resp.StatusCode,
			string(rawBody),
		)
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(rawBody, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", ErrInvalidAIResponse
	}

	return parsed.Choices[0].Message.Content, nil
}
