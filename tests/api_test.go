package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elchemista/easy_rag/internal/api"
	"github.com/elchemista/easy_rag/internal/models"
	"github.com/elchemista/easy_rag/internal/pkg/rag"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Example: Test UploadHandler
func TestUploadHandler(t *testing.T) {
	e := echo.New()

	// Create a mock for the LLM, Embeddings, and Database
	mockLLM := new(MockLLMService)
	mockEmbeddings := new(MockEmbeddingsService)
	mockDB := new(MockDatabase)

	// Setup the Rag object
	r := &rag.Rag{
		LLM:        mockLLM,
		Embeddings: mockEmbeddings,
		Database:   mockDB,
	}

	// We expect calls to these mocks in the background goroutine, for each document.

	// The request body
	requestBody := api.RequestUpload{
		Docs: []api.UploadDoc{
			{
				Content:  "Test document content",
				Link:     "http://example.com/doc",
				Filename: "doc1.txt",
				Category: "TestCategory",
				Metadata: map[string]string{"Author": "Me"},
			},
		},
	}

	// Convert requestBody to JSON
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a ResponseRecorder
	rec := httptest.NewRecorder()

	// New echo context
	c := e.NewContext(req, rec)
	// Set the rag object in context
	c.Set("Rag", r)

	// Because the UploadHandler spawns a goroutine, we only test the immediate HTTP response.
	// We can still set expectations for the calls that happen in the goroutine to ensure they're invoked.
	// For example, we expect the summary to be generated, so:

	testSummary := "Test summary from LLM"
	mockLLM.On("Generate", mock.Anything).Return(testSummary, nil).Maybe() // .Maybe() because the concurrency might not complete by the time we assert

	// The embedding vector returned from the embeddings service
	testVector := [][]float32{{0.1, 0.2, 0.3, 0.4}}

	// We'll mock calls to Vectorize() for summary and each chunk
	mockEmbeddings.On("Vectorize", mock.AnythingOfType("string")).Return(testVector, nil).Maybe()

	// The database SaveDocument / SaveEmbeddings calls
	mockDB.On("SaveDocument", mock.AnythingOfType("models.Document")).Return(nil).Maybe()
	mockDB.On("SaveEmbeddings", mock.AnythingOfType("[]models.Embedding")).Return(nil).Maybe()

	// Invoke the handler
	err := api.UploadHandler(c)

	// Check no immediate errors
	assert.NoError(t, err)

	// Check the response
	assert.Equal(t, http.StatusAccepted, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	// We expect certain fields in the JSON response
	assert.Equal(t, "v1", resp["version"])
	assert.NotEmpty(t, resp["task_id"])
	assert.Equal(t, "Processing started", resp["status"])

	// Typically, you might want to wait or do more advanced concurrency checks if you want to test
	// the logic in the goroutine, but that goes beyond a simple unit test.
	// The background process can be tested more thoroughly in integration or end-to-end tests.

	// Optionally assert that our mocks were called
	mockLLM.AssertExpectations(t)
	mockEmbeddings.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// Example: Test ListAllDocsHandler
func TestListAllDocsHandler(t *testing.T) {
	e := echo.New()

	mockLLM := new(MockLLMService)
	mockEmbeddings := new(MockEmbeddingsService)
	mockDB := new(MockDatabase)

	r := &rag.Rag{
		LLM:        mockLLM,
		Embeddings: mockEmbeddings,
		Database:   mockDB,
	}

	// Mock data
	doc1 := models.Document{
		ID:       uuid.NewString(),
		Filename: "doc1.txt",
		Summary:  "summary doc1",
	}
	doc2 := models.Document{
		ID:       uuid.NewString(),
		Filename: "doc2.txt",
		Summary:  "summary doc2",
	}
	docs := []models.Document{doc1, doc2}

	// Expect the database to return the docs
	mockDB.On("ListDocuments").Return(docs, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("Rag", r)

	err := api.ListAllDocsHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "v1", resp["version"])

	// The "docs" field should match the ones we returned
	docsInterface, ok := resp["docs"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, docsInterface, 2)

	// Verify mocks
	mockDB.AssertExpectations(t)
}

// Example: Test GetDocHandler
func TestGetDocHandler(t *testing.T) {
	e := echo.New()

	mockLLM := new(MockLLMService)
	mockEmbeddings := new(MockEmbeddingsService)
	mockDB := new(MockDatabase)

	r := &rag.Rag{
		LLM:        mockLLM,
		Embeddings: mockEmbeddings,
		Database:   mockDB,
	}

	// Mock a single doc
	docID := "123"
	testDoc := models.Document{
		ID:       docID,
		Filename: "doc3.txt",
		Summary:  "summary doc3",
	}

	mockDB.On("GetDocument", docID).Return(testDoc, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/doc/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// set path param
	c.SetParamNames("id")
	c.SetParamValues(docID)
	c.Set("Rag", r)

	err := api.GetDocHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "v1", resp["version"])

	docInterface := resp["doc"].(map[string]interface{})
	assert.Equal(t, "doc3.txt", docInterface["filename"])

	// Verify mocks
	mockDB.AssertExpectations(t)
}

// Example: Test AskDocHandler
func TestAskDocHandler(t *testing.T) {
	e := echo.New()

	mockLLM := new(MockLLMService)
	mockEmbeddings := new(MockEmbeddingsService)
	mockDB := new(MockDatabase)

	r := &rag.Rag{
		LLM:        mockLLM,
		Embeddings: mockEmbeddings,
		Database:   mockDB,
	}

	// 1) We expect to Vectorize the question
	question := "What is the summary of doc?"
	questionVector := [][]float32{{0.5, 0.2, 0.1}}
	mockEmbeddings.On("Vectorize", question).Return(questionVector, nil)

	// 2) We expect a DB search
	emb := []models.Embedding{
		{
			ID:         "emb1",
			DocumentID: "doc123",
			TextChunk:  "Relevant content chunk",
			Score:      0.99,
		},
	}
	mockDB.On("Search", questionVector).Return(emb, nil)

	// 3) We expect the LLM to generate an answer from the chunk
	generatedAnswer := "Here is an answer from the chunk"
	// The prompt we pass is something like: "Given the following information: chunk ... Answer the question: question"
	mockLLM.On("Generate", mock.AnythingOfType("string")).Return(generatedAnswer, nil)

	// Prepare request
	reqBody := api.RequestQuestion{
		Question: question,
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ask", bytes.NewReader(reqBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("Rag", r)

	// Execute
	err := api.AskDocHandler(c)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "v1", resp["version"])
	assert.Equal(t, generatedAnswer, resp["answer"])

	// The docs field should have the docID "doc123"
	docsInterface := resp["docs"].([]interface{})
	assert.Len(t, docsInterface, 1)
	assert.Equal(t, "doc123", docsInterface[0])

	// Verify mocks
	mockLLM.AssertExpectations(t)
	mockEmbeddings.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

// Example: Test DeleteDocHandler
func TestDeleteDocHandler(t *testing.T) {
	e := echo.New()
	mockLLM := new(MockLLMService)
	mockEmbeddings := new(MockEmbeddingsService)
	mockDB := new(MockDatabase)

	r := &rag.Rag{
		LLM:        mockLLM,
		Embeddings: mockEmbeddings,
		Database:   mockDB,
	}

	docID := "abc"
	mockDB.On("DeleteDocument", docID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/doc/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(docID)
	c.Set("Rag", r)

	err := api.DeleteDocHandler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "v1", resp["version"])

	// docs should be nil according to DeleteDocHandler
	assert.Nil(t, resp["docs"])

	// Verify mocks
	mockDB.AssertExpectations(t)
}
