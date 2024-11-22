package milvus

import (
	"encoding/json"
	"fmt"

	"github.com/elchemista/easy_rag/internal/models"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

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

func convertToMetadata(metadata string) map[string]string {
	var metadataMap map[string]string
	json.Unmarshal([]byte(metadata), &metadataMap)
	return metadataMap
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

func transformResultSet(rs client.ResultSet, outputFields ...string) ([]map[string]interface{}, error) {
	if rs == nil || rs.Len() == 0 {
		return nil, fmt.Errorf("empty result set")
	}

	results := []map[string]interface{}{}

	for i := 0; i < rs.Len(); i++ { // Iterate through rows
		row := map[string]interface{}{}

		for _, fieldName := range outputFields {
			column := rs.GetColumn(fieldName)
			if column == nil {
				return nil, fmt.Errorf("column %s does not exist in result set", fieldName)
			}

			switch column.Type() {
			case entity.FieldTypeInt64:
				value, err := column.GetAsInt64(i)
				if err != nil {
					return nil, fmt.Errorf("error getting int64 value for column %s, row %d: %w", fieldName, i, err)
				}
				row[fieldName] = value

			case entity.FieldTypeInt32:
				value, err := column.GetAsInt64(i)
				if err != nil {
					return nil, fmt.Errorf("error getting int64 value for column %s, row %d: %w", fieldName, i, err)
				}
				row[fieldName] = value

			case entity.FieldTypeFloat:
				value, err := column.GetAsDouble(i)
				if err != nil {
					return nil, fmt.Errorf("error getting float value for column %s, row %d: %w", fieldName, i, err)
				}
				row[fieldName] = value

			case entity.FieldTypeDouble:
				value, err := column.GetAsDouble(i)
				if err != nil {
					return nil, fmt.Errorf("error getting double value for column %s, row %d: %w", fieldName, i, err)
				}
				row[fieldName] = value

			case entity.FieldTypeVarChar:
				value, err := column.GetAsString(i)
				if err != nil {
					return nil, fmt.Errorf("error getting string value for column %s, row %d: %w", fieldName, i, err)
				}
				row[fieldName] = value

			default:
				return nil, fmt.Errorf("unsupported field type for column %s", fieldName)
			}
		}

		results = append(results, row)
	}

	return results, nil
}

func transformSearchResultSet(rs client.SearchResult, outputFields ...string) ([]map[string]interface{}, error) {
	if rs.ResultCount == 0 {
		return nil, fmt.Errorf("empty result set")
	}

	result := make([]map[string]interface{}, rs.ResultCount)

	for i := 0; i < rs.ResultCount; i++ { // Iterate through rows
		result[i] = make(map[string]interface{})
		for _, fieldName := range outputFields {
			column := rs.Fields.GetColumn(fieldName)
			result[i]["Score"] = rs.Scores[i]

			if column == nil {
				return nil, fmt.Errorf("column %s does not exist in result set", fieldName)
			}

			switch column.Type() {
			case entity.FieldTypeInt64:
				value, err := column.GetAsInt64(i)
				if err != nil {
					return nil, fmt.Errorf("error getting int64 value for column %s, row %d: %w", fieldName, i, err)
				}
				result[i][fieldName] = value

			case entity.FieldTypeInt32:
				value, err := column.GetAsInt64(i)
				if err != nil {
					return nil, fmt.Errorf("error getting int64 value for column %s, row %d: %w", fieldName, i, err)
				}
				result[i][fieldName] = value

			case entity.FieldTypeFloat:
				value, err := column.GetAsDouble(i)
				if err != nil {
					return nil, fmt.Errorf("error getting float value for column %s, row %d: %w", fieldName, i, err)
				}
				result[i][fieldName] = value

			case entity.FieldTypeDouble:
				value, err := column.GetAsDouble(i)
				if err != nil {
					return nil, fmt.Errorf("error getting double value for column %s, row %d: %w", fieldName, i, err)
				}
				result[i][fieldName] = value

			case entity.FieldTypeVarChar:
				value, err := column.GetAsString(i)
				if err != nil {
					return nil, fmt.Errorf("error getting string value for column %s, row %d: %w", fieldName, i, err)
				}
				result[i][fieldName] = value

			default:
				return nil, fmt.Errorf("unsupported field type for column %s", fieldName)
			}
		}
	}

	return result, nil
}
