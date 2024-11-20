package embeddings

type OllamaEmbeddings struct {
	Endpoint string
	Model    string
}

func NewOllamaEmbeddings(endpoint string, model string) *OllamaEmbeddings {
	return &OllamaEmbeddings{
		Endpoint: endpoint,
		Model:    model,
	}
}

func (o *OllamaEmbeddings) Generate(text string) ([]float64, error) {
	return nil, nil
}
