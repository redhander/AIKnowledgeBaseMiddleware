package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	var contentBuilder strings.Builder
	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return nil, fmt.Errorf("failed to create extractor for page %d: %w", pageNum, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return nil, fmt.Errorf("failed to extract text from page %d: %w", pageNum, err)
		}

		contentBuilder.WriteString(text)
		contentBuilder.WriteString("\n\n") // 保留页面分隔
	}

	fullContent := contentBuilder.String()
	chunks, err := p.textSplitter.Split(fullContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split PDF text: %w", err)
	}

	var documents []*Document
	for _, chunk := range chunks {
		documents = append(documents, &Document{
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
