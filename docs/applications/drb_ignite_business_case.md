Business Case: Data Manager


















Prepared by:  Michael Stewart  
Date: 15 January 2026  
Contact:  michael.stewart@espthinking.co.uk | 07525 359661
Executive Summary
Most educational data is used for administration and accountability. However, to truly drive school improvement, data must move from the spreadsheet to the classroom. This requires a quartet of skills:
1.	Read like a Headteacher: Strategic oversight, accountability, and knowing which operational levers to pull to improve Trust-wide performance.
2.	Interpret like a Scientist: Separating the signal from the noise, identifying the "tells" that precede disengagement, and diagnosing the root causes of attainment gaps.
3.	Investigate like a Business Founder: Breaking down processes and numbers to shine a light on areas to fund where help can be provided.
4.	Implement like an Architect: Building the efficient tools that capture insights without adding to staff workload.
This business case proposes a Data & Insight Manager who delivers the foundational "day job" (MIS, Census, GDPR) while building the Integrated Insight Markbook—a single tool that turns data into a genuine engine for school improvement.

 
Level 1: Operational Excellence (The "Safe Hands")
Clean data is the prerequisite for any scientific insight. I will ensure the Trust’s foundational operations are flawless:
o	MIS Integrity (Arbor/Sims): Ensuring clean, consistent data entry across all 13 schools.
o	Statutory Returns: 100% accuracy on Census, SATs, and financial data submissions.
o	GDPR & Safeguarding: Robust governance of pupil and staff data.
o	Efficiency: Automating reports to give headteachers the "standard" information faster.

 
Level 2: The Integrated Insight Mark-book (The Strategic Value)
I propose the implementation of a unified tool that replaces conventional heavy marking and reporting with real-time, actionable diagnostics. This mark-book combines two powerful frameworks into a single interface:

1.	The Skills Map (The "Sticker Album")
A visual, motivating record of competency acquisition.
o	Workload Displacement: Logging a skill achievement is 5x faster than marking a book. The "data entry" *is* the marking.
o	Live Reporting: Reports are generated automatically from live data, eliminating the "termly report crunch."
o	Primary Advantage: In lower school, where teachers see their cohorts all day, this is the most efficient way to capture the "micro-wins" that lead to success.
2.	The ETP Profile (Predictive Diagnostics)
A "Tone Log" that identifies Character Tells before they become behavioural incidents.
o	Group A (Disengaged): The system spots "Micro-Withdrawals" or "Transition Friction"—the "No Gearbox" signals—allowing for early reconnection.
o	Group B (Effort-Heavy): Identifies students who are working hard but "stalled" in a low gear. They need learning strategies, not more pressure.
o	Group C (Self-Driving): Spots the "Turbocharged" students who are cruising in their comfort zone, allowing us to challenge them with "Optimal Discomfort."

 
## Level 3: The Data Intelligence Infrastructure (The Competitive Advantage)

Beyond the classroom markbook lies a trust-wide analytics layer that transforms fragmented data into actionable intelligence. This is where "Investigate like a Business Founder" becomes reality.

### The Data Lake (Unified Source of Truth)

Currently, each school maintains separate spreadsheets for skills, concerns, finances, and operations. This creates:
- Data silos (no cross-school pattern recognition)
- Duplication and inconsistency
- Compliance risk (scattered sensitive data)
- Lost insights (patterns buried in 13 separate systems)

I propose a **Trust-wide Data Lake** built on Weaviate (vector database) + analysis engine:

**Architecture**:
- Encrypted ingestion from MIS (Arbor/Sims), Skills Markbook, and business forms
- Role-based access control (teachers see only their class; heads see their school; DRB sees trust trends)
- Automated audit logging for GDPR compliance
- Three data tiers:
  - **Raw Data**: Original entries (encrypted, school-segregated)
  - **Analytical Data**: Anonymized for pattern detection (hashed student IDs, no PII)
  - **Actionable Insights**: Conversational flags triggered when thresholds met

**Business Forms Integration** (Same API, Separate Modules):
- Estate Management: Facility condition, maintenance requests, asset tracking
- Finance: Budget, spend, cost-per-outcome analysis
- HR: Staff CPD, expertise inventory, deployment optimization
- Each captures data via forms; all feed the Data Lake for cross-analysis

### Pattern Detection & Trend Analysis

The **Analysis Engine** identifies patterns humans cannot:

**Behavioral Intelligence**:
- Detects Group A cohort signals (e.g., "Classes in Schools 3, 7, and 11 show identical disengagement pattern in Q2")
- Flags resource deployment opportunities (e.g., "These 5 students need Intervention X; cost-benefit suggests group tutoring vs. 1:1")
- Validates fairness (e.g., "Teachers in School 2 award communication skills 40% more than Trust average → bias or genuine excellence?")

**Predictive Intelligence**:
- Identifies students at risk before behavior deteriorates
- Correlates CPD attendance with skill award patterns
- Surfaces cross-school best practices (e.g., "School 5's approach to Group C challenge resulted in 3.2x better high-end skill growth")

**Operational Intelligence**:
- Finance: "Budget for intervention X has 4:1 ROI based on skill-outcome correlation"
- Estate: "Facility issues in 3 schools correlate with higher Group A flags in those locations"
- Staffing: "Expert inventory analysis shows School 10 has untapped CPD capacity; could train others"

### Three-Tier Insight Model

**Tier 1 - Classroom (Daily)**:
- Teacher sees: "[Your student Emma] demonstrated Problem Solving + [Your cohort average skill growth: +15% this term]"
- Action: Real-time feedback for instruction

**Tier 2 - School (Weekly/Termly)**:
- Head sees: "[Your school trends] vs. [Trust average]" + equity dashboard
- Action: Resource allocation, staff development priorities

**Tier 3 - Trust (Termly/Annual)**:
- DRB sees: Anonymized cross-school patterns, investment opportunities, strategic risks
- Action: Trust-wide initiatives, funding decisions, collaboration opportunities

### Security & GDPR-First Design

**Data Protection**:
- Encryption in transit (TLS 1.3) and at rest (AES-256 field-level for PII)
- Segregation by school (breach in School 1 does not expose School 2)
- Automatic de-identification for analysis (cannot re-identify without audit-locked master key)

**Compliance**:
- Data Processing Agreements (DPA) with each school (mandatory GDPR)
- Retention policy: Student data = 5 years post-exit; anonymized insights = indefinite
- Subject access rights: Parents/students can download their data anytime
- Breach protocol: Notify ICO + affected parties within 72 hours
- Privacy impact assessment: Risks + mitigations documented

**Governance**:
- Audit trail for every data access (who, when, what, why)
- Automated GDPR reports (data minimization, purpose limitation checks)
- Role-based access enforced at database level

**CPD as Expertise Library**: Same mechanisms as student skills tracking, but for staff. Teachers log their own CPD attendance and growing competencies. The system creates a trust-wide expertise inventory, linked to professional development resources and session recommendations. Over time, this becomes a powerful tool for succession planning, mentoring pairing, and identifying internal training capacity.

**Trend Detection & Collaborative Flags**: Both teachers and leadership see the data, but thresholds trigger automated escalation. When a concern pattern reaches significance (e.g., 3+ teachers flag the same student with "withdrawal" signals within a 2-week window), the system flags both teachers and leadership for conversation. Multiple observations from different adults create a richer picture than any single teacher could achieve alone.

Why This Matters: Social Gravity & Voltage

By consolidating the Skills Map, ETP Profile, and Data Lake into one integrated system, we can manage the two key determinants of school success:
- **Social Gravity**: Data can identify the "critical mass" of Group A students in a cohort who are pulling others away from engagement. We can deploy resources to "break the gravity" before it becomes systemic. Cross-school analysis reveals systemic patterns (not just classroom problems).
- **Voltage Regulation**: We move from measuring "Attainment" to managing "Optimal Challenge." We ensure every child is pushed regularly against their barriers in a measure of "meaningful discomfort" without creating a culture of anxiety. Trust-wide data ensures no student falls through the cracks.

## Implementation Roadmap: From Pilot to Trust-Wide Scale

### Phase 1: Pilot (School 1, Months 1-4)

**Month 1: Preparation**
- Data governance workshop with School 1 leadership
- Schema design for pilot school's data
- Weaviate instance + analysis engine setup
- GDPR audit and DPA finalization

**Month 2-3: Implementation**
- Skills Markbook deployment (teachers trained)
- Data ingestion from MIS (test run)
- Analysis engine calibration (ensure patterns are valid)
- Stakeholder feedback loops

**Month 4: Refinement**
- Troubleshooting + optimization
- Staff confidence building (evidence of value)
- Training materials completed

**Pilot Success Metrics**:
- 80%+ teacher adoption within 4 weeks
- 3+ actionable insights generated per term
- Zero data breaches or compliance incidents
- Training materials ready for 13-school rollout

### Phase 2: Training & Rollout (Schools 2-13, Months 5-12)

**Training Materials** (Created during pilot):
- Administrator guide (data governance, access control, GDPR)
- Teacher manual (how to use markbook, interpret insights)
- Head dashboard tutorial (reading trust-wide reports, equity analysis)
- GDPR compliance checklist
- Data quality standards and escalation procedures

**Rollout Strategy**:
- 2-3 schools per month (staggered to ensure support quality)
- Dedicated training session at each school
- On-call support during first month at each site
- Cross-school collaboration sessions (learn from early adopters)

**Redundancy & Continuity**:
- Backup Data Manager trained by Month 6 (can take over all operations)
- Automated failover for database (replication across 2 instances)
- Recovery procedures documented and tested quarterly

### Phase 3: Trust-Wide Intelligence (Months 13+)

- Full 13-school data lake operational
- Cross-school pattern analysis enabled
- Trust-wide analytics dashboard live
- Business forms integration (estates, finance, HR) begins
- Annual strategic review cycle using data insights

## Summary: Success Through Insight
I am offering to do the job you need now, while preparing the platform for the future you want:

1. **Support** (Level 1): Ensuring Arbor/MIS/Canvas integrations are frictionless for staff, with clean, accurate data as the foundational prerequisite.
2. **Synthesize** (Level 2): Interpreting the "Noise" in classroom observations to find the real stories—what skills are being built, which students need support, which teachers need CPD.
3. **Strategize** (Level 3): Implementing a low-burden, integrated mark-book backed by a trust-wide Data Lake that turns numbers into actionable intelligence. Business forms, analytics, and cross-school insights—all from a single, secure source of truth. This is where data becomes a genuine lever for school improvement.
4. **Implement like an Architect**: Building systems that work without adding staff burden, with GDPR-compliance and security built in from the foundation.

I'm ready to help DRBIgnite move from managing numbers to shaping lives through data. My approach is to be the Architect of low-risk systems (starting with one pilot school), the Scientist who validates insights against real outcomes, and the Strategist who builds staff confidence through early, tangible wins. Together, we can ensure data serves the mission, not the other way around.

 
Appendix: Credentials
o	QTS & 10+ Years Teaching: I know the classroom reality.
o	University of Birmingham Administrator: Experience with high-level data governance.
o	Founder, ESP Thinking: 12 years of intellectual property on learner diagnostics.
o	AI Solutions Architect (Endorsed 2025): Expert in technical automation and insight generation.  Endorsement from Toby Fotherby 
AWS Senior AI/ML Specialist Solutions Architect (Strategic Accounts)
o	I recently designed and led an AI Solutions Architecture program, that Michael enrolled for. He quickly stood out as one of the most thoughtful and capable practitioners that I’ve worked with.

Michael consistently demonstrated a strong grasp of AI architecture fundamentals, along with a highly effective, real-world approach to designing and delivering enterprise-grade AI solutions. He showed mastery of all key concepts in the program—from model evaluation and orchestration to scalability and responsible AI principles.

Additionally, I am deeply impressed with the innovative AI-powered teaching strategy and framework he developed prior-to and during the course. It’s an outstanding example of how AI can be applied not only to solve technical problems, but also to transform the way people learn. His approach is intuitive, adaptive, and insightful—providing personalized, efficient training that supports both advanced learners and those with unique challenges.

Michael brings together technical depth, creativity, and empathy in a way that’s rare. I’m confident he will make a significant impact in any AI or education-focused initiative he chooses to pursue.
o	

 
*I don't just file the story. I help write the next chapter.*

