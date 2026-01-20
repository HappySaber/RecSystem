package airecommendation

import "encoding/json"

type AIResponseParser interface {
	ParseContentIDs(raw string) ([]string, error)
}

type JSONContentIDParser struct{}

func (p *JSONContentIDParser) ParseContentIDs(
	raw string,
) ([]string, error) {

	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, ErrInvalidAIResponse
	}

	return ids, nil
}
