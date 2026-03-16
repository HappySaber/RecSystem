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

func (p *DefaultPromptBuilder) BuildImplicit(
	genres map[string]float64,
	limit int,
) string {

	return fmt.Sprintf(`
You are a movie recommendation engine.

User genre preferences:
%v

Recommend movies matching these genres.

Return ONLY valid JSON.
NO text before or after.

JSON schema:
["Movie name 1","Movie name 2"]

Return EXACTLY %d items.
`, genres, limit)
}
