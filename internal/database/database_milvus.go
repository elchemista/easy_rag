package database

import "github.com/elchemista/easy_rag/internal/models"

// implement database interface for milvus
type Milvus struct {
	Host string
	// Client *Client
}

func NewMilvus(host string) *Milvus {
	return &Milvus{
		Host: host,
	}
}

func (m *Milvus) SaveDocument(document models.Document) error {
	return nil
}

func (m *Milvus) GetDocumentInfo(id string) (models.Document, error) {
	return models.Document{}, nil
}

func (m *Milvus) GetDocument(id string) (models.Document, error) {
	return models.Document{}, nil
}

func (m *Milvus) Search(content string) ([]models.Embedding, error) {
	return nil, nil
}

func (m *Milvus) ListDocuments() ([]models.Document, error) {
	return nil, nil
}

func (m *Milvus) DeleteDocument(id string) error {
	return nil
}
