package query

import (
	"context"
	"fmt"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	deepseek "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/llm"
)

// Query 定义用户查询的参数
type Query struct {
	Text      string
	Embedding []float32
	TopK      int
	Filter    string // 过滤条件
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
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
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
	// docs, err := s.Repo.Search(ctx, q.Embedding, q.TopK, q.Filter)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("搜索结果DB：", docs)
	// // 拼接上下文并调用 LLM 生成回答
	// var contextText string
	// for _, doc := range docs {
	// 	contextText += doc.Content + "\n"
	// }

	// prompt := "根据以下信息回答问题：\n" + contextText + "\n问题是：" + q.Text + "\n过滤条件为：" + q.Filter

	// answer, err := s.LLM.Generate(ctx, prompt)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("LLM 生成的回答：", answer)
	// return &QueryResult{
	// 	Answer:   answer,
	// 	Sources:  docs,
	// 	Metadata: map[string]interface{}{"model": s.LLM.Model()},
	// }, nil
	docs, err := s.Repo.Search(ctx, q.Embedding, q.TopK, q.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	fmt.Println("搜索结果DB：", docs)

	// Handle case when no documents are found
	if len(docs) == 0 {
		// Fallback to general LLM response when no relevant documents found
		prompt := "请回答以下问题：\n" + q.Text
		if q.Filter != "" {
			prompt += "\n注意：请考虑以下过滤条件：" + q.Filter
		}

		answer, err := s.LLM.Generate(ctx, prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to generate answer without context: %w", err)
		}

		return &QueryResult{
			Answer:  answer,
			Sources: []*document.Document{},
			Metadata: map[string]interface{}{
				"model":   s.LLM.Model(),
				"context": "no_relevant_documents",
			},
		}, nil
	}

	// Build context with length awareness
	var contextText string
	totalLength := 0
	maxContextLength := 3000 // Adjust based on your LLM's context window

	for _, doc := range docs {
		// Check if adding this document would exceed context limit
		docLength := len([]rune(doc.Content))
		if totalLength+docLength > maxContextLength {
			// Truncate document content or skip it
			remainingSpace := maxContextLength - totalLength
			if remainingSpace > 100 { // Only add if there's meaningful space
				truncatedContent := string([]rune(doc.Content)[:remainingSpace])
				contextText += truncatedContent + "\n[文档内容已截断]\n"
			}
			break
		}

		contextText += "文档内容: " + doc.Content + "\n\n"
		totalLength += docLength
	}

	// Improved prompt structure
	prompt := fmt.Sprintf(`基于以下文档内容回答问题。如果文档内容不包含足够信息，请说明无法基于提供的文档回答该问题。文档内容:%s问题: %s`, contextText, q.Text)

	if q.Filter != "" {
		prompt += fmt.Sprintf("\n过滤条件: %s", q.Filter)
	}

	prompt += "\n\n请基于上述文档内容回答问题，如果文档中没有相关信息，请回答'根据提供的文档无法回答该问题'。"
	fmt.Printf("Prompt: %s\n", prompt)
	answer, err := s.LLM.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate answer with context: %w", err)
	}

	return &QueryResult{
		Answer:  answer,
		Sources: docs,
		Metadata: map[string]interface{}{
			"model":           s.LLM.Model(),
			"documents_found": len(docs),
			"context_length":  len(contextText),
		},
	}, nil
}
func (s *RAGQueryService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// 使用LLM客户端生成嵌入向量
	embedding, err := s.LLM.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("生成嵌入向量失败: %w", err)
	}
	return embedding, nil
}
