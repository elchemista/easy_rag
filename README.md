# README: RAG System API Documentation

## Overview

This project implements a Retrieval-Augmented Generation (RAG) system for document management, question answering, and embedding-based searches. It uses a vector database and language models to retrieve, store, and generate responses from documents.

---

## API Endpoints

### 1. **List All Documents**

- **Method**: `GET`
- **URL**: `/api/v1/docs`
- **Description**: Retrieve all stored documents.
- **Response**:
    ```json
    {
        "version": "v1",
        "docs": [
            {
                "id": "document_id",
                "filename": "text.txt",
                "summary": "Document summary",
                "metadata": { "key": "value" }
            }
        ]
    }
    ```

---

### 2. **Upload Documents**

- **Method**: `POST`
- **URL**: `/api/v1/upload`
- **Description**: Upload one or more documents for processing.
- **Request Body**:
    ```json
    {
        "docs": [
            {
                "content": "Document content",
                "link": "https://example.com/document",
                "filename": "document.txt",
                "category": "CategoryName",
                "metadata": {
                    "key1": "value1"
                }
            }
        ]
    }
    ```
- **Response**:
    ```json
    {
        "version": "v1",
        "task_id": "unique_task_id",
        "expected_time": "10m",
        "status": "Processing started"
    }
    ```

---

### 3. **Get Document by ID**

- **Method**: `GET`
- **URL**: `/api/v1/doc/{id}`
- **Description**: Retrieve details of a document by its ID.
- **Response**:
    ```json
    {
        "version": "v1",
        "doc": {
            "id": "document_id",
            "content": "Document content",
            "filename": "document.txt",
            "summary": "Document summary",
            "metadata": {
                "key1": "value1"
            }
        }
    }
    ```

---

### 4. **Ask a Question**

- **Method**: `POST`
- **URL**: `/api/v1/ask`
- **Description**: Ask a question based on stored documents.
- **Request Body**:
    ```json
    {
        "question": "What is ISO 27001?"
    }
    ```
- **Response**:
    ```json
    {
        "version": "v1",
        "docs": ["document_id_1", "document_id_2"],
        "answer": "ISO 27001 is an international information technology standard..."
    }
    ```

---

### 5. **Delete Document**

- **Method**: `DELETE`
- **URL**: `/api/v1/doc/{id}`
- **Description**: Delete a document by its ID.
- **Response**:
    ```json
    {
        "version": "v1",
        "docs": null
    }
    ```

---

## Data Structures

### **Document**

Represents a stored document.

```go
type Document struct {
    ID             string            `json:"id" milvus:"ID"`                          // Unique identifier
    Content        string            `json:"content" milvus:"Content"`                // Document content (stored as chunks)
    Link           string            `json:"link" milvus:"Link"`                      // Source link
    Filename       string            `json:"filename" milvus:"Filename"`              // Document filename
    Category       string            `json:"category" milvus:"Category"`              // Document category
    EmbeddingModel string            `json:"embedding_model" milvus:"EmbeddingModel"` // Embedding model used
    Summary        string            `json:"summary" milvus:"Summary"`                // Summary of the document
    Metadata       map[string]string `json:"metadata" milvus:"Metadata"`              // Metadata
    Vector         []float32         `json:"vector" milvus:"Vector"`                  // Embedding vector
}
```

### **Embedding**

Represents vector embeddings of document chunks.

```go
type Embedding struct {
    ID         string    `json:"id" milvus:"ID"`                  // Unique identifier
    DocumentID string    `json:"document_id" milvus:"DocumentID"` // Associated document ID
    Vector     []float32 `json:"vector" milvus:"Vector"`          // Embedding vector
    TextChunk  string    `json:"text_chunk" milvus:"TextChunk"`   // Text chunk of the document
    Dimension  int64     `json:"dimension" milvus:"Dimension"`    // Vector dimensionality
    Order      int64     `json:"order" milvus:"Order"`            // Chunk order
    Score      float32   `json:"score"`                           // Search relevance score
}
```

---

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/elchemista/easy_rag.git
   cd easy_rag
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

4. Access the API at `http://localhost:4002`.

---

## Development Notes

- **LLM Integration**: The system supports multiple LLM services (e.g., OpenAI, Ollama) via the `LLM` interface.
- **Database Flexibility**: The project allows switching between different databases (e.g., Milvus, MongoDB) by implementing the `Database` interface.
- **Chunking and Vectorization**:
  - Documents are chunked for efficient embedding and search.
  - Each chunk is vectorized and stored in the database.
