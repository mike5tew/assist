# ESP Organizer

## Project Overview
The ESP Organizer project is a semantic skills database system using Go, Weaviate (vector database), MongoDB, and a React frontend to organize educational skills and materials. **The entire system is architected around an external AI/LLM provider (like AWS Bedrock) for all semantic understanding and vectorization tasks.**

## Core Architecture: External Vectorization
A critical design choice is that Weaviate is configured with `DEFAULT_VECTORIZER_MODULE: 'none'`. This system **does not** use built-in sentence transformers. Instead, it follows a sophisticated pipeline:
1.  **Raw Content Storage**: MongoDB stores original content (medical texts, educational materials).
2.  **Semantic Meaning Extraction**: An external LLM (e.g., via AWS Bedrock) extracts relationships, concepts, and semantic links from the raw content.
3.  **Relationship Vectorization**: An external embedding model (e.g., AWS Titan) generates vectors for the **extracted meanings**, which are then stored in Weaviate.

This approach allows for a much deeper level of semantic search, focusing on the relationships between concepts rather than just raw text similarity.

## Quick Start Guide

### 🚀 **Complete System Startup (Recommended)**

```bash
# 1. Navigate to project directory
cd /Users/michaelstewart/Coding/assist/esp-organizer

# 2. Configure your environment
cp .env.example .env
# Edit .env with your AWS credentials and other LLM settings.
# AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_REGION are required.

# 3. Start databases (MongoDB + Weaviate)
make start-databases

# 4. Initialize with educational skills data (uses AWS for vectorization)
make setup

# 5. Initialize with medical test data (uses AWS for vectorization)
make setup-medical

# 6. Start the API server
make server

# 7. (Optional) Start React frontend
cd frontend && npm install && npm start
```

### 🔧 **Step-by-Step Setup**

#### **Prerequisites**
```bash
# Verify required software
go version          # Should be Go 1.23+
node --version      # Should be Node.js 18+
docker --version    # Docker for databases
aws --version       # AWS CLI, configured with credentials
```

#### **1. Environment Setup**
```bash
# Install Go dependencies
go mod tidy

# Configure environment variables
cp .env.example .env
# Edit .env with your AWS credentials for medical processing
# OPENAI_API_KEY is a fallback, but AWS is the primary for the full pipeline.
```

#### **2. Database Setup**
```bash
# Start MongoDB and Weaviate with Docker
make start-databases

# Verify databases are running
make check-databases
# Expected output:
# ✅ MongoDB is running on port 27017
# ✅ Weaviate is running on port 8081
```

#### **3. Data Initialization**
```bash
# Clear any existing data and upload educational skills
# This command now uses your configured LLM (e.g., AWS) to generate vectors.
make setup

# Set up the medical knowledge base (book, excerpts, semantic links)
make setup-medical
```

#### **4. Start API Server**
```bash
# Start the Go API server
make server

# Expected output:
# 🚀 ESP Organizer API Server starting on port 8080
# ✅ MongoDB connected successfully
# ✅ Weaviate connected successfully
# ...
```

#### **5. Test the System**
```bash
# In another terminal, test basic functionality
curl http://localhost:8080/

# Test semantic search (will use your configured LLM)
curl "http://localhost:8080/api/skills/semantic-query?q=working+memory"

# Run comprehensive tests
make test
```

### 🧪 **Development Workflows**

#### **Educational Skills Development**
```bash
# Full reset and setup of educational skills
make purge-databases && make start-databases && make setup
```

#### **Medical Content Processing (AWS Required)**
```bash
# 1. Ensure AWS credentials are in .env

# 2. Test AWS setup
make verify-aws-credentials

# 3. Run the full medical content setup
make setup-medical

# 4. Check the results
make check-medical-content
```

### 🤖 **LLM Integration (Required for Functionality)**

This system is designed to work with an external LLM. AWS Bedrock is the primary intended provider.

**Setup (Required):**
```bash
# 1. Add your AWS credentials to your environment or .env file
#    AWS_ACCESS_KEY_ID=...
#    AWS_SECRET_ACCESS_KEY=...
#    AWS_REGION=eu-west-2 # or your preferred region

# 2. Restart the server if it's running
make server

# 3. Test a query that relies on the LLM
curl "http://localhost:8080/api/skills/semantic-query?q=emotional+development"
```

### 🔧 **Troubleshooting**

#### **Port 8080 Already in Use**
```bash
# Find and kill the process on port 8080
lsof -ti:8080 | xargs kill -9
```

#### **Database Connection Issues**
```bash
# Check if databases are running
make status

# Restart databases
make restart-databases

# Full reset if needed (THIS IS DESTRUCTIVE)
make purge-databases
make start-databases
make setup
make setup-medical
```

#### **Weaviate Errors (e.g., vector dimension mismatch)**
This error usually means old data was not cleared properly. A full purge is the best solution.
```bash
make purge-databases
make start-databases
make setup
make setup-medical
```

### 🎯 **API Testing Examples**

#### **Educational Skills Search**
```bash
# Basic semantic queries (requires LLM)
curl "http://localhost:8080/api/skills/semantic-query?q=reading+skills"
curl "http://localhost:8080/api/skills/semantic-query?q=fine+motor+development"
```

#### **Medical Content Search**
```bash
# Search for medical content (requires medical setup)
make search-medical --query "X-Linked Agammaglobulinemia"
```