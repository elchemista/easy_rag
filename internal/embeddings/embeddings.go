package embeddings

// implement embeddings interface
type EmbeddingsService interface {
	// generate embedding from text
	Generate(text string) ([]float64, error)
}
