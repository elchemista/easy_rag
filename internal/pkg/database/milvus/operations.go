package milvus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elchemista/easy_rag/internal/models"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// InsertDocuments inserts documents into the "documents" collection.
func (m *Client) InsertDocuments(ctx context.Context, docs []models.Document) error {
	idColumn := entity.NewColumnVarChar("ID", extractIDs(docs))
	contentColumn := entity.NewColumnVarChar("Content", extractContents(docs))
	linkColumn := entity.NewColumnVarChar("Link", extractLinks(docs))
	filenameColumn := entity.NewColumnVarChar("Filename", extractFilenames(docs))
	categoryColumn := entity.NewColumnVarChar("Category", extractCategories(docs))
	embeddingModelColumn := entity.NewColumnVarChar("EmbeddingModel", extractEmbeddingModels(docs))
	summaryColumn := entity.NewColumnVarChar("Summary", extractSummaries(docs))
	metadataColumn := entity.NewColumnVarChar("Metadata", extractMetadata(docs))
	vectorColumn := entity.NewColumnFloatVector("Vector", 1024, extractVectorsDocs(docs))
	// Insert the data
	_, err := m.Instance.Insert(ctx, "documents", "_default", idColumn, contentColumn, linkColumn, filenameColumn,
		categoryColumn, embeddingModelColumn, summaryColumn, metadataColumn, vectorColumn)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %w", err)
	}

	// Flush the collection
	err = m.Instance.Flush(ctx, "documents", false)
	if err != nil {
		return fmt.Errorf("failed to flush documents collection: %w", err)
	}

	return nil
}

// InsertEmbeddings inserts embeddings into the "chunks" collection.
func (m *Client) InsertEmbeddings(ctx context.Context, embeddings []models.Embedding) error {
	idColumn := entity.NewColumnVarChar("ID", extractEmbeddingIDs(embeddings))
	documentIDColumn := entity.NewColumnVarChar("DocumentID", extractDocumentIDs(embeddings))
	vectorColumn := entity.NewColumnFloatVector("Vector", embeddings[0].Dimension, extractVectors(embeddings))
	textChunkColumn := entity.NewColumnVarChar("TextChunk", extractTextChunks(embeddings))
	dimensionColumn := entity.NewColumnInt32("Dimension", extractDimensions(embeddings))
	orderColumn := entity.NewColumnInt32("Order", extractOrders(embeddings))

	_, err := m.Instance.Insert(ctx, "chunks", "", idColumn, documentIDColumn, vectorColumn,
		textChunkColumn, dimensionColumn, orderColumn)
	if err != nil {
		return fmt.Errorf("failed to insert embeddings: %w", err)
	}

	err = m.Instance.Flush(ctx, "chunks", false)
	if err != nil {
		return fmt.Errorf("failed to flush chunks collection: %w", err)
	}

	return nil
}

// GetDocumentByID retrieves a document from the "documents" collection by ID.
func (m *Client) GetDocumentByID(ctx context.Context, id string) (map[string]interface{}, error) {
	collectionName := "documents"
	expr := fmt.Sprintf("ID == '%s'", id)
	projections := []string{"*"} // Fetch all fields

	results, err := m.Instance.Query(ctx, collectionName, nil, expr, projections)
	if err != nil {
		return nil, fmt.Errorf("failed to query document by ID: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("document with ID '%s' not found", id)
	}

	fmt.Println(results)
	return map[string]interface{}{}, nil
}

// GetEmbeddingByID retrieves an embedding from the "chunks" collection by ID.
func (m *Client) GetEmbeddingByID(ctx context.Context, id string) (map[string]interface{}, error) {
	collectionName := "chunks"
	expr := fmt.Sprintf("ID == '%s'", id)
	projections := []string{"*"} // Fetch all fields

	results, err := m.Instance.Query(ctx, collectionName, nil, expr, projections)
	if err != nil {
		return nil, fmt.Errorf("failed to query embedding by ID: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("embedding with ID '%s' not found", id)
	}

	fmt.Println(results)
	return map[string]interface{}{}, nil
}

// DeleteDocument deletes a document from the "documents" collection by ID.
func (m *Client) DeleteDocument(ctx context.Context, id string) error {
	collectionName := "documents"
	partitionName := "_default"
	expr := fmt.Sprintf("ID == '%s'", id)

	err := m.Instance.Delete(ctx, collectionName, partitionName, expr)
	if err != nil {
		return fmt.Errorf("failed to delete document by ID: %w", err)
	}

	return nil
}

// DeleteEmbedding deletes an embedding from the "chunks" collection by ID.
func (m *Client) DeleteEmbedding(ctx context.Context, id string) error {
	collectionName := "chunks"
	partitionName := "_default"
	expr := fmt.Sprintf("ID == '%s'", id)

	err := m.Instance.Delete(ctx, collectionName, partitionName, expr)
	if err != nil {
		return fmt.Errorf("failed to delete embedding by ID: %w", err)
	}

	return nil
}

// Helper functions for extracting data
func extractIDs(docs []models.Document) []string {
	ids := make([]string, len(docs))
	for i, doc := range docs {
		ids[i] = doc.ID
	}
	return ids
}

// extractLinks extracts the "Link" field from the documents.
func extractLinks(docs []models.Document) []string {
	links := make([]string, len(docs))
	for i, doc := range docs {
		links[i] = doc.Link
	}
	return links
}

// extractFilenames extracts the "Filename" field from the documents.
func extractFilenames(docs []models.Document) []string {
	filenames := make([]string, len(docs))
	for i, doc := range docs {
		filenames[i] = doc.Filename
	}
	return filenames
}

// extractCategories extracts the "Category" field from the documents as a comma-separated string.
func extractCategories(docs []models.Document) []string {
	categories := make([]string, len(docs))
	for i, doc := range docs {
		categories[i] = fmt.Sprintf("%v", doc.Category)
	}
	return categories
}

// extractEmbeddingModels extracts the "EmbeddingModel" field from the documents.
func extractEmbeddingModels(docs []models.Document) []string {
	models := make([]string, len(docs))
	for i, doc := range docs {
		models[i] = doc.EmbeddingModel
	}
	return models
}

// extractSummaries extracts the "Summary" field from the documents.
func extractSummaries(docs []models.Document) []string {
	summaries := make([]string, len(docs))
	for i, doc := range docs {
		summaries[i] = doc.Summary
	}
	return summaries
}

// extractMetadata extracts the "Metadata" field from the documents as a JSON string.
func extractMetadata(docs []models.Document) []string {
	metadata := make([]string, len(docs))
	for i, doc := range docs {
		metaBytes, _ := json.Marshal(doc.Metadata)
		metadata[i] = string(metaBytes)
	}
	return metadata
}

func extractContents(docs []models.Document) []string {
	contents := make([]string, len(docs))
	for i, doc := range docs {
		contents[i] = doc.Content
	}
	return contents
}

// extractEmbeddingIDs extracts the "ID" field from the embeddings.
func extractEmbeddingIDs(embeddings []models.Embedding) []string {
	ids := make([]string, len(embeddings))
	for i, embedding := range embeddings {
		ids[i] = embedding.ID
	}
	return ids
}

// extractDocumentIDs extracts the "DocumentID" field from the embeddings.
func extractDocumentIDs(embeddings []models.Embedding) []string {
	documentIDs := make([]string, len(embeddings))
	for i, embedding := range embeddings {
		documentIDs[i] = embedding.DocumentID
	}
	return documentIDs
}

// extractVectors extracts the "Vector" field from the embeddings.
// extractVectors extracts the "Vector" field from the embeddings.
func extractVectors(embeddings []models.Embedding) [][]float32 {
	vectors := make([][]float32, len(embeddings))
	for i, embedding := range embeddings {
		vectors[i] = embedding.Vector // Direct assignment since it's already []float32
	}
	return vectors
}

// extractVectorsDocs extracts the "Vector" field from the documents.
func extractVectorsDocs(docs []models.Document) [][]float32 {
	vectors := make([][]float32, len(docs))
	for i, doc := range docs {
		vectors[i] = doc.Vector // Direct assignment since it's already []float32
	}
	return vectors
}

// extractTextChunks extracts the "TextChunk" field from the embeddings.
func extractTextChunks(embeddings []models.Embedding) []string {
	textChunks := make([]string, len(embeddings))
	for i, embedding := range embeddings {
		textChunks[i] = embedding.TextChunk
	}
	return textChunks
}

// extractDimensions extracts the "Dimension" field from the embeddings.
func extractDimensions(embeddings []models.Embedding) []int32 {
	dimensions := make([]int32, len(embeddings))
	for i, embedding := range embeddings {
		dimensions[i] = int32(embedding.Dimension)
	}
	return dimensions
}

// extractOrders extracts the "Order" field from the embeddings.
func extractOrders(embeddings []models.Embedding) []int32 {
	orders := make([]int32, len(embeddings))
	for i, embedding := range embeddings {
		orders[i] = int32(embedding.Order)
	}
	return orders
}
