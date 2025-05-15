package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
)

func Logging(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("Panic recovered: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
