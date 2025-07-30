package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
)

type PDFParser struct {
	textSplitter *TextSplitter
}

func NewPDFParser(chunkSize, chunkOverlap int) *PDFParser {
	return &PDFParser{
		textSplitter: NewTextSplitter(chunkSize, chunkOverlap),
	}
}

func (p *PDFParser) Parse(filePath string) ([]*Document, error) {
	content, err := extractTextFromPDF(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// 分割文本
	chunks, err := p.textSplitter.Split(content)
	if err != nil {
		return nil, fmt.Errorf("failed to split PDF text: %w", err)
	}

	// 构建文档对象
	var documents []*Document
	for _, chunk := range chunks {
		documents = append(documents, &Document{
			ID:      uuid.New().String(),
			Content: chunk,
			Metadata: Metadata{
				Filename:    filepath.Base(filePath),
				ContentType: "application/pdf",
				Size:        int64(len(chunk)),
			},
		})
	}

	return documents, nil
}

func (p *PDFParser) SupportedExtensions() []string {
	return []string{".pdf"}
}

// 使用 pdfcpu 提取文本
func extractTextFromPDF(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	reader, err := pdf.NewReader(file, fileInfo.Size())
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var builder strings.Builder
	totalPage := reader.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", pageIndex, err)
		}
		builder.WriteString(text)
	}

	return builder.String(), nil
}
