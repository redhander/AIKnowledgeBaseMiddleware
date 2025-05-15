package queries

import (
	"context"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/embedding"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/query"
	deepseek "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/llm" // 添加导入
)

type QueryKnowledgeRequest struct {
	Text string
	TopK int
}

type QueryKnowledgeResponse struct {
	Answer  string
	Sources []*document.Document
}

type QueryKnowledgeHandler struct {
	embedder     embedding.Embedder
	docRepo      document.DocumentRepository
	queryService query.QueryService
}

func (h *QueryKnowledgeHandler) Handle(ctx context.Context, req QueryKnowledgeRequest) (*QueryKnowledgeResponse, error) {
	// 1. 嵌入查询
	embedding, err := h.embedder.Embed(ctx, req.Text)
	if err != nil {
		return nil, err
	}

	// 2. 执行查询
	result, err := h.queryService.Execute(ctx, &query.Query{
		Text:      req.Text,
		Embedding: embedding.Vector,
		TopK:      req.TopK,
	})
	if err != nil {
		return nil, err
	}

	return &QueryKnowledgeResponse{
		Answer:  result.Answer,
		Sources: result.Sources,
	}, nil
}

func NewQueryKnowledgeHandler(embedder embedding.Embedder, repo document.DocumentRepository, client *deepseek.Client) *QueryKnowledgeHandler {
	return &QueryKnowledgeHandler{
		embedder:     embedder,
		docRepo:      repo,
		queryService: query.NewRAGQueryService(client, repo), // 示例初始化逻辑
	}
}
