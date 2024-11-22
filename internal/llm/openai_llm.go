package llm

type OpenAI struct {
	APIKey   string
	Endpoint string
	Model    string
}

func NewOpenAI(apiKey string, endpoint string, model string) *OpenAI {
	return &OpenAI{
		APIKey:   apiKey,
		Endpoint: endpoint,
		Model:    model,
	}
}

func (o *OpenAI) Generate(prompt string) (string, error) {
	return "", nil
	// TODO: implement
}

func (o *OpenAI) GetModel() string {
	return o.Model
}
