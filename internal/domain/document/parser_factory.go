package document

import (
	"fmt"
	"strings"
)

type ParserFactory struct {
	parsers map[string]DocumentParser // 扩展名 -> 解析器
}

func NewParserFactory(chunkSize, chunkOverlap int) *ParserFactory {
	return &ParserFactory{
		parsers: map[string]DocumentParser{
			".pdf":  NewPDFParser(chunkSize, chunkOverlap),
			".txt":  NewTextParser(chunkSize, chunkOverlap),
			".docx": NewDOCXParser(chunkSize, chunkOverlap),
			".xlsx": NewXLSParser(chunkSize, chunkOverlap),
			".xls":  NewXLSParser(chunkSize, chunkOverlap),
		},
	}
}

func (f *ParserFactory) GetParser(fileExt string) (DocumentParser, error) {
	parser, exists := f.parsers[strings.ToLower(fileExt)]
	if !exists {
		return nil, fmt.Errorf("no parser found for extension: %s", fileExt)
	}
	return parser, nil
}

func (f *ParserFactory) SupportedExtensions() []string {
	exts := make([]string, 0, len(f.parsers))
	for ext := range f.parsers {
		exts = append(exts, ext)
	}
	return exts
}
