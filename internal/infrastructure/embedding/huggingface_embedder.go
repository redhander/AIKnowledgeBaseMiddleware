package embedding

import (
	"context"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/embedding"
)

type HuggingFaceEmbedder struct {
	apiURL string
	apiKey string
	model  string
}

func NewHuggingFaceEmbedder(apiURL, apiKey, model string) *HuggingFaceEmbedder {
	return &HuggingFaceEmbedder{
		apiURL: apiURL,
		apiKey: apiKey,
		model:  model,
	}
}

func (e *HuggingFaceEmbedder) Embed(ctx context.Context, text string) (*embedding.Embedding, error) {
	// 调用HuggingFace API实现嵌入
	// ...
	return nil, nil
}

func (e *HuggingFaceEmbedder) EmbedBatch(ctx context.Context, texts []string) ([]*embedding.Embedding, error) {
	var embeddings []*embedding.Embedding
	for _, text := range texts {
		embedding, err := e.Embed(ctx, text)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, embedding)
	}
	return embeddings, nil
}
