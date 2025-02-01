package api_test

import (
	"github.com/elchemista/easy_rag/internal/models"
	"github.com/stretchr/testify/mock"
)

// --------------------
// Mock LLM
// --------------------
type MockLLMService struct {
	mock.Mock
}

func (m *MockLLMService) Generate(prompt string) (string, error) {
	args := m.Called(prompt)
	return args.String(0), args.Error(1)
}

func (m *MockLLMService) GetModel() string {
	args := m.Called()
	return args.String(0)
}

// --------------------
// Mock Embeddings
// --------------------
type MockEmbeddingsService struct {
	mock.Mock
}

func (m *MockEmbeddingsService) Vectorize(text string) ([][]float32, error) {
	args := m.Called(text)
	return args.Get(0).([][]float32), args.Error(1)
}

func (m *MockEmbeddingsService) GetModel() string {
	args := m.Called()
	return args.String(0)
}

// --------------------
// Mock Database
// --------------------
type MockDatabase struct {
	mock.Mock
}

// SaveDocument(document Document) error
func (m *MockDatabase) SaveDocument(doc models.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

// SaveEmbeddings([]Embedding) error
func (m *MockDatabase) SaveEmbeddings(emb []models.Embedding) error {
	args := m.Called(emb)
	return args.Error(0)
}

// ListDocuments() ([]Document, error)
func (m *MockDatabase) ListDocuments() ([]models.Document, error) {
	args := m.Called()
	return args.Get(0).([]models.Document), args.Error(1)
}

// GetDocument(id string) (Document, error)
func (m *MockDatabase) GetDocument(id string) (models.Document, error) {
	args := m.Called(id)
	return args.Get(0).(models.Document), args.Error(1)
}

// DeleteDocument(id string) error
func (m *MockDatabase) DeleteDocument(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Search(vector []float32) ([]models.Embedding, error)
func (m *MockDatabase) Search(vector [][]float32) ([]models.Embedding, error) {
	args := m.Called(vector)
	return args.Get(0).([]models.Embedding), args.Error(1)
}
