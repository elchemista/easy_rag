# Rag Go Simplified

## Prerequisite

- malvius DB
- ollama

## models

- mistral:7b   
- bge-m3  (embeddings)


## For testing later
```go
err := rag.Database.SaveDocument(models.Document{
	ID:             "1",
	Content:        "Hello, World!",
	Link:           "https://www.google.com",
	Filename:       "hello.txt",
	Category:       []string{"test"},
	EmbeddingModel: "bge-m3",
	Summary:        "This is a test document",
	Vector:         generateMockVector(1024, 0.0, 0.1),
	Metadata:       map[string]string{"key": "value"},
})

fmt.Println(err)


func generateMockVector(dimension int, start, step float32) []float32 {
	vector := make([]float32, dimension)
	for i := 0; i < dimension; i++ {
		vector[i] = start + step*float32(i)
	}
	return vector
}
```
