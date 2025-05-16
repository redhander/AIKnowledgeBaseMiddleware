package main

import (
	"context"
	"io"
	"log"
	httpO "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/commands"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/queries" // 添加这一行
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	config "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/config"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/embedding"
	deepseek "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/llm" // 添加deepseek包导入
	logger "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
	milvus "github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/persistence/milvus" // 添加milvus包导入
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http/handler"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http/middleware"
)

func main() {
	// 1. 初始化配置
	cfg, err := config.Load("../../configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志系统
	logger := initLogger(cfg.Logging)
	// defer logger.Sync() // 确保日志缓冲区被刷新

	rootCtx := context.WithValue(context.Background(), "logger", logger)
	logger.Info("Application starting...")
	logger.Infof("Loaded configuration: %+v", cfg.Sanitized()) // 确保敏感配置被过滤

	// 3. 初始化基础设施组件
	milvusClient, err := initMilvus(rootCtx, cfg.Milvus)
	if err != nil {
		logger.Errorf("Failed to initialize Milvus: %v", err)
		return
	}
	if milvusClient == nil {
		logger.Errorf("Failed to initialize Milvus: %v", err)
		return
	}
	embedder := embedding.NewHuggingFaceEmbedder(cfg.Embedding.ModelName, cfg.Embedding.APIKey, cfg.Embedding.Model)
	if embedder == nil {
		logger.Fatalf("Failed to initialize embedder")
	}
	logger.Infof("Embedder initialized successfully with model: %s", embedder)
	deepseekClient := initDeepSeek(cfg.DeepSeek)
	if deepseekClient == nil {
		logger.Fatalf("Failed to initialize DeepSeek client")
	}
	logger.Infof("DeepSeek client initialized successfully with model: %s", deepseekClient.Model())
	// 4. 初始化解析器工厂
	parserFactory := document.NewParserFactory(
		cfg.Document.ChunkSize,
		cfg.Document.ChunkOverlap,
	)

	// 5. 初始化存储库
	docRepo := milvus.NewMilvusDocumentRepository(
		milvusClient,
		cfg.Milvus.CollectionName,
	)

	// 6. 初始化应用层
	uploadHandler := commands.NewUploadDocumentHandler(
		parserFactory,
		embedder,
		docRepo,
	)

	queryHandler := queries.NewQueryKnowledgeHandler(
		embedder,
		docRepo,
		deepseekClient,
	)

	// 7. 初始化HTTP服务
	httpHandler := handler.NewKnowledgeHandler(uploadHandler, queryHandler)
	router := http.NewRouter(
		httpHandler,
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)

	srv := &httpO.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. 启动服务
	go func() {
		logger.Infof("Server listening on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != httpO.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// 9. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(rootCtx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("Server exited properly")
}

// 初始化日志系统
func initLogger(cfg config.LoggingConfig) logger.Logger {
	var output io.Writer = os.Stdout

	if cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		if cfg.Console {
			output = io.MultiWriter(os.Stdout, file)
		} else {
			output = file
		}
	}

	return logger.New(output).WithFields(logger.Fields{
		"app": "enterprise-knowledge",
		"env": cfg.Environment,
	})
}

// 初始化Milvus客户端
func initMilvus(ctx context.Context, cfg config.MilvusConfig) (*milvus.MilvusClient, error) {
	options := []milvus.Option{
		milvus.WithAddress(cfg.Address),
	}

	if cfg.Username != "" {
		options = append(options, milvus.WithAuth(cfg.Username, cfg.Password))
	}

	return milvus.NewMilvusClient(ctx, cfg)
}

// 初始化DeepSeek客户端
func initDeepSeek(cfg config.DeepSeekConfig) *deepseek.Client {
	return deepseek.NewClient(
		cfg.BaseURL,
		deepseek.WithAPIKey(cfg.APIKey),
		deepseek.WithModel(cfg.Model),
	)
}
