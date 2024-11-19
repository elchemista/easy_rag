package database

// database interface

// Database defines the interface for interacting with a database
type Database interface {
	SaveDocument(document Document) error        // the content will be chunked and saved
	GetDocument(id string) (Document, error)     // it will return the document with the given id with content assembled
	GetFullDocument(id string) (Document, error) // it will return the document with the given id with content assembled
	Search(content string) ([]Embedding, error)
	// to implement	in future
	// SearchByCategory(category []string) ([]Embedding, error)
	// SearchByMetadata(metadata map[string]string) ([]Embedding, error)
	// GetAllEmbeddingByDocumentID(documentID string) ([]Embedding, error)
}
