package query

import (
	"context"

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
	if _, err := s.knowledgeRepo.FindByID(kbID); err != nil {
		return nil, err
	}

	// 2. 执行查询
	return s.queryService.Execute(ctx, q)
}
