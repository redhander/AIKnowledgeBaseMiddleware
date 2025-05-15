package knowledge

import (
	"time"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/domain/document"
)

type KnowledgeBase struct {
	ID          string
	Name        string
	Description string
	Documents   []*document.Document
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type KnowledgeRepository interface {
	AddDocument(kbID string, doc *document.Document) error
	RemoveDocument(kbID string, docID string) error
	FindByID(kbID string) (*KnowledgeBase, error)
}
