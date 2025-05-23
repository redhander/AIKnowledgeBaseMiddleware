package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/commands"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/application/queries"
	"github.com/redhander/AIKnowledgeBaseMiddleware/internal/infrastructure/logger"
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
	// 解析请求体
	r.ParseMultipartForm(10 << 20)

	// 获取文件句柄
	file, handler, err := r.FormFile("file")
	if err != nil {
		logger.Errorf("Error retrieving the file: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	logger.Infof("Uploaded File: %+v\n", handler.Filename)
	logger.Infof("File Size: %d\n", handler.Size)
	logger.Infof("MIME Header: %v\n", handler.Header)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		logger.Errorf("Failed to read file: %v", err)
		http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		return
	}
	logger.Infof("File content read successfully, size: %d bytes", len(fileContent))
	// 转换为 UploadDocumentCommand
	cmd := commands.UploadDocumentCommand{
		FileContent: fileContent,
		Filename:    handler.Filename,
	}
	// 调用 uploadHandler.Handle()
	ctx := r.Context()
	if err := h.uploadHandler.Handle(ctx, cmd); err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload document: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Document uploaded successfully"))
}

func (h *KnowledgeHandler) QueryKnowledge(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body)) //
	var req queries.QueryKnowledgeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Errorf("JSON decode error: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// 调用 queryHandler.Handle()
	ctx := r.Context()
	result, err := h.queryHandler.Handle(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query knowledge: %v", err), http.StatusInternalServerError)
		return
	}

	// 返回查询结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
