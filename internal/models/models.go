package models

// type VectorEmbedding [][]float32
// type Vector []float32
// Document represents the data structure for storing documents
type Document struct {
	ID             string            `json:"id" milvus:"ID"`                          // Unique identifier for the document
	Content        string            `json:"content" milvus:"Content"`                // Text content of the document become chunks of data will not be saved
	Link           string            `json:"link" milvus:"Link"`                      // Link to the document
	Filename       string            `json:"filename" milvus:"Filename"`              // Filename of the document
	Category       string            `json:"category" milvus:"Category"`              // Category of the document
	EmbeddingModel string            `json:"embedding_model" milvus:"EmbeddingModel"` // Embedding model used to generate the embedding
	Summary        string            `json:"summary" milvus:"Summary"`                // Summary of the document
	Metadata       map[string]string `json:"metadata" milvus:"Metadata"`              // Additional metadata (e.g., author, timestamp)
	Vector         []float32         `json:"vector" milvus:"Vector"`
}

// Embedding represents the vector embedding for a document or query
type Embedding struct {
	ID         string    `json:"id" milvus:"ID"`                  // Unique identifier
	DocumentID string    `json:"document_id" milvus:"DocumentID"` // Unique identifier linked to a Document
	Vector     []float32 `json:"vector" milvus:"Vector"`          // The embedding vector
	TextChunk  string    `json:"text_chunk" milvus:"TextChunk"`   // Text chunk of the document
	Dimension  int64     `json:"dimension" milvus:"Dimension"`    // Dimensionality of the vector
	Order      int64     `json:"order" milvus:"Order"`            // Order of the embedding to build the content back
	Score      float32   `json:"score"`                           // Score of the embedding
}
