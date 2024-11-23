package database

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/MaxwellGroup/ragexp1/internal/models"
	"github.com/MaxwellGroup/ragexp1/internal/pkg/database/milvus"
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

	return m.Client.InsertDocuments(ctx, []models.Document{document})
}

func (m *Milvus) SaveEmbeddings(embeddings []models.Embedding) error {
	ctx := context.Background()
	return m.Client.InsertEmbeddings(ctx, embeddings)
}

func (m *Milvus) GetDocumentInfo(id string) (models.Document, error) {
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
		Link:           doc["Link"].(string),
		Filename:       doc["Filename"].(string),
		Category:       doc["Category"].(string),
		EmbeddingModel: doc["EmbeddingModel"].(string),
		Summary:        doc["Summary"].(string),
		Metadata:       doc["Metadata"].(map[string]string),
	}, nil
}

func (m *Milvus) GetDocument(id string) (models.Document, error) {
	ctx := context.Background()
	doc, err := m.Client.GetDocumentByID(ctx, id)
	if err != nil {
		return models.Document{}, err
	}

	embeds, err := m.Client.GetAllEmbeddingByDocID(ctx, id)

	if err != nil {
		return models.Document{}, err
	}

	// order embed by order
	sort.Slice(embeds, func(i, j int) bool {
		return embeds[i].Order < embeds[j].Order
	})

	// concatenate text chunks
	var buf bytes.Buffer
	for _, embed := range embeds {
		buf.WriteString(embed.TextChunk)
	}

	textChunks := buf.String()

	if len(doc) == 0 {
		return models.Document{}, nil
	}
	return models.Document{
		ID:             doc["ID"].(string),
		Content:        textChunks,
		Link:           doc["Link"].(string),
		Filename:       doc["Filename"].(string),
		Category:       doc["Category"].(string),
		EmbeddingModel: doc["EmbeddingModel"].(string),
		Summary:        doc["Summary"].(string),
		Metadata:       doc["Metadata"].(map[string]string),
	}, nil
}

func (m *Milvus) Search(vector [][]float32) ([]models.Embedding, error) {
	ctx := context.Background()
	results, err := m.Client.Search(ctx, vector, 10)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (m *Milvus) ListDocuments() ([]models.Document, error) {
	ctx := context.Background()

	docs, err := m.Client.GetAllDocuments(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get docs: %w", err)
	}

	return docs, nil
}

func (m *Milvus) DeleteDocument(id string) error {
	ctx := context.Background()
	err := m.Client.DeleteDocument(ctx, id)
	if err != nil {
		return err
	}

	err = m.Client.DeleteEmbedding(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
