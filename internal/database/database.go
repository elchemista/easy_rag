package database

import "github.com/MaxwellGroup/ragexp1/internal/models"

// database interface

// Database defines the interface for interacting with a database
type Database interface {
	SaveDocument(document models.Document) error        // the content will be chunked and saved
	GetDocumentInfo(id string) (models.Document, error) // return the document with the given id without content
	GetDocument(id string) (models.Document, error)     // return the document with the given id with content assembled
	Search(vector [][]float32) ([]models.Embedding, error)
	ListDocuments() ([]models.Document, error)
	DeleteDocument(id string) error
	SaveEmbeddings(embeddings []models.Embedding) error
	// to implement	in future
	// SearchByCategory(category []string) ([]Embedding, error)
	// SearchByMetadata(metadata map[string]string) ([]Embedding, error)
	// GetAllEmbeddingByDocumentID(documentID string) ([]Embedding, error)
}
