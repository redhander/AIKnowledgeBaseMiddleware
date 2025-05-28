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
  
  - cd /deployments
  
  - docker-compose build huggingface
  
  - docker-compose up -d

3.Install Dependencies
  
  - cd cmd/server
  
  - go mod tidy

4.Configure Environment

  - cd configs
  
  - edit confg.yaml with your own params
  
  - edit with your API keys and Milvus config

5.Run System
  
  - go run main.go

6. Web access
   
https://github.com/redhander/AIKnowlegeBaseWeb
node version: v22.16.0
- clone the web repo
- npm install
- npm run serve
- http://localhost:5173

------------------------------------------------------------------------------------------------------------------------------------------------------
# AIKnowledgeBaseMiddleware
# RAG System with Go, DeepSeek, Hugging Face & Milvus
用于智能问答的检索增强生成（RAG）系统，集成了 DeepSeek LLM、Hugging Face 嵌入和 Milvus 向量数据库。

## 🌟 核心功能

- 多模态检索**  
  支持文本/PDF/docx/xlsx 解析 → 矢量化 → 语义搜索
- 混合生成***  
  结合DeepSeek的生成和检索上下文（Top-K + MMR重排）
- 实时知识更新***  
  动态更新Milvus知识库，增量索引
- 多语言支持**  
  通过 “All-MiniLM-L12-v2 ”处理多语言内容
- 性能优化***  
  GPU加速推理+批处理管道

##🛠️ 快速入门

### 先决条件

- Go 1.24.3+
- Milvus 2.5.12+ (Docker 部署。基于 Minio,ETCD)
- Deepseek-（使用 ollama 进行本地部署）
- Huggingface（本地部署或在线）
- 英伟达™（NVIDIA®）GPU（推荐）或纯 CPU 模式

#### 安装
1. 克隆仓库
   ```bash
    https://github.com/redhander/AIKnowledgeBaseMiddleware.git
    cd cmd/server

2.安装 Milvus/deepseek (Docker)

2.1 使用 ollama 在本地部署 deepkseek

2.2 部署其他
  
  - cd /deployments
  
  - docker-compose build huggingface
  
  - docker-compose up -d

3.安装依赖
  
  - cd cmd/server
  
  - go mod tidy

4.配置环境

  - cd configs
  
  - 使用自己的参数编辑 confg.yaml
  
  - 使用 API 密钥和 Milvus 配置进行编辑

5.运行系统
  
  - go run main.go

6. 前端页面访问
   
https://github.com/redhander/AIKnowlegeBase
 6.1 克隆仓库
 6.2 npm install 
 6.3 npm run dev
 6.4 访问 http://localhost:5173



   
