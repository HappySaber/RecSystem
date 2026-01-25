package airecommendation

import "encoding/json"

type JSONContentIDParser struct{}

func NewAIResponseParser() (*JSONContentIDParser, error) {
	return &JSONContentIDParser{}, nil
}

func (p *JSONContentIDParser) ParseContentIDs(
	raw string,
) ([]string, error) {
	// 1️⃣ пробуем идеальный вариант
	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err == nil {
		return ids, nil
	}

	// 2️⃣ пробуем объект-обёртку
	var wrapped struct {
		ContentIDs []string `json:"content_ids"`
		Results    []string `json:"results"`
		Items      []string `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &wrapped); err != nil {
		return nil, ErrInvalidAIResponse
	}

	switch {
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
