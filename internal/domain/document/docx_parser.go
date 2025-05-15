package document

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/unidoc/unioffice/document"
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
	// Open the DOCX file
	doc, err := document.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX file: %w", err)
	}

	var contentBuilder strings.Builder

	// Process all paragraphs in the document
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			contentBuilder.WriteString(run.Text())
		}
		contentBuilder.WriteString("\n\n") // Maintain paragraph separation
	}

	// Process all tables in the document
	for _, tbl := range doc.Tables() {
		for _, row := range tbl.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						contentBuilder.WriteString(run.Text())
					}
				}
				contentBuilder.WriteString("\t") // Cell separator
			}
			contentBuilder.WriteString("\n") // Row separator
		}
		contentBuilder.WriteString("\n\n") // Table separator
	}

	fullContent := contentBuilder.String()

	// Split the content into chunks
	chunks, err := p.textSplitter.Split(fullContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split text: %w", err)
	}

	// Create domain documents from chunks
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
