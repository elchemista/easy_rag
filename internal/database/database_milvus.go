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
	doc, err := m.Client.GetDocumentByID(ctx, id)

	if err != nil {
		return models.Document{}, err
	}

	if len(doc) == 0 {
		return models.Document{}, nil
	}
	return models.Document{
		ID:             doc["ID"].(string),
		Content:        doc["Content"].(string),
		Link:           doc["Link"].(string),
		Filename:       doc["Filename"].(string),
		Category:       doc["Category"].(string),
		EmbeddingModel: doc["EmbeddingModel"].(string),
		Summary:        doc["Summary"].(string),
		Metadata:       doc["Metadata"].(map[string]string),
	}, nil
}

func (m *Milvus) Search(content string) ([]models.Embedding, error) {
	return nil, nil
}

func (m *Milvus) ListDocuments() ([]models.Document, error) {
	ctx := context.Background()
	docs, err := m.Client.GetAllDocuments(ctx)
	if err != nil {
		return nil, err
	}
	var documents []models.Document
	for _, doc := range docs {
		documents = append(documents, models.Document{
			ID:             doc["ID"].(string),
			Content:        doc["Content"].(string),
			Link:           doc["Link"].(string),
			Filename:       doc["Filename"].(string),
			Category:       doc["Category"].(string),
			EmbeddingModel: doc["EmbeddingModel"].(string),
			Summary:        doc["Summary"].(string),
			Metadata:       doc["Metadata"].(map[string]string),
		})
	}
	return documents, nil
}

func (m *Milvus) DeleteDocument(id string) error {
	ctx := context.Background()
	err := m.Client.DeleteDocument(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
