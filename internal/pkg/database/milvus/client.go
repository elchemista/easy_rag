package milvus

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Client struct {
	Instance client.Client
}

// InitMilvusClient initializes the Milvus client and returns a wrapper around it.
func NewClient(milvusAddr string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{Address: milvusAddr})
	if err != nil {
		log.Printf("Failed to connect to Milvus: %v", err)
		return nil, err
	}

	client := &Client{Instance: c}

	err = client.EnsureCollections(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// EnsureCollections ensures that the required collections ("documents" and "chunks") exist.
// If they don't exist, it creates them based on the predefined structs.
func (m *Client) EnsureCollections(ctx context.Context) error {
	collections := []struct {
		Name       string
		Schema     *entity.Schema
		IndexField string
		IndexType  string
		MetricType entity.MetricType
		Nlist      int
	}{
		{
			Name:       "documents",
			Schema:     createDocumentSchema(),
			IndexField: "Vector", // Indexing the Vector field for similarity search
			IndexType:  "IVF_FLAT",
			MetricType: entity.L2,
			Nlist:      128, // Number of clusters for IVF_FLAT index
		},
		{
			Name:       "chunks",
			Schema:     createEmbeddingSchema(),
			IndexField: "Vector", // Indexing the Vector field for similarity search
			IndexType:  "IVF_FLAT",
			MetricType: entity.L2,
			Nlist:      128,
		},
	}

	for _, collection := range collections {
		// Ensure the collection exists
		exists, err := m.Instance.HasCollection(ctx, collection.Name)
		if err != nil {
			return fmt.Errorf("failed to check collection existence: %w", err)
		}

		if !exists {
			err := m.Instance.CreateCollection(ctx, collection.Schema, entity.DefaultShardNumber)
			if err != nil {
				return fmt.Errorf("failed to create collection '%s': %w", collection.Name, err)
			}
			log.Printf("Collection '%s' created successfully", collection.Name)
		} else {
			log.Printf("Collection '%s' already exists", collection.Name)
		}

		// Ensure the default partition exists
		hasPartition, err := m.Instance.HasPartition(ctx, collection.Name, "_default")
		if err != nil {
			return fmt.Errorf("failed to check default partition for collection '%s': %w", collection.Name, err)
		}

		if !hasPartition {
			err = m.Instance.CreatePartition(ctx, collection.Name, "_default")
			if err != nil {
				return fmt.Errorf("failed to create default partition for collection '%s': %w", collection.Name, err)
			}
			log.Printf("Default partition created for collection '%s'", collection.Name)
		}

		// Skip index creation if IndexField is empty
		if collection.IndexField == "" {
			continue
		}

		// Ensure the index exists
		log.Printf("Creating index on field '%s' for collection '%s'", collection.IndexField, collection.Name)

		idx, err := entity.NewIndexIvfFlat(collection.MetricType, collection.Nlist)
		if err != nil {
			return fmt.Errorf("failed to create IVF_FLAT index: %w", err)
		}

		err = m.Instance.CreateIndex(ctx, collection.Name, collection.IndexField, idx, false)
		if err != nil {
			return fmt.Errorf("failed to create index on field '%s' for collection '%s': %w", collection.IndexField, collection.Name, err)
		}

		log.Printf("Index created on field '%s' for collection '%s'", collection.IndexField, collection.Name)
	}

	return nil
}

// Helper functions for creating schemas
func createDocumentSchema() *entity.Schema {
	return entity.NewSchema().
		WithName("documents").
		WithDescription("Collection for storing documents").
		WithField(entity.NewField().WithName("ID").WithDataType(entity.FieldTypeVarChar).WithIsPrimaryKey(true).WithMaxLength(512)).
		WithField(entity.NewField().WithName("Content").WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535)).
		WithField(entity.NewField().WithName("Link").WithDataType(entity.FieldTypeVarChar).WithMaxLength(512)).
		WithField(entity.NewField().WithName("Filename").WithDataType(entity.FieldTypeVarChar).WithMaxLength(512)).
		WithField(entity.NewField().WithName("Category").WithDataType(entity.FieldTypeVarChar).WithMaxLength(8048)).
		WithField(entity.NewField().WithName("EmbeddingModel").WithDataType(entity.FieldTypeVarChar).WithMaxLength(256)).
		WithField(entity.NewField().WithName("Summary").WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535)).
		WithField(entity.NewField().WithName("Metadata").WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535)).
		WithField(entity.NewField().WithName("Vector").WithDataType(entity.FieldTypeFloatVector).WithDim(1024)) // bge-m3
}

func createEmbeddingSchema() *entity.Schema {
	return entity.NewSchema().
		WithName("chunks").
		WithDescription("Collection for storing document embeddings").
		WithField(entity.NewField().WithName("ID").WithDataType(entity.FieldTypeVarChar).WithIsPrimaryKey(true).WithMaxLength(512)).
		WithField(entity.NewField().WithName("DocumentID").WithDataType(entity.FieldTypeVarChar).WithMaxLength(512)).
		WithField(entity.NewField().WithName("Vector").WithDataType(entity.FieldTypeFloatVector).WithDim(1024)). // bge-m3
		WithField(entity.NewField().WithName("TextChunk").WithDataType(entity.FieldTypeVarChar).WithMaxLength(65535)).
		WithField(entity.NewField().WithName("Dimension").WithDataType(entity.FieldTypeInt32)).
		WithField(entity.NewField().WithName("Order").WithDataType(entity.FieldTypeInt32))
}

// Close closes the Milvus client connection.
func (m *Client) Close() {
	m.Instance.Close()
}
