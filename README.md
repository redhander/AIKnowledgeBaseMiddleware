# AIKnowledgeBaseMiddleware
# RAG System with Go, DeepSeek, Hugging Face & Milvus

[![Go Version](https://img.shields.io/badge/Go-1.24.3+-00ADD8?logo=go)](https://golang.org/)
[![Milvus](https://img.shields.io/badge/Milvus-2.5.12+-00B4D8)](https://milvus.io/)
MIT License +　Additional Clause
A Retrieval-Augmented Generation (RAG) system for intelligent Q&A, integrating DeepSeek LLM, Hugging Face embeddings, and Milvus vector database.

## 🌟 Core Features

- **Multi-modal Retrieval**  
  Supports text/PDF/docx/xlsx parsing → vectorization → semantic search
- **Hybrid Generation**  
  Combines DeepSeek's generation with retrieved context (Top-K + MMR reranking)
- **Real-time Knowledge Updates**  
  Dynamic Milvus knowledge base updates with incremental indexing
- **Multilingual Support**  
  Handles multilingual content via `All-MiniLM-L12-v2`
- **Performance Optimized**  
  GPU-accelerated inference + batch processing pipeline

## 🛠️ Quick Start

### Prerequisites

- Go 1.24.3+
- Milvus 2.5.12+ (Docker deployment)
- NVIDIA GPU (recommended) or CPU-only mode

### Installation
1. **Clone Repository**
   ```bash
    https://github.com/redhander/AIKnowledgeBaseMiddleware.git
    cd rag-system

2.Setup Milvus/deepseek (Docker)
2.1 deploy deepkseek on local with ollama
2.2 bash
cd /deployments
docker-compose build huggingface
docker-compose up -d

3.Install Dependencies

bash
cd cmd/server
go mod tidy

4.Configure Environment
 cd configs
 edit confg.yaml with your own params
# Edit with your API keys and Milvus config

5.Run System

bash
go run main.go












   
