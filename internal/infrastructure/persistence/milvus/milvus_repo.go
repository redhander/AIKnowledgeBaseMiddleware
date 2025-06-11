package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
)

type MilvusDocumentRepository struct {
	Client         client.Client
	CollectionName string
}

func NewMilvusDocumentRepository(milvusClient *MilvusClient, collectionName string) *MilvusDocumentRepository {
	return &MilvusDocumentRepository{
		Client:         milvusClient.Client,
		CollectionName: collectionName,
	}
}

func (r *MilvusDocumentRepository) Store(ctx context.Context, doc *document.Document) error {
	// 准备Milvus插入数据
	columns := []entity.Column{
		entity.NewColumnString("id", []string{doc.ID}),
		entity.NewColumnVarChar("filename", []string{doc.Metadata.Filename}),
		entity.NewColumnFloatVector("vector", 768, [][]float32{doc.Vector}),
		entity.NewColumnJSONBytes("metadata", [][]byte{[]byte(doc.Metadata.String())}),
		entity.NewColumnString("content", []string{doc.Content}),
	}
	// 插入数据到Milvus
	_, err := r.Client.Insert(ctx, r.CollectionName, "", columns...)
	if err != nil {
		return fmt.Errorf("failed to insert document into Milvus: %w", err)
	}

	// 刷新集合确保数据可搜索
	err = r.Client.Flush(ctx, r.CollectionName, false)
	if err != nil {
		return fmt.Errorf("failed to flush Milvus collection: %w", err)
	}
	return nil
}

func (r *MilvusDocumentRepository) StoreBatch(ctx context.Context, docs []*document.Document) error {
	if len(docs) == 0 {
		return nil
	}

	// 准备批量插入数据
	idCol := make([]string, len(docs))
	filenameCol := make([]string, len(docs))
	vectorCol := make([][]float32, len(docs))
	metadataCol := make([][]byte, len(docs))
	contentCol := make([]string, len(docs))

	for i, doc := range docs {
		idCol[i] = doc.ID
		filenameCol[i] = doc.Metadata.Filename
		vectorCol[i] = doc.Vector
		metadataCol[i] = []byte(doc.Metadata.String())
		contentCol[i] = doc.Content
	}

	// 创建批量插入列
	columns := []entity.Column{
		entity.NewColumnString("id", idCol),
		entity.NewColumnVarChar("filename", filenameCol),
		entity.NewColumnFloatVector("vector", 768, vectorCol),
		entity.NewColumnJSONBytes("metadata", metadataCol),
		entity.NewColumnString("content", contentCol),
	}

	// 批量插入数据到Milvus
	_, err := r.Client.Insert(ctx, r.CollectionName, "", columns...)
	if err != nil {
		return fmt.Errorf("failed to batch insert documents into Milvus: %w", err)
	}

	// 刷新集合确保数据可搜索
	err = r.Client.Flush(ctx, r.CollectionName, false)
	if err != nil {
		return fmt.Errorf("failed to flush Milvus collection: %w", err)
	}

	return nil
}

// Implement all required methods
func (r *MilvusDocumentRepository) Save(ctx context.Context, doc *document.Document) error {
	// 先尝试查找文档是否存在
	existingDoc, err := r.FindByID(ctx, doc.ID)
	if err != nil {
		return fmt.Errorf("failed to check document existence: %w", err)
	}

	// 如果文档不存在，直接存储
	if existingDoc == nil {
		return r.Store(ctx, doc)
	}

	// 文档存在则更新
	columns := []entity.Column{
		entity.NewColumnString("id", []string{doc.ID}),
		entity.NewColumnVarChar("filename", []string{doc.Metadata.Filename}),
		entity.NewColumnFloatVector("vector", 768, [][]float32{doc.Vector}),
		entity.NewColumnJSONBytes("metadata", [][]byte{[]byte(doc.Metadata.String())}),
		entity.NewColumnString("content", []string{doc.Content}),
	}
	// 执行更新操作
	_, err = r.Client.Upsert(ctx, r.CollectionName, "", columns...)
	if err != nil {
		return fmt.Errorf("failed to update document in Milvus: %w", err)
	}

	// 刷新集合确保数据可搜索
	err = r.Client.Flush(ctx, r.CollectionName, false)
	if err != nil {
		return fmt.Errorf("failed to flush Milvus collection: %w", err)
	}

	return nil
}

func (r *MilvusDocumentRepository) FindByID(ctx context.Context, id string) (*document.Document, error) {
	expr := fmt.Sprintf("id == \"%s\"", id)
	outputFields := []string{"id", "filename", "content", "metadata", "vector"}

	result, err := r.Client.Query(
		ctx,
		r.CollectionName,
		[]string{},
		expr,
		outputFields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query document from Milvus: %w", err)
	}

	if result == nil || len(result) == 0 {
		return nil, nil
	}
	row := result[0]

	idVal, err := row.Get(0) // Assuming 'id' is the first field
	if err != nil {
		return nil, fmt.Errorf("failed to get id: %w", err)
	}

	// _, err = row.Get(1) // Assuming 'filename' is the second field
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get filename: %w", err)
	// }

	vectorVal, err := row.Get(2) // Assuming 'vector' is the third field
	if err != nil {
		return nil, fmt.Errorf("failed to get vector: %w", err)
	}
	logger.Info("vectorVal: %v", vectorVal)
	metadataVal, err := row.Get(3) // Assuming 'metadata' is the fourth field
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}
	logger.Info("metadataVal: %v", metadataVal)
	contentVal, err := row.Get(4) // Assuming 'content' is the fifth field
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}
	logger.Info("contentVal: %v", contentVal)
	// Unmarshal metadata JSON
	var metadata document.Metadata
	if err := json.Unmarshal(metadataVal.([]byte), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	doc := &document.Document{
		ID:      idVal.(string),
		Content: contentVal.(string),
		Vector:  vectorVal.([]float32),
		Metadata: document.Metadata{
			Filename:     metadata.Filename,
			Custom:       metadata.Custom,
			UploadTime:   metadata.UploadTime,
			OriginalFile: metadata.OriginalFile,
			ContentType:  metadata.ContentType,
			Size:         metadata.Size,
		},
	}
	return doc, nil
}

func (r *MilvusDocumentRepository) Search(ctx context.Context, embedding []float32, topK int) ([]*document.Document, error) {
	// Define the search parameters
	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, fmt.Errorf("failed to create search params: %w", err)
	}

	// Define the vector field name in your Milvus collection
	vectorFieldName := "vector"
	outputFields := []string{"id", "filename", "content", "metadata", "vector"}

	results, err := r.Client.Search(
		ctx,
		r.CollectionName,
		[]string{}, // partition names
		"",         // filter expression
		outputFields,
		[]entity.Vector{entity.FloatVector(embedding)},
		vectorFieldName,
		entity.L2, // distance metric
		topK,
		searchParam,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform Milvus vector search: %w", err)
	}

	// Process search results
	if len(results) == 0 {
		return []*document.Document{}, nil
	}

	// Convert IDs column to proper type
	idCol, ok := results[0].IDs.(*entity.ColumnVarChar)
	if !ok {
		return nil, fmt.Errorf("failed to convert IDs column to VarChar type")
	}

	// Extract result IDs and build output documents
	var docs []*document.Document
	for i := 0; i < idCol.Len(); i++ {
		id, err := idCol.ValueByIdx(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get ID at index %d: %w", i, err)
		}

		distance := results[0].Scores[i]

		// Fetch full document by ID
		doc, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get document by ID %s: %w", id, err)
		}

		if doc != nil {
			doc.Score = distance
			docs = append(docs, doc)
		}
	}

	return docs, nil
}
