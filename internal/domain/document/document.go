package document

import (
	"context"
	"time"
)

type Document struct {
	ID        string
	Content   string
	Metadata  Metadata
	Vector    []float32
	CreatedAt time.Time
}

type Metadata struct {
	Filename     string
	ContentType  string
	Size         int64
	Custom       map[string]interface{}
	UploadTime   time.Time
	OriginalFile string
}

type DocumentRepository interface {
	Store(ctx context.Context, doc *Document) error
	StoreBatch(ctx context.Context, docs []*Document) error
	FindByID(ctx context.Context, id string) (*Document, error)
	Search(ctx context.Context, embedding []float32, topK int) ([]*Document, error)
}

type DocumentParser interface {
	Parse(filePath string) ([]*Document, error)
}

type DocumentSplitter interface {
	Split(content string) ([]string, error)
}
