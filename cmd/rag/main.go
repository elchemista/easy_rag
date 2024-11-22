package main

import (
	"github.com/elchemista/easy_rag/api"
	"github.com/elchemista/easy_rag/config"
	"github.com/elchemista/easy_rag/internal/database"
	"github.com/elchemista/easy_rag/internal/embeddings"
	"github.com/elchemista/easy_rag/internal/llm"
	"github.com/elchemista/easy_rag/internal/pkg/rag"
	"github.com/labstack/echo/v4"
)

// Rag is the main struct for the rag application

func main() {
	cfg := config.NewConfig()

	llm := llm.NewOpenAI(cfg.OpenAIAPIKey, cfg.OpenAIEndpoint, cfg.OpenAIModel)
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
