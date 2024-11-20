package llm

// implement llm interface
type LLMService interface {
	// generate text from prompt
	Generate(prompt string) (string, error)
}
