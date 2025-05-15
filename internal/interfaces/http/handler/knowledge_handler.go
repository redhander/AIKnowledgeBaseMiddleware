package handler

import (
	"net/http"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/commands"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/queries"
)

type KnowledgeHandler struct {
	uploadHandler *commands.UploadDocumentHandler
	queryHandler  *queries.QueryKnowledgeHandler
}

func NewKnowledgeHandler(upload *commands.UploadDocumentHandler, query *queries.QueryKnowledgeHandler) *KnowledgeHandler {
	return &KnowledgeHandler{
		uploadHandler: upload,
		queryHandler:  query,
	}
}

func (h *KnowledgeHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	// 解析请求并调用uploadHandler.Handle()
	// ...
}

func (h *KnowledgeHandler) QueryKnowledge(w http.ResponseWriter, r *http.Request) {
	// 解析请求并调用queryHandler.Handle()
	// ...
}
