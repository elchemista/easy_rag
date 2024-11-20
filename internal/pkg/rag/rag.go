package rag

import (
	"github.com/elchemista/easy_rag/internal/database"
	"github.com/elchemista/easy_rag/internal/embeddings"
	"github.com/elchemista/easy_rag/internal/llm"
)

type Rag struct {
	llm        llm.LLMService
	embeddings embeddings.EmbeddingsService
	database   database.Database
}

func NewRag(llm llm.LLMService, embeddings embeddings.EmbeddingsService, database database.Database) *Rag {
	return &Rag{
		llm:        llm,
		embeddings: embeddings,
		database:   database,
	}
}
