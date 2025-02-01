package milvus

import (
	"context"
	"fmt"
	"sort"

	"github.com/elchemista/easy_rag/internal/models"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
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
	vectorColumn := entity.NewColumnFloatVector("Vector", 1024, extractVectors(embeddings))
	textChunkColumn := entity.NewColumnVarChar("TextChunk", extractTextChunks(embeddings))
	dimensionColumn := entity.NewColumnInt32("Dimension", extractDimensions(embeddings))
	orderColumn := entity.NewColumnInt32("Order", extractOrders(embeddings))

	_, err := m.Instance.Insert(ctx, "chunks", "_default", idColumn, documentIDColumn, vectorColumn,
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
	projections := []string{"ID", "Content", "Link", "Filename", "Category", "EmbeddingModel", "Summary", "Metadata"} // Fetch all fields

	results, err := m.Instance.Query(ctx, collectionName, nil, expr, projections)
	if err != nil {
		return nil, fmt.Errorf("failed to query document by ID: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("document with ID '%s' not found", id)
	}

	mp, err := transformResultSet(results, "ID", "Content", "Link", "Filename", "Category", "EmbeddingModel", "Summary", "Metadata")

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	// convert metadata to map
	mp[0]["Metadata"] = convertToMetadata(mp[0]["Metadata"].(string))

	return mp[0], err
}

// GetAllDocuments retrieves all documents from the "documents" collection.
func (m *Client) GetAllDocuments(ctx context.Context) ([]models.Document, error) {
	collectionName := "documents"
	projections := []string{"ID", "Content", "Link", "Filename", "Category", "EmbeddingModel", "Summary", "Metadata"} // Fetch all fields
	expr := ""

	rs, err := m.Instance.Query(ctx, collectionName, nil, expr, projections, client.WithLimit(1000))
	if err != nil {
		return nil, fmt.Errorf("failed to query all documents: %w", err)
	}

	if len(rs) == 0 {
		return nil, fmt.Errorf("no documents found in the collection")
	}

	results, err := transformResultSet(rs, "ID", "Content", "Link", "Filename", "Category", "EmbeddingModel", "Summary", "Metadata")

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal all documents: %w", err)
	}

	var docs []models.Document = make([]models.Document, len(results))
	for i, result := range results {
		docs[i] = models.Document{
			ID:             result["ID"].(string),
			Content:        result["Content"].(string),
			Link:           result["Link"].(string),
			Filename:       result["Filename"].(string),
			Category:       result["Category"].(string),
			EmbeddingModel: result["EmbeddingModel"].(string),
			Summary:        result["Summary"].(string),
			Metadata:       convertToMetadata(results[0]["Metadata"].(string)),
		}
	}

	return docs, nil
}

// GetAllEmbeddingByDocID retrieves all embeddings linked to a specific DocumentID from the "chunks" collection.
func (m *Client) GetAllEmbeddingByDocID(ctx context.Context, documentID string) ([]models.Embedding, error) {
	collectionName := "chunks"
	projections := []string{"ID", "DocumentID", "TextChunk", "Order"} // Fetch all fields
	expr := fmt.Sprintf("DocumentID == '%s'", documentID)

	rs, err := m.Instance.Query(ctx, collectionName, nil, expr, projections, client.WithLimit(1000))

	if err != nil {
		return nil, fmt.Errorf("failed to query embeddings by DocumentID: %w", err)
	}

	if rs.Len() == 0 {
		return nil, fmt.Errorf("no embeddings found for DocumentID '%s'", documentID)
	}

	results, err := transformResultSet(rs, "ID", "DocumentID", "TextChunk", "Order")

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal all documents: %w", err)
	}

	var embeddings []models.Embedding = make([]models.Embedding, rs.Len())

	for i, result := range results {
		embeddings[i] = models.Embedding{
			ID:         result["ID"].(string),
			DocumentID: result["DocumentID"].(string),
			TextChunk:  result["TextChunk"].(string),
			Order:      result["Order"].(int64),
		}
	}

	return embeddings, nil
}

func (m *Client) Search(ctx context.Context, vectors [][]float32, topK int) ([]models.Embedding, error) {
	const (
		collectionName = "chunks"
		vectorDim      = 1024 // Replace with your actual vector dimension
	)
	projections := []string{"ID", "DocumentID", "TextChunk", "Order"}
	metricType := entity.L2 // Default metric type

	// Validate and convert input vectors
	searchVectors, err := validateAndConvertVectors(vectors, vectorDim)
	if err != nil {
		return nil, err
	}

	// Set search parameters
	searchParams, err := entity.NewIndexIvfFlatSearchParam(16) // 16 is the number of clusters for IVF_FLAT index
	if err != nil {
		return nil, fmt.Errorf("failed to create search params: %w", err)
	}

	// Perform the search
	searchResults, err := m.Instance.Search(ctx, collectionName, nil, "", projections, searchVectors, "Vector", metricType, topK, searchParams, client.WithLimit(10))
	if err != nil {
		return nil, fmt.Errorf("failed to search collection: %w", err)
	}

	// Process search results
	embeddings, err := processSearchResults(searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to process search results: %w", err)
	}

	return embeddings, nil
}

// validateAndConvertVectors validates vector dimensions and converts them to Milvus-compatible format.
func validateAndConvertVectors(vectors [][]float32, expectedDim int) ([]entity.Vector, error) {
	searchVectors := make([]entity.Vector, len(vectors))
	for i, vector := range vectors {
		if len(vector) != expectedDim {
			return nil, fmt.Errorf("vector dimension mismatch: expected %d, got %d", expectedDim, len(vector))
		}
		searchVectors[i] = entity.FloatVector(vector)
	}
	return searchVectors, nil
}

// processSearchResults transforms and aggregates the search results into embeddings and sorts by score.
func processSearchResults(results []client.SearchResult) ([]models.Embedding, error) {
	var embeddings []models.Embedding

	for _, result := range results {
		for i := 0; i < result.ResultCount; i++ {
			embeddingMap, err := transformSearchResultSet(result, "ID", "DocumentID", "TextChunk", "Order")
			if err != nil {
				return nil, fmt.Errorf("failed to transform search result set: %w", err)
			}

			for _, embedding := range embeddingMap {
				embeddings = append(embeddings, models.Embedding{
					ID:         embedding["ID"].(string),
					DocumentID: embedding["DocumentID"].(string),
					TextChunk:  embedding["TextChunk"].(string),
					Order:      embedding["Order"].(int64), // Assuming 'Order' is a float64 type
					Score:      embedding["Score"].(float32),
				})
			}
		}
	}

	// Sort embeddings by score in descending order (higher is better)
	sort.Slice(embeddings, func(i, j int) bool {
		return embeddings[i].Score > embeddings[j].Score
	})

	return embeddings, nil
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
	expr := fmt.Sprintf("DocumentID == '%s'", id)

	err := m.Instance.Delete(ctx, collectionName, partitionName, expr)
	if err != nil {
		return fmt.Errorf("failed to delete embedding by DocumentID: %w", err)
	}

	return nil
}
