# Product Vision: The Integrated Learning Ecosystem

## The Strategic Vision
To build an **Adaptive Learning Architecture** with multiple interconnected pathways, moving beyond isolated features to create a coherent, all-in-one educational operating system.

---

## Core Product Pillars

### 1. The Coach Product (GCSE Revision Tool)
A knowledge development and tracking tool designed to make learning content as efficient as possible. It is built on the CHISG knowledge graph (Semantic Links).

**Core Features**:
-   **Lesson Overviews & Narratives**: AI-generated summaries and contextual stories to explain complex topics simply.
-   **Spreeder Mode**: A speed-reading interface that flashes words one at a time to increase content consumption speed and enable rapid repetition of key concepts.
-   **Downloadable Audio Summaries**: On-demand text-to-speech generation of lesson overviews, creating downloadable MP3 files for on-the-go learning.
-   **Automated Glossaries**: Automatic extraction and definition of key terms for any given subject.
-   **Structured Question Practice**: An interface for students to practice long-form and multi-step questions, with AI-assisted feedback.
-   **Curriculum-Aware Study Planner**: An intelligent timetable that holds the official curriculum, tracks progress against it, and schedules practice sessions using spaced repetition.
-   **Curriculum Intelligence (Advanced)**: Utilizes CHISG and an LLM to analyze the curriculum itself, identifying conceptual gaps, optimal learning pathways, and prerequisite chains.

---

### GCSE Revision Tool: Review Modes & Philosophy

**Data Model**:  
- All content (in MySQL or MongoDB) is structured as:
  - **Keywords & Definitions**
  - **Facts**
  - **Narratives/Overviews**

**Review Modes**:
- **45-second Spreeder**: Rapid, low-pressure review of key facts or overviews (words/phrases flashed sequentially).
- **3-minute Audio Overview**: Text-to-speech (e.g., AWS Polly) of the narrative, for passive listening.
- **10 Keyword-Definition Drill**: Quick recall of essential terms.
- **Narrative Overviews**: Readable, AI-generated summaries for context.
- **Multiple Choice Questions**: Lightweight, formative assessment.

**Philosophy**:  
- **"Little and Often"**: Sessions are intentionally short and non-intimidating to lower barriers to entry and maximize cost-benefit for the student.
- **Low Intimidation, High Value**: The tool is designed to make revision feel achievable, not overwhelming, encouraging frequent, bite-sized practice.

**Additional Review Modes (for future consideration):**
- Spaced repetition flashcards (offline-capable)
- Cloze (fill-in-the-blank) exercises
- Short written recall (summarize in your own words)
- Confidence-based self-assessment
- Peer-generated questions (if social features added)
- Visual summaries (mind maps, diagrams)
- Mini-quizzes (2-3 questions per session)
- Progress streaks and gamification

---

#### Points for Discussion / Further Review Options

- **Other Review Modes**:  
  - Flashcards (spaced repetition)
  - Fill-in-the-blank (cloze) exercises
  - Peer-generated questions
  - Visual summaries (mind maps, diagrams)
  - Confidence-based self-assessment (rate your recall before/after)
- **Gamification**:  
  - Streaks, points, badges for consistency
- **Personalization**:  
  - Adaptive review based on past performance
- **Social/Community**:  
  - Group challenges, shared progress

---

### 2. Skills Tree Rising (STR)
A skill development and tracking tool designed for competency-based training and assessment. It is built on the Skills Map model.

**Core Features**:
-   **Dynamic Skill Mapping**: Automatically updates and adjusts the Skills Map based on student interactions and progress.
-   **Competency-Based Progression**: Students advance upon mastering skills, not time spent on a topic.
-   **Integrated Assessment Tools**: Seamless incorporation of formative and summative assessments linked to skill development.
-   **Real-Time Feedback Loops**: Immediate, actionable feedback for students and educators to guide learning.
-   **Customizable Learning Pathways**: Tailored educational experiences that adapt to individual student needs and goals.

#### Related Internal Tool: The Skills Collector
The content for Skills Tree Rising is managed via a dedicated internal authoring tool known as the **Skills Collector** (`skills-map-platform`). This application allows team members and curriculum designers to visually create skills, define relationships, build skill trees, and associate them with courses. It is the content management system (CMS) for our skills-based products.

### 3. Teacher Development 2.0
-   **Collaborative Skill Mapping**: Teachers use the same Skills Map as students to model their own professional growth, creating peer-led development communities.
-   **Teachers as Innovators**: The platform provides tools for teachers to design and share their own effective learning exercises.

### 4. Parent Ecosystem Strategy
-   **Transparency & Empowerment**: Parents get real-time access to their child's Skills Map, with linked resources on how to support learning at home.
-   **Community Cooperatives**: A vision for institution-partnered parent cooperatives for childcare and learning support, reducing costs and improving communication.

---

## The "Magic Wand" Tablet: A Complete Solution

### The Vision
A low-cost, managed Android tablet that runs our educational OS exclusively. This isn't about selling hardware; it's about controlling the environment to deliver a perfect learning experience.

### It Solves Three Problems at Once:
1.  **Device Management Chaos**: Ends the "bring your own device" and mobile phone battles.
2.  **AI Anxiety**: Provides a safe, walled-garden environment to teach AI literacy.
3.  **Budget Pressure**: Eliminates massive photocopying and textbook costs, with a clear ROI in 12-18 months.

### Target Specifications
-   **Screen**: 8-10" HD
-   **RAM**: 4GB
-   **Storage**: 32GB + microSD
-   **Battery**: All-day (6000-7000mAh)
-   **Durability**: Reinforced for classroom use
-   **OS**: Android Enterprise for school device management
-   **Cost**: Target bulk price of ~$150/unit.

---

## AI as a Skills Coach: The Phased Rollout

1.  **Phase 1: Manual Mapping**: We start by manually mapping common student responses and errors to our psychological frameworks, building a rich dataset.
2.  **Phase 2: Pattern Recognition**: The AI learns from this data to detect patterns and automatically suggest interventions.
3.  **Phase 3: Emergent Strategy**: The AI begins to identify novel learning patterns and strategies that even expert teachers might miss, creating a system that learns and improves as it teaches.
4.  **Phase 4: Self-Evaluation & Refinement**: The system uses LLMs to evaluate its own performance against two key metrics:
    *   **Extraction Quality**: It measures its ability to transform new, unstructured text into high-quality semantic links by comparing its output against a "golden set" of demonstrative examples. This creates a feedback loop for improving its "understanding".
    *   **Groundedness (Hallucination Prevention)**: It measures how well its final answers are grounded in the structured knowledge provided by CHISG, ensuring it doesn't invent information.

---

### Mobile App Architecture (React Native)

- **Local-first, Offline-first**: All progress, lesson completion, and core content stored locally (using MongoDB Realm or SQLite).
- **Audio files**: Download on demand, cache locally, allow user to clear cache.
- **LLM Feedback**:  
  - Offline: Use rules-based encouragement and feedback.
  - Online: Sync progress and use cloud LLM for richer feedback.
- **Sync**: When online, sync progress and answers to the cloud for backup and analytics.

---

### Subject Chat Mode & Engagement-First Design

**Conversational LLM Review ("Subject Chat Mode")**
- The LLM can operate in a "chat" mode for each subject or topic.
- Instead of only delivering facts or questions, the LLM engages the student in a lightweight, friendly conversation about the topic.
- The chat can include:
  - Explaining concepts in a conversational way
  - Asking the student simple, low-pressure questions ("What do you remember about X?")
  - Offering encouragement and positive feedback
  - Occasionally slipping in a quick quiz or challenge, but always in a supportive, non-intimidating manner

**Engagement Scoring & Adaptive Activity**
- Each interaction is graded against an "engagement score" (e.g., response time, message length, voluntary participation, sentiment).
- The system starts with the absolute minimum effort required to get a "win" (e.g., just replying at all).
- As engagement increases, the LLM gradually introduces more challenge or depth, but always within the same short time frame (e.g., 10 minutes).
- The goal is to maximize the reward-to-effort ratio, especially for students who are resistant to traditional homework or revision.
- Gamification elements (streaks, badges, progress bars) are layered on top to make even minimal participation feel rewarding.

**Philosophy: "Little and Often, Least Effort to Start"**
- The system is designed for students who find homework or revision antithetical to their habits.
- The initial barrier to entry is extremely low—just open the app and reply to a message.
- Over time, the system "sneaks in" more engagement, but always keeps sessions short and positive.
- The aim is to make 10 minutes on the bus or wherever feel like a complete, valuable revision session.

**Future Directions**
- Adaptive chat personalities (serious, funny, motivational, etc.)
- Dynamic adjustment of chat/question ratio based on engagement trends
- Integration with other review/test modes for seamless transitions

---

### Should This Be a Separate Project?

- If the mobile app is a major standalone product, create a new `/mobile` directory or repo.
- If it’s just a new UI for the same backend, keep it in the monorepo.

---
