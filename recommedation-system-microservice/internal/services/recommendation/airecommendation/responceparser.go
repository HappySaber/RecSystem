package airecommendation

import "encoding/json"

type JSONContentIDParser struct{}

func NewAIResponseParser() (*JSONContentIDParser, error) {
	return &JSONContentIDParser{}, nil
}

func (p *JSONContentIDParser) ParseContentIDs(
	raw string,
) ([]string, error) {

	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, ErrInvalidAIResponse
	}

	return ids, nil
}
