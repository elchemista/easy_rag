package llm

type Ollama struct {
	Endpoint string
	Model    string
}

func NewOllama(endpoint string, model string) *Ollama {
	return &Ollama{
		Endpoint: endpoint,
		Model:    model,
	}
}

func (o *Ollama) Generate(prompt string) (string, error) {
	return "", nil
	// TODO: implement
}
