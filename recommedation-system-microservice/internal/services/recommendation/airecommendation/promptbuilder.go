package airecommendation

import "fmt"

type PromptBuilder interface {
	BuildExplicit(query string, limit int) string
}

type DefaultPromptBuilder struct{}

func (p *DefaultPromptBuilder) BuildExplicit(query string, limit int) string {
	return fmt.Sprintf(`
You are a recommendation engine.

User request:
"%s"

Return EXACTLY %d content IDs from catalog.
Format response as JSON array of strings.
Example:
["anime-1","book-7"]

DO NOT add explanations.
`, query, limit)
}
