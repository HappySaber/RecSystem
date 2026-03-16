// internal/services/recommendation/airecommendation/stub.go
package airecommendation

import "context"

// StubAIClient — заглушка для тестов, не требует API ключа
type StubAIClient struct{}

func NewStubAIClient() *StubAIClient {
	return &StubAIClient{}
}

func (s *StubAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	// возвращаем валидный пустой JSON чтобы парсер не падал
	return "[]", nil
}
