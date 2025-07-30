package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
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
		entity.NewColumnFloatVector("vector", 384, [][]float32{doc.Vector}),
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
	//filenameCol := make([]string, len(docs))
	vectorCol := make([][]float32, len(docs))
	metadataCol := make([][]byte, len(docs))
	contentCol := make([]string, len(docs))
	for i, doc := range docs {
		fmt.Printf("Storing document id: %s\n", doc.ID)
		if doc.ID == "" {
			return fmt.Errorf("document at index %d has empty ID", i)
		}
		idCol[i] = doc.ID
		//filenameCol[i] = doc.Metadata.Filename
		vectorCol[i] = doc.Vector
		// Sanitize content to ensure valid UTF-8
		contentCol[i] = sanitizeString(doc.Content)
		metadataCol[i] = []byte(sanitizeString(doc.Metadata.String()))
	}
	//logger.Debugf("metadata=%s", metadataCol)
	// 创建批量插入列
	columns := []entity.Column{
		entity.NewColumnVarChar("id", idCol),
		entity.NewColumnFloatVector("vector", 384, vectorCol),
		entity.NewColumnJSONBytes("metadata", metadataCol),
		entity.NewColumnVarChar("content", contentCol),
	}

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

// Helper function to sanitize strings to valid UTF-8
func sanitizeString(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// If invalid UTF-8 found, replace invalid sequences
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			// Replace invalid UTF-8 with replacement character
			buf = append(buf, []byte("")...)
		} else {
			buf = append(buf, s[i:i+size]...)
		}
		i += size
	}
	return string(buf)
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
		entity.NewColumnFloatVector("vector", 384, [][]float32{doc.Vector}),
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
	outputFields := []string{"id", "metadata", "content"}
	result, err := r.Client.Query(ctx, r.CollectionName, []string{}, expr, outputFields)
	if err != nil {
		return nil, fmt.Errorf("failed to query document: %w", err)
	}

	if len(result) == 0 {
		return nil, nil
	}

	// 直接按顺序获取字段值
	fieldMap := make(map[string]entity.Column)
	for _, col := range result {
		fieldMap[col.Name()] = col
	}
	idCol := fieldMap["id"]
	metadataCol := fieldMap["metadata"]
	contentCol := fieldMap["content"]
	// 获取ID
	idVal, err := idCol.Get(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get id: %w", err)
	}

	// 处理id
	var idStr string
	switch v := idVal.(type) {
	case string:
		idStr = v
	case []byte:
		idStr = string(v)
	default:
		return nil, fmt.Errorf("unexpected type for id field: %T", v)
	}

	// 获取metadata
	metadataVal, err := metadataCol.Get(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	// 解析metadata
	var metadataBytes []byte
	switch v := metadataVal.(type) {
	case string:
		metadataBytes = []byte(sanitizeString(v)) // 添加字符串清理
	case []byte:
		metadataBytes = sanitizeBytes(v) // 添加字节数组清理
	default:
		return nil, fmt.Errorf("unexpected type for metadata field: %T", v)
	}

	var metadata document.Metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata (raw: %s): %w", string(metadataBytes), err)
	}
	// 获取content
	contentVal, err := contentCol.Get(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// 处理content
	var content string
	switch v := contentVal.(type) {
	case string:
		content = v
	case []byte:
		content = string(v)
	default:
		return nil, fmt.Errorf("unexpected type for content field: %T", v)
	}

	doc := &document.Document{
		ID:       idStr,
		Content:  content,
		Metadata: metadata,
	}
	return doc, nil
}

func (r *MilvusDocumentRepository) Search(ctx context.Context, embedding []float32, topK int, filter string) ([]*document.Document, error) {
	// Define the search parameters
	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, fmt.Errorf("failed to create search params: %w", err)
	}

	// Define the vector field name in your Milvus collection
	vectorFieldName := "vector"
	outputFields := []string{"id", "content", "metadata", "vector"}

	results, err := r.Client.Search(
		ctx,
		r.CollectionName,
		[]string{}, // partition names
		filter,     // filter expression
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

// 新增 sanitizeBytes 函数
func sanitizeBytes(b []byte) []byte {
	if utf8.Valid(b) {
		return b
	}

	// 清理无效UTF-8序列
	var buf bytes.Buffer
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError {
			// 跳过无效字符
			i++
		} else {
			buf.WriteRune(r)
			i += size
		}
	}
	return buf.Bytes()
}
