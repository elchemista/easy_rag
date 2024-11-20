package rag

import (
	"github.com/elchemista/easy_rag/internal/database"
	"github.com/elchemista/easy_rag/internal/embeddings"
	"github.com/elchemista/easy_rag/internal/llm"
)

type Rag struct {
	LLM        llm.LLMService
	Embeddings embeddings.EmbeddingsService
	Database   database.Database
}

func NewRag(llm llm.LLMService, embeddings embeddings.EmbeddingsService, database database.Database) *Rag {
	return &Rag{
		LLM:        llm,
		Embeddings: embeddings,
		Database:   database,
	}
}
