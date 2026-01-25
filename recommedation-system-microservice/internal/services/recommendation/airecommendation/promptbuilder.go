package airecommendation

import "fmt"

type DefaultPromptBuilder struct{}

func NewPromptBuilder() (*DefaultPromptBuilder, error) {
	return &DefaultPromptBuilder{}, nil
}

func (p *DefaultPromptBuilder) BuildExplicit(query string, limit int) string {
	return fmt.Sprintf(`
You are a recommendation engine.

User request:
"%s"

Return ONLY valid JSON.
NO text before or after.

JSON schema:
["Movie name 1","Movie name 2"]

Return EXACTLY %d items.
`, query, limit)
}
