package api

import (
	"sync"

	"esp-organizer/internal/aws/llm"
)

var (
	embeddingClientOnce sync.Once
	embeddingClient     *llm.LlamaClient
)

func getEmbeddingClient() *llm.LlamaClient {
	embeddingClientOnce.Do(func() {
		embeddingClient = llm.NewLlamaClient()
	})
	return embeddingClient
}
