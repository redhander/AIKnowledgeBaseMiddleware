package deepseek

import (
	"context"
)

// Client 是 DeepSeek 的客户端实现
type Client struct {
	service *DeepSeekQueryService
}
type ClientOption func(*Client)

// NewClient 创建一个新的 DeepSeek 客户端
func NewClient(baseURL string, opts ...ClientOption) *Client {
	client := &Client{
		service: &DeepSeekQueryService{
			baseURL: baseURL,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.service.apiKey = apiKey
	}
}

func WithModel(model string) ClientOption {
	return func(c *Client) {
		c.service.model = model
	}
}

// Generate 调用 DeepSeek 生成回答
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	// 这里调用 DeepSeekQueryService 的实现
	return c.service.Generate(ctx, prompt)
}

// Model 返回当前使用的模型名称
func (c *Client) Model() string {
	return c.service.model
}
