# AIKnowledgeBaseMiddleware
# RAG System with Go, DeepSeek, Hugging Face & Milvus

[![Go Version](https://img.shields.io/badge/Go-1.24.3+-00ADD8?logo=go)](https://golang.org/)
[![Milvus](https://img.shields.io/badge/Milvus-2.5.12+-00B4D8)](https://milvus.io/)
MIT License +　Additional Clause.
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
- Milvus 2.5.12+ (Docker deployment. Based on Minio,ETCD)
- Deepseek-(local deployment with ollama)
- Huggingface(local deployment or online)
- NVIDIA GPU (recommended) or CPU-only mode

### Installation
1. Clone Repository
   ```bash
    https://github.com/redhander/AIKnowledgeBaseMiddleware.git
    cd cmd/server

2.Setup Milvus/deepseek (Docker)

2.1 deploy deepkseek on local with ollama

2.2 deploy others
  
  -cd /deployments
  
  -docker-compose build huggingface
  
  -docker-compose up -d

3.Install Dependencies
  
  -cd cmd/server
  
  -go mod tidy

4.Configure Environment

  -cd configs
  
  -edit confg.yaml with your own params
  
  -edit with your API keys and Milvus config

5.Run System
  
  -go run main.go

6. Web access
   
https://github.com/redhander/AIKnowlegeBaseWeb
node version: v22.16.0
- clone the web repo
- npm install
- npm run serve
- http://localhost:5173





   
