package document

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/qifengzhang007/gooxml/document"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
)

// DOCXParser implements DocumentParser for Microsoft Word (.docx) files
type DOCXParser struct {
	textSplitter *TextSplitter
}

// NewDOCXParser creates a new DOCX parser instance
func NewDOCXParser(chunkSize, chunkOverlap int) *DOCXParser {
	return &DOCXParser{
		textSplitter: NewTextSplitter(chunkSize, chunkOverlap),
	}
}

// Parse extracts text content from a DOCX file and splits it into chunks
func (p *DOCXParser) Parse(filePath string) ([]*Document, error) {
	logger.Infof("Parsing DOCX file: %s", filePath)
	doc, err := document.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("opened DOCX document is nil")
	}
	var contentBuilder strings.Builder
	// 提取所有段落文本
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			contentBuilder.WriteString(run.Text())
		}
		contentBuilder.WriteString("\n\n") // 段落之间保留空行
	}

	// 提取表格内容（可选）
	for _, tbl := range doc.Tables() {
		for _, row := range tbl.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						contentBuilder.WriteString(run.Text())
					}
				}
				contentBuilder.WriteString("\t") // 单元格分隔符
			}
			contentBuilder.WriteString("\n") // 行结束换行
		}
		contentBuilder.WriteString("\n\n") // 表格之间留空行
	}

	fullContent := contentBuilder.String()

	// 使用 textSplitter 进行分块处理
	chunks, err := p.textSplitter.Split(fullContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split text: %w", err)
	}

	// 构建返回结果
	var documents []*Document
	for _, chunk := range chunks {
		documents = append(documents, &Document{
			Content: chunk,
			Metadata: Metadata{
				Filename:    filepath.Base(filePath),
				ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				Size:        int64(len(chunk)),
			},
		})
	}

	return documents, nil
}

// SupportedExtensions returns the file extensions this parser supports
func (p *DOCXParser) SupportedExtensions() []string {
	return []string{".docx"}
}
