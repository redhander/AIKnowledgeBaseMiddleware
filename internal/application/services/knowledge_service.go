package query

import (
	"context"
	"fmt"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/knowledge"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/query"
)

type KnowledgeService struct {
	docRepo       document.DocumentRepository
	knowledgeRepo knowledge.KnowledgeRepository
	queryService  query.QueryService
}

func (s *KnowledgeService) UploadDocument(ctx context.Context, kbID string, doc *document.Document) error {
	// 1. 存储文档
	if err := s.docRepo.Store(ctx, doc); err != nil {
		return err
	}

	// 2. 添加到知识库
	return s.knowledgeRepo.AddDocument(kbID, doc)
}

func (s *KnowledgeService) QueryKnowledge(ctx context.Context, kbID string, q *query.Query) (*query.QueryResult, error) {
	// 1. 验证知识库存在
	kb, err := s.knowledgeRepo.FindByID(kbID)
	if err != nil {
		return nil, fmt.Errorf("知识库查询失败: %w", err)
	}
	if kb == nil {
		return nil, fmt.Errorf("知识库不存在")
	}

	// 2. 如果查询没有嵌入向量，则生成嵌入
	if len(q.Embedding) == 0 && q.Text != "" {
		embedding, err := s.queryService.GenerateEmbedding(ctx, q.Text)
		if err != nil {
			return nil, fmt.Errorf("生成嵌入向量失败: %w", err)
		}
		q.Embedding = embedding
	}

	// 3. 添加知识库过滤条件到查询
	q.Filter = fmt.Sprintf("knowledge_base_id == \"%s\"", kbID)

	// 4. 执行查询
	result, err := s.queryService.Execute(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("查询执行失败: %w", err)
	}

	return result, nil
}
