package document

import (
	"fmt"
	"path/filepath"
	"strings"
)

type UniversalParser struct {
	docxParser *DOCXParser
	pdfParser  *PDFParser
	textParser *TextParser
}

func NewUniversalParser(chunkSize, chunkOverlap int) *UniversalParser {
	return &UniversalParser{
		docxParser: NewDOCXParser(chunkSize, chunkOverlap),
		pdfParser:  NewPDFParser(chunkSize, chunkOverlap),
		textParser: NewTextParser(chunkSize, chunkOverlap),
	}
}

func (p *UniversalParser) Parse(filePath string) ([]*Document, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return p.pdfParser.Parse(filePath)
	case ".txt":
		return p.textParser.Parse(filePath)
	case ".docx":
		return p.docxParser.Parse(filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func (p *UniversalParser) SupportedExtensions() []string {
	return []string{".pdf", ".txt", ".docx"}
}
