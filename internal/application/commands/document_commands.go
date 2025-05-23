package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/embedding"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
)

type UploadDocumentRequest struct {
	FilePath string            `json:"file_path"` // Path to the document file
	Metadata map[string]string `json:"metadata"`  // Optional metadata
}

// Metadata 结构
type Metadata map[string]interface{}

// Document 增强定义（domain层）
type Document struct {
	ID       string
	Content  string
	Vector   []float32
	Metadata DocumentMetadata
}

type DocumentMetadata struct {
	Filename     string
	ContentType  string
	UploadTime   time.Time
	OriginalFile string
	CustomFields map[string]interface{}
}
type UploadDocumentCommand struct {
	FileContent []byte   // 文件二进制内容
	Filename    string   // 原始文件名
	UserID      string   // 上传用户标识（可选）
	Attributes  Metadata // 自定义属性（可选）
}
type UploadDocumentHandler struct {
	parserFactory *document.ParserFactory
	embedder      embedding.Embedder
	docRepo       document.DocumentRepository
}

//	type UploadDocumentHandler struct {
//		docParser     document.DocumentParser
//		docSplitter   document.DocumentSplitter
//		embedder      embedding.Embedder
//		docRepo       document.DocumentRepository
//		knowledgeRepo knowledge.KnowledgeRepository
//	}
func NewUploadDocumentHandler(
	factory *document.ParserFactory,
	embedder embedding.Embedder,
	repo document.DocumentRepository,
) *UploadDocumentHandler {
	if factory == nil {
		panic("parserFactory cannot be nil")
	}
	if embedder == nil {
		panic("embedder cannot be nil")
	}
	if repo == nil {
		panic("docRepo cannot be nil")
	}
	return &UploadDocumentHandler{
		parserFactory: factory,
		embedder:      embedder,
		docRepo:       repo,
	}
}

// type UploadDocumentCommand struct {
// 	FileContent []byte
// 	Filename    string
// }

func (h *UploadDocumentHandler) Handle(ctx context.Context, cmd UploadDocumentCommand) error {
	startTime := time.Now()
	log := logger.FromContext(ctx).WithFields(map[string]interface{}{
		"filename": cmd.Filename,
		"size":     len(cmd.FileContent),
	})

	// 1. 验证文件扩展名
	ext := strings.ToLower(filepath.Ext(cmd.Filename))
	if ext == "" {
		log.Warn("File has no extension")
		return fmt.Errorf("file must have an extension")
	}

	// 2. 获取对应解析器
	parser, err := h.parserFactory.GetParser(ext)
	if err != nil {
		log.Warnf("Unsupported file type: %s", ext)
		return errors.New("unsupported file type")
	}

	// 3. 创建临时文件（带随机后缀防止冲突）
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, fmt.Sprintf("upload_*%s", ext))
	if err != nil {
		log.Errorf("Failed to create temp file: %v", err)
		return fmt.Errorf("failed to create temporary storage")
	}
	tmpPath := tmpFile.Name()

	// 4. 确保清理临时文件
	defer func() {
		if err := os.Remove(tmpPath); err != nil {
			log.Warnf("Failed to remove temp file %s: %v", tmpPath, err)
		}
	}()

	// 5. 写入临时文件
	if _, err := tmpFile.Write(cmd.FileContent); err != nil {
		log.Errorf("Failed to write temp file: %v", err)
		return fmt.Errorf("failed to prepare document for processing")
	}
	if err := tmpFile.Close(); err != nil {
		log.Warnf("Failed to close temp file: %v", err)
	}

	// 6. 解析文档内容
	log.Info("Start parsing document")
	docs, err := parser.Parse(tmpPath)
	if err != nil {
		log.Errorf("Failed to parse document: %v", err)
		return fmt.Errorf("document parsing failed: %w", err)
	}
	log.Infof("Parsed into %d chunks", len(docs))
	// 7. 生成向量嵌入
	vectorsGenerated := 0
	for _, doc := range docs {
		select {
		case <-ctx.Done():
			log.Info("Processing cancelled by context")
			return ctx.Err()
		default:
			// 设置文档元数据
			doc.Metadata.UploadTime = time.Now()
			doc.Metadata.OriginalFile = cmd.Filename

			// 生成嵌入向量
			embedding, err := h.embedder.Embed(ctx, doc.Content)
			if err != nil {
				log.Warnf("Failed to generate embedding for chunk: %v", err)
				continue // 跳过失败的分块或 return err 根据业务需求
			}
			if embedding == nil {
				log.Warn("Received nil embedding")
				continue
			}
			doc.Vector = embedding.Vector
			vectorsGenerated++
		}
	}

	if vectorsGenerated == 0 {
		log.Error("No vectors generated for document")
		return fmt.Errorf("failed to generate any embeddings")
	}

	// 8. 存储到向量数据库
	if err := h.docRepo.StoreBatch(ctx, docs); err != nil {
		log.Errorf("Failed to store documents: %v", err)
		return fmt.Errorf("failed to save document knowledge: %w", err)
	}

	// 9. 记录处理指标
	duration := time.Since(startTime)
	log.Infof("Successfully processed document in %v (chunks: %d, vectors: %d)",
		duration, len(docs), vectorsGenerated)

	return nil
}
