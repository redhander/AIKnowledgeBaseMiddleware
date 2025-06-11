package document

import (
	"context"
	"encoding/json"
	"time"
)

type Document struct {
	ID        string
	Content   string
	Metadata  Metadata
	Vector    []float32
	Score     float32
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

func (m *Metadata) String() string {
	metadata := map[string]interface{}{
		"filename":      m.Filename,
		"content_type":  m.ContentType,
		"upload_time":   m.UploadTime.Format(time.RFC3339),
		"original_file": m.OriginalFile,
		"custom_fields": m.Custom,
	}

	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return "{}" // Return empty JSON object if serialization fails
	}
	return string(jsonData)
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
