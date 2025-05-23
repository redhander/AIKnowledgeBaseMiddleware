package embedding

import "context"

type Embedding struct {
	Vector []float32
	Model  string
}

type Embedder interface {
	Embed(context.Context, string) (*Embedding, error)
	EmbedBatch(context.Context, []string) ([]*Embedding, error)
}
