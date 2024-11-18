package models

// Document represents the data structure for storing documents
type Document struct {
	ID       string            `json:"id"`       // Unique identifier for the document
	Content  string            `json:"content"`  // Text content of the document
	Metadata map[string]string `json:"metadata"` // Additional metadata (e.g., author, timestamp)
}

// Embedding represents the vector embedding for a document or query
type Embedding struct {
	ID        string    `json:"id"`        // Unique identifier linked to a Document
	Vector    []float64 `json:"vector"`    // The embedding vector
	Dimension int       `json:"dimension"` // Dimensionality of the vector
}
