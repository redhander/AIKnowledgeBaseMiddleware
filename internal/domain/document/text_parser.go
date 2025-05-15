package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type TextParser struct {
	textSplitter *TextSplitter
}

func NewTextParser(chunkSize, chunkOverlap int) *TextParser {
	return &TextParser{
		textSplitter: NewTextSplitter(chunkSize, chunkOverlap),
	}
}

func (p *TextParser) Parse(filePath string) ([]*Document, error) {
	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %w", err)
	}

	// 自动检测编码并转换为UTF-8
	encoding, name, _ := charset.DetermineEncoding(rawContent, "")
	var utf8Content string
	//todo,下面代码可能需要将name改为encoding
	fmt.Println("encoding:", encoding, "name:", name)
	switch {
	case strings.Contains(name, "GBK"):
		decoder := simplifiedchinese.GBK.NewDecoder()
		utf8Content, _, err = transform.String(decoder, string(rawContent))
		if err != nil {
			return nil, fmt.Errorf("failed to convert GBK to UTF-8: %w", err)
		}
	default:
		utf8Content = string(rawContent)
	}

	chunks, err := p.textSplitter.Split(utf8Content)
	if err != nil {
		return nil, fmt.Errorf("failed to split text content: %w", err)
	}

	var documents []*Document
	for _, chunk := range chunks {
		documents = append(documents, &Document{
			Content: chunk,
			Metadata: Metadata{
				Filename:    filepath.Base(filePath),
				ContentType: "text/plain",
				Size:        int64(len(chunk)),
			},
		})
	}

	return documents, nil
}

func (p *TextParser) SupportedExtensions() []string {
	return []string{".txt"}
}
