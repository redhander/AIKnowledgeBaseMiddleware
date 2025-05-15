package http

import (
	"github.com/gorilla/mux"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http/handler"
)

func NewRouter(kh *handler.KnowledgeHandler, middlewares ...mux.MiddlewareFunc) *mux.Router {
	r := mux.NewRouter()
	// Apply middlewares
	for _, m := range middlewares {
		r.Use(m)
	}
	r.HandleFunc("/api/documents", kh.UploadDocument).Methods("POST")
	r.HandleFunc("/api/query", kh.QueryKnowledge).Methods("POST")

	return r
}
