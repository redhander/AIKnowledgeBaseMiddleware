package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	infrastructure "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/config"
)

type MilvusClient struct {
	Client         client.Client
	CollectionName string
}

// Option defines a function type for configuring the Milvus client
type Option func(*client.Config)

func NewMilvusClient(cfg infrastructure.MilvusConfig) (*MilvusClient, error) {
	// 创建Milvus客户端
	milvusClient, err := client.NewClient(context.Background(), client.Config{
		Address:  cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %v", err)
	}

	// 检查集合是否存在
	exists, err := milvusClient.HasCollection(context.Background(), cfg.CollectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check collection existence: %v", err)
	}

	// 如果集合不存在，则创建
	if !exists {
		schema := &entity.Schema{
			CollectionName: cfg.CollectionName,
			Description:    "Enterprise knowledge documents",
			AutoID:         false,
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeVarChar,
					PrimaryKey: true,
					AutoID:     false,
					TypeParams: map[string]string{
						"max_length": "64",
					},
				},
				{
					Name:     "content",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						"max_length": "65535",
					},
				},
				{
					Name:     "metadata",
					DataType: entity.FieldTypeJSON,
				},
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": "768", // 根据嵌入模型调整维度
					},
				},
			},
		}

		err := milvusClient.CreateCollection(context.Background(), schema, entity.DefaultShardNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %v", err)
		}

		// 创建向量索引
		index, err := entity.NewIndexIvfFlat(entity.L2, 128)
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %v", err)
		}

		err = milvusClient.CreateIndex(context.Background(), cfg.CollectionName, "vector", index, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %v", err)
		}
	}

	return &MilvusClient{
		Client:         milvusClient,
		CollectionName: cfg.CollectionName,
	}, nil
}

func (mc *MilvusClient) InsertDocuments(docs []document.Document) error {
	// 准备数据
	var ids []string
	var contents []string
	var metadatas []string
	var vectors [][]float32

	for _, doc := range docs {
		ids = append(ids, doc.ID)
		contents = append(contents, doc.Content)
		metadataJSON, err := json.Marshal(doc.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %v", err)
		}
		metadatas = append(metadatas, string(metadataJSON))
		vectors = append(vectors, doc.Vector)
	}

	// 插入数据
	idCol := entity.NewColumnVarChar("id", ids)
	contentCol := entity.NewColumnVarChar("content", contents)
	var metadataBytes [][]byte
	for _, metadata := range metadatas {
		metadataBytes = append(metadataBytes, []byte(metadata))
	}
	metadataCol := entity.NewColumnJSONBytes("metadata", metadataBytes)
	vectorCol := entity.NewColumnFloatVector("vector", 768, vectors)

	_, err := mc.Client.Insert(context.Background(), mc.CollectionName, "", idCol, contentCol, metadataCol, vectorCol)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}

	return nil
}

func (mc *MilvusClient) Search(queryVector []float32, topK int) ([]document.Document, error) {
	// 搜索向量
	sp, _ := entity.NewIndexIvfFlatSearchParam(16)
	results, err := mc.Client.Search(
		context.Background(),
		mc.CollectionName,
		[]string{},
		"",
		[]string{"id", "content", "metadata"},
		[]entity.Vector{entity.FloatVector(queryVector)},
		"vector",
		entity.L2,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %v", err)
	}

	// 处理结果
	var docs []document.Document
	for _, result := range results {
		for i := 0; i < result.ResultCount; i++ {
			id, _ := result.IDs.GetAsString(i)
			content, err := result.Fields.GetColumn("content").GetAsString(i)
			if err != nil {
				return nil, fmt.Errorf("failed to get content: %v", err)

			}
			metadataStr, err := result.Fields.GetColumn("metadata").GetAsString(i)
			if err != nil {
				return nil, fmt.Errorf("failed to get metadata as string: %v", err)
			}
			metadataJSON := []byte(metadataStr)

			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %v", err)
			}
			//todo, 这里屏蔽了代码，需要搞清楚情况
			//score := result.Scores[i]

			docs = append(docs, document.Document{
				ID:      id,
				Content: content,
				//Metadata: metadata,
				//Score:    float64(score),
			})
		}
	}

	return docs, nil
}

// WithAddress sets the Milvus server address
func WithAddress(addr string) Option {
	return func(cfg *client.Config) {
		cfg.Address = addr
	}
}

// WithAuth sets the authentication credentials
func WithAuth(username, password string) Option {
	return func(cfg *client.Config) {
		cfg.Username = username
		cfg.Password = password
	}
}
