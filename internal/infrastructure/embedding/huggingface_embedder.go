package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

// HuggingFaceRequest 定义发送给本地 Hugging Face API 的请求体
type HuggingFaceRequest struct {
	Texts []string `json:"texts"`
}
type HuggingFaceResponse struct {
	Vectors [][]float32 `json:"vectors"`
}

func (e *HuggingFaceEmbedder) Embed(ctx context.Context, text string) (*embedding.Embedding, error) {
	if text == "" {
		return nil, fmt.Errorf("cannot embed empty text")
	}

	// 构建请求体
	reqBody := HuggingFaceRequest{
		Texts: []string{text},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", e.apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求（设置超时时间）
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inference server returned status code %d: %s", resp.StatusCode, respBody)
	}

	// 使用新的结构体解析响应
	var respData HuggingFaceResponse
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(respData.Vectors) == 0 || len(respData.Vectors[0]) == 0 {
		return nil, fmt.Errorf("received empty embedding vector")
	}

	// 返回结果
	return &embedding.Embedding{
		Vector: respData.Vectors[0],
		Model:  e.model,
	}, nil
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
