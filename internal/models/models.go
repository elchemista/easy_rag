package models

// Document represents the data structure for storing documents
type Document struct {
	ID             string            `json:"id"`              // Unique identifier for the document
	Content        string            `json:"content"`         // Text content of the document become chunks of data will not be saved
	Link           string            `json:"link"`            // Link to the document
	Filename       string            `json:"filename"`        // Filename of the document
	Category       []string          `json:"category"`        // Category of the document
	EmbeddingModel string            `json:"embedding_model"` // Embedding model used to generate the embedding
	Summary        string            `json:"summary"`         // Summary of the document
	Vector         []float32         `json:"vector"`          // The embedding vector
	Metadata       map[string]string `json:"metadata"`        // Additional metadata (e.g., author, timestamp)
}

// Embedding represents the vector embedding for a document or query
type Embedding struct {
	ID         string    `json:"id"`          // Unique identifier
	DocumentID string    `json:"document_id"` // Unique identifier linked to a Document
	Vector     []float32 `json:"vector"`      // The embedding vector
	TextChunk  string    `json:"text_chunk"`  // Text chunk of the document
	Dimension  int       `json:"dimension"`   // Dimensionality of the vector
	Order      int       `json:"order"`       // Order of the embedding to build the content back

	// maybe later adding summary and metadata
}
