package airecommendation_test

import (
	"testing"

	"rec-system-microservice/internal/services/recommendation/airecommendation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newParser() *airecommendation.JSONTitlesParser {
	p, _ := airecommendation.NewAIResponseParser()
	return p
}

func TestParseContentIDs_PlainArray(t *testing.T) {
	parser := newParser()

	result, err := parser.ParseContentIDs(`["Naruto","Bleach","One Piece"]`)

	require.NoError(t, err)
	assert.Equal(t, []string{"Naruto", "Bleach", "One Piece"}, result)
}

func TestParseContentIDs_WrappedContentIDs(t *testing.T) {
	parser := newParser()

	result, err := parser.ParseContentIDs(`{"content_ids":["id-1","id-2"]}`)

	require.NoError(t, err)
	assert.Equal(t, []string{"id-1", "id-2"}, result)
}

func TestParseContentIDs_WrappedResults(t *testing.T) {
	parser := newParser()

	result, err := parser.ParseContentIDs(`{"results":["id-1","id-2"]}`)

	require.NoError(t, err)
	assert.Equal(t, []string{"id-1", "id-2"}, result)
}

func TestParseContentIDs_WrappedItems(t *testing.T) {
	parser := newParser()

	result, err := parser.ParseContentIDs(`{"items":["id-1"]}`)

	require.NoError(t, err)
	assert.Equal(t, []string{"id-1"}, result)
}

func TestParseContentIDs_InvalidJSON(t *testing.T) {
	parser := newParser()

	_, err := parser.ParseContentIDs(`not json at all`)

	require.Error(t, err)
	assert.ErrorIs(t, err, airecommendation.ErrInvalidAIResponse)
}

func TestParseContentIDs_EmptyWrappedObject(t *testing.T) {
	parser := newParser()

	_, err := parser.ParseContentIDs(`{"unknown_field":["id-1"]}`)

	require.Error(t, err)
	assert.ErrorIs(t, err, airecommendation.ErrInvalidAIResponse)
}

func TestParseContentIDs_MarkdownBlock(t *testing.T) {
	parser := newParser()

	raw := "```json\n[\"The Shining\", \"Hereditary\"]\n```"
	result, err := parser.ParseContentIDs(raw)

	require.NoError(t, err)
	assert.Equal(t, []string{"The Shining", "Hereditary"}, result)
}

func TestParseContentIDs_WrappedTitles(t *testing.T) {
	parser := newParser()

	result, err := parser.ParseContentIDs(`{"titles":["Film A","Film B"]}`)

	require.NoError(t, err)
	assert.Equal(t, []string{"Film A", "Film B"}, result)
}
