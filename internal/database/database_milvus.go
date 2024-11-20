package database

import (
	"context"

	"github.com/elchemista/easy_rag/internal/models"
	"github.com/elchemista/easy_rag/internal/pkg/database/milvus"
)

// implement database interface for milvus
type Milvus struct {
	Host   string
	Client *milvus.Client
}

func NewMilvus(host string) *Milvus {

	milviusClient, err := milvus.NewClient(host)
	if err != nil {
		panic(err)
	}

	return &Milvus{
		Host:   host,
		Client: milviusClient,
	}
}

func (m *Milvus) SaveDocument(document models.Document) error {
	// for now lets use context background
	ctx := context.Background()
	err := m.Client.EnsureCollections(ctx)
	if err != nil {
		return err
	}
	err = m.Client.InsertDocuments(ctx, []models.Document{document})
	if err != nil {
		return err
	}
	return nil
}

func (m *Milvus) GetDocumentInfo(id string) (models.Document, error) {
	ctx := context.Background()
	docs, err := m.Client.GetDocumentByID(ctx, id)
	if err != nil {
		return models.Document{}, err
	}
	if len(docs) == 0 {
		return models.Document{}, nil
	}
	return models.Document{}, nil
}

func (m *Milvus) GetDocument(id string) (models.Document, error) {
	ctx := context.Background()
	docs, err := m.Client.GetDocumentByID(ctx, id)
	if err != nil {
		return models.Document{}, err
	}
	if len(docs) == 0 {
		return models.Document{}, nil
	}
	return models.Document{}, nil
}

func (m *Milvus) Search(content string) ([]models.Embedding, error) {
	return nil, nil
}

func (m *Milvus) ListDocuments() ([]models.Document, error) {
	return nil, nil
}

func (m *Milvus) DeleteDocument(id string) error {
	ctx := context.Background()
	err := m.Client.DeleteDocument(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
