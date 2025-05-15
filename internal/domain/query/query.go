package query

import (
	"context"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	deepseek "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/llm"
)

// Query 定义用户查询的参数
type Query struct {
	Text      string
	Embedding []float32
	TopK      int
}

// QueryResult 定义查询返回结果
type QueryResult struct {
	Answer   string
	Sources  []*document.Document
	Metadata map[string]interface{}
}

// QueryService 定义查询服务接口
type QueryService interface {
	Execute(ctx context.Context, query *Query) (*QueryResult, error)
}

// ===== 以下是 RAG 查询服务实现 =====

// RAGQueryService 是基于检索增强生成的查询服务实现
type RAGQueryService struct {
	LLM  *deepseek.Client
	Repo document.DocumentRepository
}

// NewRAGQueryService 创建一个新的 RAG 查询服务实例
func NewRAGQueryService(client *deepseek.Client, repo document.DocumentRepository) QueryService {
	return &RAGQueryService{
		LLM:  client,
		Repo: repo,
	}
}

// Execute 实现 QueryService 接口的执行方法
func (s *RAGQueryService) Execute(ctx context.Context, q *Query) (*QueryResult, error) {
	// 示例逻辑：搜索相关文档
	docs, err := s.Repo.Search(ctx, q.Embedding, q.TopK)
	if err != nil {
		return nil, err
	}

	// 拼接上下文并调用 LLM 生成回答
	var contextText string
	for _, doc := range docs {
		contextText += doc.Content + "\n"
	}

	prompt := "根据以下信息回答问题：\n" + contextText + "\n问题：" + q.Text

	answer, err := s.LLM.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &QueryResult{
		Answer:   answer,
		Sources:  docs,
		Metadata: map[string]interface{}{"model": s.LLM.Model()},
	}, nil
}
