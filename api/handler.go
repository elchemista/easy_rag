package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elchemista/easy_rag/internal/models"
	"github.com/elchemista/easy_rag/internal/pkg/rag"
	"github.com/elchemista/easy_rag/internal/pkg/textprocessor"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UploadDoc struct {
	Content  string            `json:"content"`
	Link     string            `json:"link"`
	Filename string            `json:"filename"`
	Category string            `json:"category"`
	Metadata map[string]string `json:"metadata"`
}

type RequestUpload struct {
	Docs []UploadDoc `json:"docs"`
}

type RequestQuestion struct {
	Question string `json:"question"`
}

type ResposeQuestion struct {
	Version string            `json:"version"`
	Docs    []models.Document `json:"docs"`
	Answer  string            `json:"answer"`
}

func UploadHandler(c echo.Context) error {
	// Retrieve the RAG instance from context
	rag := c.Get("Rag").(*rag.Rag)

	var request RequestUpload
	if err := c.Bind(&request); err != nil {
		return ErrorHandler(err, c)
	}

	// Generate a unique task ID
	taskID := uuid.NewString()

	// Launch the upload process in a separate goroutine
	go func(taskID string, request RequestUpload) {
		log.Printf("Task %s: started processing", taskID)
		defer log.Printf("Task %s: completed processing", taskID)

		var docs []models.Document

		for idx, doc := range request.Docs {
			// Generate a unique ID for each document
			docID := uuid.NewString()
			log.Printf("Task %s: processing document %d with generated ID %s (filename: %s)", taskID, idx, docID, doc.Filename)

			// Step 1: Create chunks from document content
			chunks := textprocessor.CreateChunks(doc.Content)
			log.Printf("Task %s: created %d chunks for document %s", taskID, len(chunks), docID)

			// Step 2: Generate summary for the document
			var summaryChunks string
			if len(chunks) < 4 {
				summaryChunks = doc.Content
			} else {
				summaryChunks = textprocessor.ConcatenateStrings(chunks[:3])
			}

			log.Printf("Task %s: generating summary for document %s", taskID, docID)
			summary, err := rag.LLM.Generate(fmt.Sprintf("Give me only summary of the following text: %s", summaryChunks))
			if err != nil {
				log.Printf("Task %s: error generating summary for document %s: %v", taskID, docID, err)
				return
			}
			log.Printf("Task %s: generated summary for document %s", taskID, docID)

			// Step 3: Vectorize the summary
			log.Printf("Task %s: vectorizing summary for document %s", taskID, docID)
			vectorSum, err := rag.Embeddings.Vectorize(summary)
			if err != nil {
				log.Printf("Task %s: error vectorizing summary for document %s: %v", taskID, docID, err)
				return
			}
			log.Printf("Task %s: vectorized summary for document %s", taskID, docID)

			// Step 4: Save the document
			document := models.Document{
				ID:             docID, // Use generated ID
				Content:        "",
				Link:           doc.Link,
				Filename:       doc.Filename,
				Category:       doc.Category,
				EmbeddingModel: rag.Embeddings.GetModel(),
				Summary:        summary,
				Vector:         vectorSum[0],
				Metadata:       doc.Metadata,
			}
			log.Printf("Task %s: saving document %s", taskID, docID)
			if err := rag.Database.SaveDocument(document); err != nil {
				log.Printf("Task %s: error saving document %s: %v", taskID, docID, err)
				return
			}
			log.Printf("Task %s: saved document %s", taskID, docID)

			// Step 5: Process and save embeddings for each chunk
			var embeddings []models.Embedding
			for order, chunk := range chunks {
				log.Printf("Task %s: vectorizing chunk %d for document %s", taskID, order, docID)
				vectorEmbedding, err := rag.Embeddings.Vectorize(chunk)
				if err != nil {
					log.Printf("Task %s: error vectorizing chunk %d for document %s: %v", taskID, order, docID, err)
					return
				}
				log.Printf("Task %s: vectorized chunk %d for document %s", taskID, order, docID)

				embedding := models.Embedding{
					ID:         uuid.NewString(),
					DocumentID: docID,
					Vector:     vectorEmbedding[0],
					TextChunk:  chunk,
					Dimension:  int64(1024),
					Order:      int64(order),
				}
				embeddings = append(embeddings, embedding)
			}

			log.Printf("Task %s: saving %d embeddings for document %s", taskID, len(embeddings), docID)
			if err := rag.Database.SaveEmbeddings(embeddings); err != nil {
				log.Printf("Task %s: error saving embeddings for document %s: %v", taskID, docID, err)
				return
			}
			log.Printf("Task %s: saved embeddings for document %s", taskID, docID)

			docs = append(docs, document)
		}
	}(taskID, request)

	// Return the task ID and expected completion time
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"version":       APIVersion,
		"task_id":       taskID,
		"expected_time": "10m",
		"status":        "Processing started",
	})
}

func ListAllDocsHandler(c echo.Context) error {
	rag := c.Get("Rag").(*rag.Rag)
	docs, err := rag.Database.ListDocuments()
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"docs":    docs,
	})
}

func GetDocHandler(c echo.Context) error {
	rag := c.Get("Rag").(*rag.Rag)
	id := c.Param("id")
	doc, err := rag.Database.GetDocument(id)
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"doc":     doc,
	})
}

func AskDocHandler(c echo.Context) error {
	rag := c.Get("Rag").(*rag.Rag)

	var request RequestQuestion
	err := c.Bind(&request)

	if err != nil {
		return ErrorHandler(err, c)
	}

	questionV, err := rag.Embeddings.Vectorize(request.Question)

	if err != nil {
		return ErrorHandler(err, c)
	}

	embeddings, err := rag.Database.Search(questionV)

	if err != nil {
		return ErrorHandler(err, c)
	}

	if len(embeddings) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"version": APIVersion,
			"docs":    nil,
			"answer":  "Don't found any relevant documents",
		})
	}

	answer, err := rag.LLM.Generate(fmt.Sprintf("Given the following information: %s \nAnswer the question: %s", embeddings[0].TextChunk, request.Question))

	if err != nil {
		return ErrorHandler(err, c)
	}

	// Use a map to track unique DocumentIDs
	docSet := make(map[string]struct{})
	for _, embedding := range embeddings {
		docSet[embedding.DocumentID] = struct{}{}
	}

	// Convert the map keys to a slice
	docs := make([]string, 0, len(docSet))
	for docID := range docSet {
		docs = append(docs, docID)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"docs":    docs,
		"answer":  answer,
	})
}

func DeleteDocHandler(c echo.Context) error {
	rag := c.Get("Rag").(*rag.Rag)
	id := c.Param("id")
	err := rag.Database.DeleteDocument(id)
	if err != nil {
		return ErrorHandler(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"docs":    nil,
	})
}

func ErrorHandler(err error, c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": err.Error(),
	})
}
