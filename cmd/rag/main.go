package main

import (
	"github.com/elchemista/easy_rag/api"
	"github.com/elchemista/easy_rag/config"
	"github.com/elchemista/easy_rag/internal/database"
	"github.com/elchemista/easy_rag/internal/embeddings"
	"github.com/elchemista/easy_rag/internal/llm"
	"github.com/elchemista/easy_rag/internal/pkg/rag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Rag is the main struct for the rag application

func main() {
	cfg := config.NewConfig()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	llm := llm.NewOpenAI(cfg.OpenAIAPIKey, cfg.OpenAIEndpoint, cfg.OpenAIModel)
	embeddings := embeddings.NewOpenAIEmbeddings(cfg.OpenAIEmbeddingAPIKey, cfg.OpenAIEmbeddingEndpoint, cfg.OpenAIEmbeddingModel)
	database := database.NewMilvus(cfg.MilvusHost)

	rag := rag.NewRag(llm, embeddings, database)

	api.StartServer(e, rag)
}
