package rag

import (
	"github.com/MaxwellGroup/ragexp1/internal/database"
	"github.com/MaxwellGroup/ragexp1/internal/embeddings"
	"github.com/MaxwellGroup/ragexp1/internal/llm"
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
