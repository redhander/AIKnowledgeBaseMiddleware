package http

import (
	"github.com/gorilla/mux"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http/handler"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/interfaces/http/middleware"
)

func NewRouter(kh *handler.KnowledgeHandler, logger logger.Logger) *mux.Router {
	r := mux.NewRouter()

	// ✅ 先注册中间件
	r.Use(middleware.CORS(
		[]string{"http://localhost:5173"},
		[]string{"POST", "OPTIONS"},
		[]string{"Content-Type", "Authorization"},
	))
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recovery(logger))

	// ✅ 再注册路由
	r.HandleFunc("/api/documents", kh.UploadDocument).Methods("POST")
	r.HandleFunc("/api/query", kh.QueryKnowledge).Methods("POST")

	return r
}
