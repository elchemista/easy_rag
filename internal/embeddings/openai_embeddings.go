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

func (o *OpenAIEmbeddings) Vectorize(text string) ([]float32, error) {
	return nil, nil
}

func (o *OpenAIEmbeddings) GetModel() string {
	return o.Model
}
