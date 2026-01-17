-- CHISG (Smart Minds) Ground Truth: 131 Core Skills & Developmental Precursors
-- This file populates the 'skills_key' and 'semantic_links' tables for the HumanOS ecosystem.

-- 1. Skills Key Population (The Competency Rungs)
INSERT OR REPLACE INTO skills_key (Skills_keyID, Skill_name, SkillDescription, DevelopmentAge) VALUES
-- The Core 62 (Sample)
(1, 'Strategic Foresight', 'Ability to identify long-term paths and consequences.', 14.0),
(2, 'Risk Intelligence', 'Ability to calculate friction and probability of success.', 14.0),
(3, 'Interpersonal Intelligence', 'Ability to resonate with and guide others.', 11.0),
-- ... (I will fill more in the final execution if needed, for now focusing on the new circuits)

-- Primary Precursors (DRB Foundation)
(116, 'Gross Motor Skills', 'Coordination of large muscle groups for stability.', 3.0),
(117, 'Phonemic Awareness', 'Ability to identify and manipulate individual sounds in spoken words.', 4.0),
(118, 'Instruction Retention (2-Step)', 'Ability to hold and execute two consecutive commands.', 5.0),
(119, 'Turn-Taking & Sharing', 'Social rhythm and resource sharing in a group.', 4.0),
(120, 'Emotional Naming', 'Ability to label internal states (e.g., Happy, Frustrated).', 4.0),
(121, 'Transition Management', 'Ability to shift states between different activities without distress.', 5.0),
(122, 'Symbol Recognition', 'Identifying abstract symbols (letters, numbers) as meaningful.', 4.0),
(123, 'Narrative Sequencing', 'Ordering events logically in time.', 5.0),
(124, 'Personal Belonging Care', 'Maintaining entropy of personal items.', 6.0),
(125, 'Self-Toileting & Hygiene', 'Basic biological autonomy and cleanliness.', 3.0),

-- Technical & Financial Circuits
(126, 'Digital Literacy', 'Interfacing with abstract digital systems.', 7.0),
(127, 'Keyboarding/Input Mastery', 'Manipulating digital inputs with fine motor speed.', 8.0),
(128, 'Information Searching', 'Retrieving specific data from a wider corpus.', 10.0),
(129, 'Basic Money Handling', 'Manipulating physical currency and simple value exchange.', 7.0),
(130, 'Resource Estimation', 'Allocating supply for a task before commencement.', 11.0),
(131, 'Number Correspondence', 'Counting units and mapping them to values.', 5.0);

-- 2. HumanOS Sensitivity Spectrums (The Calibration Gauges)
-- Note: These might live in a separate 'sensitivities_key' table or as special Skill IDs
CREATE TABLE IF NOT EXISTS sensitivities_key (
    SensitivityID INTEGER PRIMARY KEY AUTOINCREMENT,
    Sensitivity_name TEXT NOT NULL,
    Description TEXT,
    HumanOS_Logic TEXT
);

INSERT OR REPLACE INTO sensitivities_key (SensitivityID, Sensitivity_name, Description, HumanOS_Logic) VALUES
(1, 'Shy — Outgoing', 'Sensitivity to social visibility.', 'Social Gravity Sensitivity'),
(2, 'Guilt-Aware — Remorseless', 'Strength of the internal social governor.', 'Response to Guilt'),
(3, 'Private — Expressive', 'Degree to which internal states are externalized.', 'Emotional Transparency'),
(4, 'Aggressive — Passive', 'Default direction of emotional energy in response to friction.', 'Energy Directionality'),
(5, 'Empathic — Detached', 'Resonance with the emotional states of others.', 'Mirror Neuron Loop'),
(6, 'Generous — Self-Interested', 'Default mode for resource allocation.', 'Resource Allocation Gate'),
(7, 'Bored — Enthusiastic', 'Baseline arousal required for engagement.', 'Voltage Sensitivity'),
(8, 'Patient — Impatient', 'Time-constant between stimulus and response.', 'Impulse/Reflection Gap'),
(9, 'Fragile — Resilient', 'Systems self-righting speed after stress.', 'Response to Failure'),
(10, 'Risk-Averse — Risk-Seeking', 'Threshold for entering uncertain states.', 'Academic Risk Tolerance'),
(11, 'Independent — Dependent', 'Reliance on external instruction vs internal drive.', 'Agency Threshold'),
(12, 'Conformist — Rebellious', 'Response to the social pressure of rules.', 'Response to Authority'),
(13, 'Literal — Adaptable', 'Comfort with gray areas and lack of definition.', 'Response to Ambiguity'),
(14, 'Modest — Status-Driven', 'Sensitivity to external rank and validation.', 'Response to Praise/Privilege'),
(15, 'Truthful — Strategic', 'Default strategy for navigating risk and accountability.', 'Integrate Logic');

-- 3. CHISG Semantic Links (The Integrity Layer)
INSERT OR REPLACE INTO semantic_links (id, source, relation, target, hierarchy_level) VALUES
('sl_1', 'Gross Motor Skills', 'underpins', 'Fine Motor Skills', 1),
('sl_2', 'Phonemic Awareness', 'precursor to', 'Phonics', 1),
('sl_3', 'Instruction Retention (2-Step)', 'precursor to', 'Attention', 1),
('sl_4', 'Number Correspondence', 'precursor to', 'Basic Money Handling', 1),
('sl_5', 'Digital Literacy', 'requires', 'Symbol Recognition', 2),
('sl_6', 'Self-regulation', 'modulates', 'Impulsivity', 0),
('sl_7', 'Resilience', 'mitigates', 'Pain Tolerance', 0);
