package embeddings

// implement embeddings interface
type EmbeddingsService interface {
	// generate embedding from text
	Vectorize(text string) ([][]float32, error)
	GetModel() string
}
