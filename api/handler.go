package api

import (
	"fmt"
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
	// upload json list of models.document
	rag := c.Get("Rag").(*rag.Rag)

	var request RequestUpload
	err := c.Bind(&request)

	if err != nil {
		return ErrorHandler(err, c)
	}

	var docs []models.Document
	for _, doc := range request.Docs {
		chunks := textprocessor.CreateChunks(doc.Content)
		summary_chunks := textprocessor.ConcatenateStrings(chunks[:4])

		summary, err := rag.LLM.Generate(summary_chunks)

		if err != nil {
			return ErrorHandler(err, c)
		}

		vectorSum, err := rag.Embeddings.Vectorize(summary)

		if err != nil {
			return ErrorHandler(err, c)
		}

		document := models.Document{
			ID:             uuid.NewString(),
			Content:        "",
			Link:           doc.Link,
			Filename:       doc.Filename,
			Category:       doc.Category,
			EmbeddingModel: rag.Embeddings.GetModel(),
			Summary:        summary,
			Vector:         vectorSum[0],
			Metadata:       doc.Metadata,
		}

		err = rag.Database.SaveDocument(document)

		if err != nil {
			return ErrorHandler(err, c)
		}

		var embeddings []models.Embedding
		for order, chunk := range chunks {
			vectorEmbedding, err := rag.Embeddings.Vectorize(chunk)
			if err != nil {
				return ErrorHandler(err, c)
			}

			embeddings = append(embeddings, models.Embedding{
				ID:         uuid.NewString(),
				DocumentID: document.ID,
				Vector:     vectorEmbedding[0],
				TextChunk:  chunk,
				Dimension:  int64(1024),
				Order:      int64(order),
			})
		}

		err = rag.Database.SaveEmbeddings(embeddings)
		if err != nil {
			return ErrorHandler(err, c)
		}

		docs = append(docs, document)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"docs":    docs,
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

	var chunks []string
	for _, embedding := range embeddings {
		chunks = append(chunks, embedding.TextChunk)
	}

	answer, err := rag.LLM.Generate(fmt.Sprintf("Given the following information: %s \nAnswer the question: %s", textprocessor.ConcatenateStrings(chunks), request.Question))

	if err != nil {
		return ErrorHandler(err, c)
	}

	docs := make([]string, len(embeddings))
	for i, embedding := range embeddings {
		docs[i] = embedding.DocumentID
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": APIVersion,
		"docs":    docs,
		"answer":  answer,
	})
}

func ErrorHandler(err error, c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"error": err.Error(),
	})
}
