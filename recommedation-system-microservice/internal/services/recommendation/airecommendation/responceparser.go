package airecommendation

import (
	"encoding/json"
	"strings"
)

type JSONTitlesParser struct{}

func NewAIResponseParser() (*JSONTitlesParser, error) {
	return &JSONTitlesParser{}, nil
}

func normalizeRawResponse(raw string) string {
	raw = strings.TrimSpace(raw)

	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```json")
		raw = strings.TrimPrefix(raw, "```JSON")
		raw = strings.TrimPrefix(raw, "```")
		if end := strings.LastIndex(raw, "```"); end >= 0 {
			raw = raw[:end]
		}
		raw = strings.TrimSpace(raw)
	}

	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start >= 0 && end > start {
		return raw[start : end+1]
	}

	return raw
}

func (p *JSONTitlesParser) ParseTitles(raw string) ([]string, error) {
	normalized := normalizeRawResponse(raw)

	var titles []string
	if err := json.Unmarshal([]byte(normalized), &titles); err == nil {
		return titles, nil
	}

	var wrapped struct {
		Titles     []string `json:"titles"`
		ContentIDs []string `json:"content_ids"`
		Results    []string `json:"results"`
		Items      []string `json:"items"`
	}
	if err := json.Unmarshal([]byte(normalized), &wrapped); err != nil {
		return nil, ErrInvalidAIResponse
	}

	switch {
	case len(wrapped.Titles) > 0:
		return wrapped.Titles, nil
	case len(wrapped.ContentIDs) > 0:
		return wrapped.ContentIDs, nil
	case len(wrapped.Results) > 0:
		return wrapped.Results, nil
	case len(wrapped.Items) > 0:
		return wrapped.Items, nil
	default:
		return nil, ErrInvalidAIResponse
	}
}

// ParseContentIDs — алиас для обратной совместимости
func (p *JSONTitlesParser) ParseContentIDs(raw string) ([]string, error) {
	return p.ParseTitles(raw)
}
