package document

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type XLSParser struct {
	textSplitter *TextSplitter
}

func NewXLSParser(chunkSize, chunkOverlap int) *XLSParser {
	return &XLSParser{
		textSplitter: NewTextSplitter(chunkSize, chunkOverlap),
	}
}

func (p *XLSParser) Parse(filePath string) ([]*Document, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	var contentBuilder strings.Builder

	// 处理所有工作表
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows from sheet %s: %w", sheet, err)
		}

		for i, row := range rows {
			var rowContent []string
			for _, cell := range row {
				rowContent = append(rowContent, strings.TrimSpace(cell))
			}
			contentBuilder.WriteString(strings.Join(rowContent, "\t"))
			contentBuilder.WriteString("\n")

			// 每处理100行添加一个分隔符
			if (i+1)%100 == 0 {
				contentBuilder.WriteString("\n\n")
			}
		}
		contentBuilder.WriteString("\n\n") // 工作表分隔
	}

	fullContent := contentBuilder.String()
	chunks, err := p.textSplitter.Split(fullContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split Excel content: %w", err)
	}

	var documents []*Document
	for _, chunk := range chunks {
		documents = append(documents, &Document{
			Content: chunk,
			Metadata: Metadata{
				Filename:    filepath.Base(filePath),
				ContentType: p.getContentType(filePath),
				Size:        int64(len(chunk)),
			},
		})
	}

	return documents, nil
}

func (p *XLSParser) getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	default:
		return "application/octet-stream"
	}
}

func (p *XLSParser) SupportedExtensions() []string {
	return []string{".xlsx", ".xls"}
}
