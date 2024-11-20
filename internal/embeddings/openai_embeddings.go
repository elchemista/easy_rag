package embeddings

type OpenAIEmbeddings struct {
	APIKey   string
	Endpoint string
	Model    string
}

func NewOpenAIEmbeddings(apiKey string, endpoint string, model string) *OpenAIEmbeddings {
	return &OpenAIEmbeddings{
		APIKey:   apiKey,
		Endpoint: endpoint,
		Model:    model,
	}
}

func (o *OpenAIEmbeddings) Generate(text string) ([]float64, error) {
	return nil, nil
}
