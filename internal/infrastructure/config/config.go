package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Config 应用全局配置
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Milvus    MilvusConfig    `yaml:"milvus"`
	Embedding EmbeddingConfig `yaml:"embedding"`
	DeepSeek  DeepSeekConfig  `yaml:"deepseek"`
	Document  DocumentConfig  `yaml:"document"`
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Environment string `yaml:"environment"` // production/development
	FilePath    string `yaml:"file_path"`   // 日志文件路径
	Console     bool   `yaml:"console"`     // 是否同时输出到控制台
}

// MilvusConfig Milvus向量数据库配置
type MilvusConfig struct {
	Address        string    `yaml:"address"`
	Username       string    `yaml:"username"`
	Password       string    `yaml:"password"`
	CollectionName string    `yaml:"collection_name"`
	TLS            TLSConfig `yaml:"tls"`
}

// TLSConfig TLS安全配置
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertPath string `yaml:"cert_path"`
}

// EmbeddingConfig 文本嵌入配置
type EmbeddingConfig struct {
	ModelName string        `yaml:"model_name"`
	APIKey    string        `yaml:"api_key"`
	Timeout   time.Duration `yaml:"timeout"`
	ApiURL    string        `yaml:"api_url"`
	Model     string        `yaml:"model"`
}

// DeepSeekConfig DeepSeek LLM配置
type DeepSeekConfig struct {
	BaseURL string        `yaml:"base_url"`
	APIKey  string        `yaml:"api_key"`
	Model   string        `yaml:"model"`
	Timeout time.Duration `yaml:"timeout"`
}

// DocumentConfig 文档处理配置
type DocumentConfig struct {
	ChunkSize    int    `yaml:"chunk_size"`    // 文本分块大小
	ChunkOverlap int    `yaml:"chunk_overlap"` // 分块重叠大小
	MaxFileSize  string `yaml:"max_file_size"` // 最大文件大小(如10MB)
}

// Load 从YAML文件加载配置
func Load(configPath string) (*Config, error) {
	// 解析文件路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute config path: %w", err)
	}

	// 读取文件内容
	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML
	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Sanitized 返回脱敏后的配置(用于日志记录)
func (c *Config) Sanitized() Config {
	sanitized := *c
	// 脱敏敏感字段
	sanitized.Embedding.APIKey = "***"
	sanitized.DeepSeek.APIKey = "***"
	sanitized.Milvus.Password = "***"
	return sanitized
}

// validate 配置验证
func (c *Config) validate() error {
	// 服务器配置验证
	if c.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	// Milvus配置验证
	if c.Milvus.Address == "" {
		return fmt.Errorf("milvus address is required")
	}
	if c.Milvus.CollectionName == "" {
		return fmt.Errorf("milvus collection name is required")
	}

	// 嵌入模型验证
	if c.Embedding.ModelName == "" {
		return fmt.Errorf("embedding model name is required")
	}

	// DeepSeek验证
	if c.DeepSeek.BaseURL == "" {
		return fmt.Errorf("deepseek base URL is required")
	}
	if c.DeepSeek.Model == "" {
		return fmt.Errorf("deepseek model is required")
	}

	// 文档处理验证
	if c.Document.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}
	if c.Document.ChunkOverlap < 0 {
		return fmt.Errorf("chunk overlap cannot be negative")
	}
	if c.Document.ChunkOverlap >= c.Document.ChunkSize {
		return fmt.Errorf("chunk overlap must be smaller than chunk size")
	}

	return nil
}

// GetMaxFileSizeBytes 解析最大文件大小字符串为字节数
func (d *DocumentConfig) GetMaxFileSizeBytes() (int64, error) {
	// 实现字符串如"10MB"到字节数的转换
	// 示例实现(需根据实际需求完善):
	return 10 * 1024 * 1024, nil // 默认10MB
}
