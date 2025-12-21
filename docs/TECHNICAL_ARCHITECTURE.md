# ESP Organizer - Technical Architecture

**Purpose**: Single source of truth for system architecture, data models, and integration patterns.

**Last Updated**: 2025-01-27
**Version**: 1.3 - Production Infrastructure & Routing

---

## Table of Contents

1.  [Monorepo Structure](#monorepo-structure)
2.  [System Overview](#system-overview)
3.  [Core Components](#core-components)
4.  [Conceptual Data Models](#conceptual-data-models)
5.  [Data Storage Architecture](#data-storage-architecture)
6.  [Production Infrastructure (Vultr)](#production-infrastructure-vultr)
7.  [Deployment Architecture](#deployment-architecture)
8.  [Architectural Evolution: The Modular Monolith](#architectural-evolution-the-modular-monolith)
9.  [AI Self-Evaluation Architecture](#ai-self-evaluation-architecture)

---

## Monorepo Structure

The entire HumanOS ecosystem is developed within a single monorepo located at `/assist`. This is a strategic decision to eliminate code duplication and documentation fragmentation.

-   `/assist/docs/`: **The single source of truth.** All project documentation (vision, roadmap, technical, daily logs) lives here.
-   `/assist/esp-organizer/`: The monolithic Go backend that serves all products and features.
-   `/assist/frontend/`: The single React frontend application that contains the UI for all products.
-   `/assist/skills-map-platform/`: The self-contained "Skills Collector" application (React frontend, Go backend, and configs).
-   `/assist/tools/`: Standalone Go utilities for seeding data, running migrations, etc.
-   `/assist/docker-compose.yml`: The single entry point for running the entire stack locally.

All new products, like the "Revision Assistant" (Coach Tool), will be built as new modules *within* this structure, not in new top-level folders.

### GitHub Repository Strategy

To support the monorepo structure, we will adhere to the following repository strategy:
-   **Single Source of Truth**: The `github.com/m5tew/assist` repository is the one and only active repository for this entire ecosystem.
-   **Archive Legacy Repositories**: All other related repositories (e.g., the original `skills-map-platform`, `skillsdocker`, `HumanOS`, `espdata`) should be **archived** on GitHub. This makes them read-only, preserving their history while clearly directing all new contributions to the `assist` monorepo.
-   **No `.gitignore` for Sub-projects**: We will not use `.gitignore` to exclude sub-projects like `skills-map-platform`. All code within the `/assist` directory should be committed to the `assist` repository.

---

## System Overview

The HumanOS ecosystem is a multi-layered platform designed to provide psychologically-aware, AI-driven educational tools. It is built as a monolithic Go backend API serving a React frontend, with specialized data stores for different types of information.

### The Three-Layer Stack

-   **Presentation Layer**: React frontend, served as static assets via Nginx.
-   **Application Layer**: Go backend (`esp-organizer`) providing a unified API for all features. This includes the HumanOS core (psychological frameworks) and the CHISG (knowledge graph) services.
-   **Data Layer**: A hybrid model using Weaviate for vector search, MongoDB for document and metadata storage, and MySQL for structured relational data from the Skills Map platform.

---

## Core Components

-   **Frontend**: React, TypeScript, Material-UI
-   **API**: Go, Gorilla Mux (`esp-organizer` monorepo backend)
-   **Vector Database**: Weaviate (for semantic search and graph relationships)
-   **Metadata Store**: MongoDB (for raw documents, user data, extraction samples)
-   **Relational Database**: MySQL (for Skills Map data, student records, course structures)
-   **LLM & Embedding Services**: AWS Bedrock (Claude for generation, Titan for embeddings)
-   **Document Processing**: AWS Textract (for OCR on scanned PDFs)
-   **Deployment**: Docker, Docker Compose, Nginx
-   **Internal Authoring Tool**: **Skills Collector** (`skills-map-platform`), a separate application for internal teams to create and manage skills, courses, and their relationships.

---

## Conceptual Data Models

This section defines the core data concepts that power the ecosystem.

#### 1. Semantic Links: The Atoms of Knowledge

-   **Purpose**: To represent a single, atomic relationship between two pieces of knowledge. This is the backbone of the **Coach Product (GCSE Tool)**.
-   **Definition**: A data object capturing a relationship like `(Source Term) --[Relation Type]--> (Target Term)`.
-   **Example**: `("X-Linked Agammaglobulinemia") --[causes]--> ("B-cell deficiency")`.
-   **Use Case**: Powers the knowledge graph (CHISG) for answering questions, generating summaries, and explaining concepts.

#### 2. Skills Map: The Pathways of Competency

-   **Purpose**: To model skill development hierarchies with prerequisites and outcomes. This is the backbone of **Skills Tree Rising**.
-   **Definition**: A specialized graph where nodes are skills, and edges represent learning dependencies. It is a specialized application of the semantic link concept.
-   **Structure**: A Skill node contains properties like `name`, `description`, `parent_skills` (prerequisites), `offspring_skills` (what this skill enables), and `assessment_criteria`.
-   **Example**: `("Fine Motor Skills") --[is_prerequisite_for]--> ("Handwriting")`.
-   **Use Case**: Powers skill tracking, personalized learning plans, and competency-based assessment.

#### 3. The Unified Vision

The long-term goal is to merge these two concepts into a single, unified graph that tracks both **knowledge** and its practical application as **skills**.

---

## Data Storage Architecture

The project uses a hybrid data store model, using the right database for the right job.

-   **Weaviate**: Stores `SemanticLink` vectors. Used for fast, semantic-based queries to find related concepts. The "what" and "how" of knowledge.
-   **MongoDB**: Stores raw `Document` content, `SummaryChunks`, `ExtractionJobs`, and other application metadata. The "source of truth" for content.
-   **MySQL**: Stores structured, relational data for the `skills-map-platform`, including `individuals` (users), `courses`, and `skill_link` relationships. The "who" and "where" of the learning structure.

---

## Production Infrastructure (Vultr)

This section documents the setup of our live production server. The architecture uses a main Nginx reverse proxy to route traffic to different applications running in Docker containers.

-   **Provider**: Vultr
-   **Instance Type**: High Frequency Compute
-   **Location**: London
-   **Operating System**: Alpine Linux
-   **Key Software**: Docker, Docker Compose, OpenSSH Server
-   **Primary Directory**: All project files are intended to be managed from the `~/assist/` directory.

### Vultr Container Architecture

-   **Nginx Reverse Proxy (Main)**: A top-level Nginx container that listens on ports 80/443. It routes traffic based on the domain/subdomain to the appropriate application container. This is the single entry point to the server from the outside world.
-   **`assist` Application**:
    -   `assist-frontend`: The Nginx container serving the main application's React frontend.
    -   `assist-api`: The `esp-organizer` Go backend container.
-   **`skills-map-platform` Application (Skills Collector)**:
    -   `skills-collector-frontend`: A separate Nginx container serving the admin/authoring tool's React frontend.
    -   `skills-api`: The backend container for the Skills Collector application.
-   **Databases**: All applications connect to the shared set of database containers (Weaviate, MongoDB, MySQL).

-   **Domain/URL**: Currently accessed via IP address. A domain name (e.g., `esp-world.co.uk`) and subdomains (e.g., `admin.esp-world.co.uk`) will be configured via DNS to point to the server's IP, with the reverse proxy handling routing.

---

## Go Module Strategy

We use a **Go Workspace (`go.work`)** for local development and rely on **canonical module paths (`go.mod`)** for production builds and sharing.

-   **Local Development (`go.work`)**: Enables simultaneous development across multiple local modules (`esp-organizer` and `tools`). The `use` directive tells the Go compiler to use the local directories instead of fetching from GitHub, allowing for rapid iteration.
-   **Production & Sharing (`go.mod`)**: Each module (`esp-organizer`, `tools`) has a `go.mod` file that declares its official import path (e.g., `github.com/m5tew/assist/esp-organizer`). This ensures reproducible builds in CI/CD pipelines and allows other projects to import them correctly.

---

## Deployment Architecture

This section outlines the process for running the project locally and deploying it to the production server.

**For detailed, step-by-step commands for both full and partial deployments, see the `DEPLOYMENT_PLAYBOOK.md` document.**

### Core Principles
- **Build Natively for Local, Cross-Compile for Production**: We build images for our native architecture (`arm64` on an M-series Mac) for local development. We explicitly build for `linux/amd64` when creating images for the Vultr server.
- **Container Registry as Intermediary**: Docker Hub is the central repository for our production-ready images. The server only pulls from this registry.

### Local Development Setup (Running on your Mac)
These steps are for running the entire `assist` application stack on your local machine using the main `docker-compose.yml`.

**1. Build Local Images (Run on your Mac)**
Build images for each service without specifying a platform. Docker will automatically build for your Mac's architecture (`arm64`).
```bash
# Build the assist-api for your Mac
cd /Users/michaelstewart/Coding/assist/esp-organizer
docker build -t mike5tew/assist-api:latest .

# Build the assist-frontend for your Mac
cd /Users/michaelstewart/Coding/assist/frontend
docker build -t mike5tew/assist-frontend:latest .
```

**2. Run Locally (Run on your Mac)**
From the root of the `assist` project, use `docker compose` to start all services.
```bash
# (Run from /Users/michaelstewart/Coding/assist/)
docker compose up -d
```

---

### Production Deployment to Vultr

The production environment on Vultr is managed via a single `docker-compose.yml` file located in `/opt/esp`. The process involves:
1.  Building `linux/amd64` images for all services on a local machine.
2.  Pushing these images to Docker Hub.
3.  Syncing configuration files (`docker-compose.prod.yml`, `.env`, etc.) to the server.
4.  Running `docker compose` commands on the server to pull the new images and restart the necessary services.

**For the exact commands for full and partial updates, refer to `docs/DEPLOYMENT_PLAYBOOK.md`.**

### Frontend Routing (SPA 404 Fix)
To solve the issue where refreshing a page like `/login` results in a 404 error, the Nginx server is configured with a `try_files` directive. This tells Nginx to first look for a file that matches the URL, then a directory, and if neither is found, to fall back to serving `/index.html`. This allows the React Router to handle the client-side routing.

---

## Architectural Evolution: The Modular Monolith

The current architecture is a **Modular Monolith**. This is a deliberate choice to prioritize development speed and operational simplicity while the product is in its early stages.

-   **Current State**: All backend logic resides within the single `esp-organizer` Go service. However, it is internally structured into distinct packages representing logical business domains (e.g., `auth`, `knowledge`, `skills`, `barriers`).
-   **Future State**: As the system grows and specific domains become complex or require independent scaling, these internal modules can be extracted into separate microservices. The clear boundaries established now will make this future migration significantly easier.
-   **The Trigger for Change**: We will consider moving to a microservices architecture when we experience specific scaling bottlenecks, or when team size grows to a point where independent service development becomes more efficient. We will not adopt microservices prematurely.

---

## AI Self-Evaluation Architecture

To ensure the CHISG system is both accurate and reliable, we will implement an LLM-based evaluation framework that assesses its performance on two critical axes.

### 1. Extraction Quality Assessment

This process evaluates how well the system can transform unstructured text into high-quality semantic links and metadata.

**Architecture**:
1.  **Golden Set**: A curated MongoDB collection of high-quality `(input_text, expected_semantic_links)` pairs that serve as demonstrative examples. This set is manually created by referencing curriculum-aligned textbooks, but **does not store the source textbook text itself**, only the derived links and the minimal text snippets they connect.
2.  **Extraction**: The system processes a new piece of text to generate a set of candidate semantic links.
3.  **LLM-as-Judge**: An evaluator LLM is prompted to compare the candidate links against the style, structure, and quality of the "Golden Set". It answers questions like, "Is the granularity of these links consistent with the examples?" and "Is the relationship type appropriate?".
4.  **Scoring & Refinement**: The evaluator's feedback generates a quality score. High-scoring outputs become candidates for inclusion in the Golden Set, creating a continuous improvement loop.

### 2. Groundedness Assessment (Hallucination Prevention)

This process evaluates how well the final, user-facing answer is grounded in the factual context provided by the CHISG knowledge graph.

**Architecture**:
1.  **Context Generation**: For a given user query, CHISG assembles a structured knowledge context (a set of relevant semantic links).
2.  **Answer Generation**: A primary LLM generates a user-facing answer, with a strict instruction to *only* use the information provided in the CHISG context.
3.  **LLM-as-Cross-Examiner**: A second, independent evaluator LLM is given the final answer and the original CHISG context. It is asked a simple question: "Does this answer contain any substantive information not present in the provided context?"
4.  **Groundedness Score**: The evaluator's "yes/no" response provides a clear, binary score for hallucination, which can be tracked over time to measure the system's reliability.

