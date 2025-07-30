package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
)

type DeepSeekQueryService struct {
	baseURL string
	apiKey  string
	model   string
}

func NewDeepSeekQueryService(baseURL, apiKey, model string) *DeepSeekQueryService {
	return &DeepSeekQueryService{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
	}
}

func (s *DeepSeekQueryService) Generate(ctx context.Context, prompt string) (string, error) {
	logger.Info("deepseek_service.Generates.model: ", s.model)
	requestBody := map[string]interface{}{
		"model":  s.model,
		"prompt": prompt,
		"stream": false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/generate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Response, nil
}

func (s *DeepSeekQueryService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	logger.Info("Generating embedding for text with model: ", s.model)

	requestBody := map[string]interface{}{
		"model": s.model,
		"input": text,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	return result.Data[0].Embedding, nil
}
