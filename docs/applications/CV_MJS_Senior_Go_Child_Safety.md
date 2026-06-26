# Michael Stewart

**69 Dovecote Road, Bromsgrove, B61 7BP**  
07525 359661 | michael.stewart@espthinking.co.uk  
[linkedin.com/in/michael-stewart-espthinking](https://www.linkedin.com/in/michael-stewart-espthinking/)  
Portfolio: [espthinking.co.uk](https://espthinking.co.uk)

---

## Professional Summary

Full-stack engineer specialising in **Go microservices** and **child wellbeing systems**, with 12+ years of classroom experience informing every technical decision. I architect and build cloud-native platforms in Go, React/TypeScript, and SQL — but my differentiator is deep domain expertise in child development, having independently designed and built an integrated child observation and skills framework (12 biological spectra, semantic knowledge graph, assessment pipeline) used to map child wellbeing across home and school contexts. I bring the rare combination of someone who can both build the backend and explain why the data model matters for safeguarding outcomes.

---

## Technical Skills

| Category | Technologies |
|----------|-------------|
| **Backend** | Go (Golang) — primary language. gorilla/mux, go-chi, MongoDB driver, Weaviate client, JWT auth, REST API design |
| **Frontend** | React, TypeScript, Material UI, React Native (Expo), HTML/CSS |
| **Databases** | MySQL (212-table schema, complex multi-table JOINs), MongoDB, SQLite, Weaviate (vector DB) |
| **Cloud & DevOps** | AWS (Bedrock, S3, EC2), Docker, docker-compose, Vultr VPS, Nginx reverse proxy, Makefile-based build automation, shell scripting |
| **Architecture** | RESTful microservices, semantic search pipelines, RAG (retrieval-augmented generation), Weaviate knowledge graph |
| **AI/ML** | LLM orchestration (AWS Bedrock/Claude), vector embeddings, semantic extraction pipelines, data enrichment workflows |
| **Other** | Python (scripting/ML pipelines), WinDev Mobile, PHP, ASP |

---

## Professional Experience

### Director & Lead Engineer | ESP Thinking Ltd | Jan 2016 – Present

Founded and built an integrated child wellbeing technology platform — designing, architecting, and coding the full stack. All backend services written in **Go**, all frontends in **React/TypeScript**. Earlier products (ESP Seating, ESP Behaviour Lite) built in **WinDev Mobile** and published to iOS and Android (2018).

**Go Backend Services (5 codebases):**

- **ESP Organizer API** (Go, gorilla/mux, MongoDB, Weaviate, AWS Bedrock) — Core platform service powering the child observation and skills engine. RESTful endpoints for student profiles, skill tracking, ETP (Emotional Trigger Points) assessment, and semantic search. Integrates MongoDB for document storage and Weaviate vector DB for knowledge graph queries.

- **Skills Map Platform API** (Go 1.24, gorilla/mux, MySQL, JWT) — CHISG-aligned skills tracking system with role-based auth and competency mapping. 212-table MySQL schema with complex multi-table JOINs across resources, competencies, skills, and markbook data. Serves structured skill data (1,500+ CHISG-linked elements) to a React frontend via REST endpoints.

- **HumanOS API** (Go 1.23, go-chi, Weaviate, AWS Bedrock) — Backend for the biological profiling engine mapping 12 developmental spectra. Consumes LLM outputs via AWS Bedrock for semantic enrichment. Weaviate schema design for vector similarity search across child development observations.

- **DRB Monitor API** (Go 1.23, MongoDB) — Multi-Academy Trust estates oversight tool aggregating financial and asset data. MongoDB aggregation pipelines for cross-school reporting.

- **ESP Data API** (Go, gorilla/mux) — Science education content service with RESTful endpoints, serving structured curriculum data to mobile and web clients.

**React/TypeScript Frontends:**

- **Portfolio & Product Site** (React, TypeScript, MUI) — 7 product landing pages (ToddlerOS, PrimaryOS, CareerOS, ETP, CHISG, LAO, ESP World) with SEO, responsive design, and nested routing. Deployed via Docker to Vultr VPS.

- **Skills Map Frontend** (React, TypeScript) — Interactive skill tree explorer, graph visualisation, competency assignment interface, and export tools.

- **DRB Frontend** (React, TypeScript, Vite) — Dashboard for estate management data with real-time monitoring views.

- **ToddlerOS Mobile App** (React Native, Expo SDK 54, TypeScript) — Parent-facing child observation app. 36 themed activity weeks, 12-spectrum sensor questions, local-first storage with `AsyncStorage`. Published to iOS via EAS Build.

- **LAO Mobile App** (React Native, Expo, TypeScript) — GCSE Science revision app with speed reading, flashcards, and structured content delivery. Published to iOS App Store.

- **ESP Seating** (WinDev Mobile, 2018) — Timetable management, seating plans, and focus registers for classroom logistics. Published to iOS and Android.

- **ESP Behaviour Lite** (WinDev Mobile, 2018) — Gamified behaviour management app with flick-action recorder for real-time scoring, synced via Firebase to classroom display. Published to iOS and Android.

**Data & ML Pipelines:**

- Built a CHISG extraction pipeline (Python, AWS Bedrock/Claude) processing 1,500+ educational elements through LLM-powered semantic analysis, outputting structured JSON for Weaviate ingestion.
- Built a RAG pipeline for case study ingestion — extracting, chunking, and storing case study content in Weaviate for retrieval-augmented queries via Go endpoints.
- Designed Weaviate schemas for semantic search across skills, observations, and curriculum content — enabling vector similarity queries across the knowledge graph.
- Seeded and maintained vector databases across development environments.

**Infrastructure:**

- Docker and docker-compose for all services (dev and production configurations).
- Vultr VPS deployment with Nginx reverse proxy, SSL termination, environment management.
- Makefile-based build automation across all Go and frontend projects.

**Child Wellbeing Domain Work:**

- Designed the **12-spectrum ETP framework** — a biological observation model mapping child development across spectra like Social Gravity, Voltage Sensitivity, Risk Tolerance, and Mirror Neuron Tuning. Each spectrum has response matrices, recommended language patterns, and assessment sensors.
- Built a **36-week activity curriculum** (ToddlerOS) with 3 difficulty cycles per spectrum, sticker-based progress tracking, and parent-facing observation prompts.
- Designed a **Three-Check Sensor** assessment loop (Recovery, Friction, Agency) with noise filtering to distinguish genuine developmental patterns from transient states.
- Created the **PrimaryOS framework** bridging classroom behaviour observation to the same 12-spectrum model — giving teachers and parents a shared vocabulary.
- Architected a **7-capacity foundational skills model** (inhibitory control, joint attention, emotional labelling, turn-taking, sustained attention, verbal interaction density, frustration tolerance) mapped to CHISG skill definitions.

---

### Electrical Labourer | Sandunn Projects Ltd | Sep 2024 – Present (intermittent)

Cable installation, tray/trunking/conduit work, and general site support on commercial electrical projects including University of Warwick. Maintained ESP Thinking development in parallel.

---

### ISU Administrator | University of Birmingham | Mar 2023 – Jan 2024

Managed deployment of role-players for medical communication skills training across the Institute of Clinical Sciences. High-paced environment coordinating classes, exams, and individual coaching sessions. Stepped in as facilitator and role-player when needed.

---

### Science Teacher (QTS) | Various Schools, West Midlands | Sep 2010 – Dec 2022

12+ years teaching Science across Key Stages 2–5, including:
- **NQT year** at OFSTED 'Outstanding' rated Kenilworth School
- Physics to A-Level at St George's School, Edgbaston
- Inner-city teaching at Lordswood Boys School and Eden Boys East
- Long-term cover at Ridgeway Secondary and Walkwood Middle

This classroom experience directly informs the child observation frameworks I now build in software — every spectrum, every activity, and every assessment sensor originates from watching thousands of children reveal their wiring through their behaviour.

---

### Director | David Manners Ltd | Jul 2000 – Aug 2010

Grew subsidiary company (Abingdon MG Parts) from **£230K to £3M turnover in 5 years**. Director of two companies within a group of five (50+ employees). Designed warehouse systems, built cataloguing and stock management software (WinDev, PHP, ASP), and created the company website. Rapid M&A integration of acquired business within one month.

---

### Content Producer | Harvey Project, Wayne State University | Jan 2000 – Jul 2000

Produced animations and distance-learning material for an online physiology resource. Presented to professors from across the US at Stanford University on teaching physiology via the web.

---

## Education

| Qualification | Institution | Date |
|--------------|-------------|------|
| **AI Solutions Architect** | Elvtr | Nov 2024 – Feb 2025 |
| **Science Teacher GTP (QTS)** | Titan Partnership, Birmingham | Sep 2011 – Aug 2012 |
| **Biophysics PhD** (research completed) | National Institute for Medical Research | Sep 1996 – Mar 2000 |
| **Biochemistry BSc (2:1)** | University College London | Sep 1993 – Aug 1996 |

QTS Number: 10/82800 | Enhanced DBS: 001755820981

---

## Publication

Jones, S., **Stewart, M.**, Michie, A., Swindells, M.B., Orengo, C., & Thornton, J.M. (1998). "Domain assignment for protein structures using a consensus approach: Characterization and analysis." *Protein Science*, 7(2), 233–242.

---

## Why This Role

Your platform delivers child safety and wellbeing solutions. I have spent years building exactly that — not as a side interest, but as 5 Go backends, 5 React frontends, a 12-spectrum observation framework, a 36-week activity curriculum, a semantic knowledge graph, and a mobile app in the App Store. I bring both the engineering skills and the domain expertise to contribute from day one — and the teaching background to explain why the data model matters for safeguarding outcomes.

---

*References available on request.*
