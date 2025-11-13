# ESP Organizer Project Structure

## Overview

This project consists of a React frontend and a Go backend, working together to organize and process medical education content. The application uses MongoDB for document storage and Weaviate for vector search.

## Directory Structure

.
├── backend
│   └── internal
│       └── integration
│           ├── chisg_client.go
│           └── chisg_types.go
├── cmd
│   └── list-project
│       └── main.go
├── docs
│   ├── Action Plan_ HumanOS Ecosystem Develop.md
│   ├── CHISG project and progress
│   ├── CHISG Project.md
│   ├── CHISG_CLIENT_EXPLANATION.md
│   └── PROJECT_STATUS.md
├── esp-organizer
│   ├── backend
│   │   └── internal
│   │       └── integration
│   │           └── chisg_client.go
│   ├── backups
│   │   ├── mongodb
│   │   └── weaviate
│   ├── byteops
│   │   └── byteops.go
│   ├── cmd
│   │   ├── api-server
│   │   │   └── main.go
│   │   ├── aws-diagnostics
│   │   │   └── main.go
│   │   ├── book-create-if-missing
│   │   │   └── main.go
│   │   ├── check-data
│   │   │   └── main.go
│   │   ├── check-json
│   │   │   └── main.go
│   │   ├── check-medical
│   │   │   └── main.go
│   │   ├── check-routes
│   │   │   └── main.go
│   │   ├── clear-db
│   │   │   └── main.go
│   │   ├── clear-educational
│   │   │   └── main.go
│   │   ├── clear-medical
│   │   │   └── main.go
│   │   ├── clear-medical-mock
│   │   │   └── main.go
│   │   ├── debug-bedrock
│   │   │   └── main.go
│   │   ├── debug-embeddings
│   │   │   └── main.go
│   │   ├── debug-relationships
│   │   │   └── main.go
│   │   ├── demo-working-features
│   │   │   └── main.go
│   │   ├── download-adobe-pdfs
│   │   │   └── main.go
│   │   ├── extract-semantic-links
│   │   │   └── main.go
│   │   ├── fix-medical-vectors
│   │   │   └── main.go
│   │   ├── fix-skills-relationships
│   │   ├── migrate-semantic-links
│   │   │   └── main.go
│   │   ├── pipeline-monitor
│   │   │   └── main.go
│   │   ├── process-medical-pdf
│   │   │   └── main.go
│   │   ├── process-pdf
│   │   │   └── main.go
│   │   ├── query-medical-content
│   │   │   └── main.go
│   │   ├── query-test
│   │   │   └── main.go
│   │   ├── reprocess-job
│   │   │   └── main.go
│   │   ├── revectorize
│   │   │   └── main.go
│   │   ├── rollback-batch
│   │   │   └── main.go
│   │   ├── run-full-test
│   │   │   └── main.go
│   │   ├── search-medical
│   │   │   └── main.go
│   │   ├── seed-test-data
│   │   │   └── main.go
│   │   ├── server
│   │   │   └── main.go
│   │   ├── setup-hsg-schema
│   │   │   └── main.go
│   │   ├── setup-weaviate-schema
│   │   │   └── main.go
│   │   ├── test-hsg-search
│   │   │   └── main.go
│   │   ├── test-medical-vectorization
│   │   │   └── main.go
│   │   ├── test-textract
│   │   │   └── main.go
│   │   ├── upload
│   │   ├── upload-skills
│   │   │   └── main.go
│   │   ├── utility
│   │   ├── validate-curated-links
│   │   │   └── main.go
│   │   ├── validate-llm-extraction
│   │   │   └── main.go
│   │   ├── verify-aws-config
│   │   │   └── main.go
│   │   ├── verify-routes
│   │   │   └── main.go
│   │   ├── weaviate-add-property
│   │   │   └── main.go
│   │   ├── weaviate-reset
│   │   │   └── main.go
│   │   └── go_version_test.go
│   ├── config
│   │   └── config.yaml
│   ├── data
│   │   ├── mongodb
│   │   └── weaviate
│   ├── docker
│   │   ├── Dockerfile
│   │   └── nginx.Dockerfile
│   ├── docs
│   │   ├── hsg_threads
│   │   ├── semantic_links
│   │   └── training_data
│   ├── frontend
│   │   ├── build
│   │   ├── public
│   │   ├── src
│   │   └── ...
│   ├── internal
│   │   ├── api
│   │   ├── InfoFlow
│   │   ├── models
│   │   └── ...
│   ├── logs
│   ├── resources
│   ├── scripts
│   └── ...
├── ...

