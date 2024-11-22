package config

import cfg "github.com/eschao/config"

type Config struct {
	// LLM
	OpenAIAPIKey   string `env:"OPENAI_API_KEY"`
	OpenAIEndpoint string `env:"OPENAI_ENDPOINT"`
	OpenAIModel    string `env:"OPENAI_MODEL"`
	OllamaEndpoint string `env:"OLLAMA_ENDPOINT"`
	OllamaModel    string `env:"OLLAMA_MODEL"`

	// Embeddings
	OpenAIEmbeddingAPIKey   string `env:"OPENAI_EMBEDDING_API_KEY"`
	OpenAIEmbeddingEndpoint string `env:"OPENAI_EMBEDDING_ENDPOINT"`
	OpenAIEmbeddingModel    string `env:"OPENAI_EMBEDDING_MODEL"`
	OllamaEmbeddingEndpoint string `env:"OLLAMA_EMBEDDING_ENDPOINT"`
	OllamaEmbeddingModel    string `env:"OLLAMA_EMBEDDING_MODEL"`

	// Database
	MilvusHost string `env:"MILVUS_HOST"`
}

func NewConfig() Config {
	config := Config{
		MilvusHost:              "localhost:19530",
		OllamaEmbeddingEndpoint: "http://localhost:11434",
		OllamaEmbeddingModel:    "bge-m3",
		OllamaEndpoint:          "http://localhost:11434/api/chat",
		OllamaModel:             "llama3.2:3b",
	}
	cfg.ParseEnv(&config)
	return config
}
