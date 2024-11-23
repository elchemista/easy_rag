package main

import (
	"github.com/MaxwellGroup/ragexp1/api"
	"github.com/MaxwellGroup/ragexp1/config"
	"github.com/MaxwellGroup/ragexp1/internal/database"
	"github.com/MaxwellGroup/ragexp1/internal/embeddings"
	"github.com/MaxwellGroup/ragexp1/internal/llm"
	"github.com/MaxwellGroup/ragexp1/internal/pkg/rag"
	"github.com/labstack/echo/v4"
)

// Rag is the main struct for the rag application

func main() {
	cfg := config.NewConfig()

	llm := llm.NewOllama(cfg.OllamaEndpoint, cfg.OllamaModel)
	embeddings := embeddings.NewOllamaEmbeddings(cfg.OllamaEmbeddingEndpoint, cfg.OllamaEmbeddingModel)
	database := database.NewMilvus(cfg.MilvusHost)

	// Rag instance
	rag := rag.NewRag(llm, embeddings, database)

	// Echo WebServer instance
	e := echo.New()

	// Wrapper for API
	api.NewAPI(e, rag)

	// Start Server
	e.Logger.Fatal(e.Start(":4002"))
}
