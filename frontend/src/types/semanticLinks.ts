/**
 * types/semanticLinks.ts
 * 
 * TypeScript interfaces for semantic link extraction feature
 * 
 * Core Model:
 * - Statement: The full, vectorizable fact ("XLA causes reduction in Ig")
 * - Terms: Source and target concepts for graph traversal
 * - Relations: Bidirectional (forward + inverse)
 * - Conditions: Circumstances under which the relationship is valid
 */

// Common relationship types with their inverses
export const RELATIONSHIP_PAIRS: Record<string, string> = {
  'causes': 'is_caused_by',
  'causes_reduction_in': 'is_reduced_by',
  'causes_increase_in': 'is_increased_by',
  'treats': 'is_treated_by',
  'prevents': 'is_prevented_by',
  'is_a': 'has_subtype',
  'part_of': 'contains',
  'requires': 'is_required_by',
  'enables': 'is_enabled_by',
  'results_in': 'results_from',
  'contrasts_with': 'contrasts_with', // Symmetric
  'similar_to': 'similar_to', // Symmetric
  'prerequisite_for': 'has_prerequisite',
  'associated_with': 'associated_with', // Symmetric
};

// Condition relation types (many are symmetric)
export const CONDITION_RELATIONS: Record<string, string> = {
  'occurs_when': 'occurs_when', // Symmetric
  'occurs_at': 'occurs_at', // Symmetric - temporal
  'occurs_in': 'occurs_in', // Symmetric - spatial/context
  'requires': 'is_required_by',
  'except_when': 'except_when', // Symmetric
  'only_if': 'only_if', // Symmetric
};

export type RelationshipType = keyof typeof RELATIONSHIP_PAIRS | string;
export type ConditionRelationType = keyof typeof CONDITION_RELATIONS | string;

export type QualityFlag = 'high_confidence' | 'needs_verification' | 'problematic';

// Condition attached to a semantic link
export interface LinkCondition {
  term: string;
  forwardRelation: ConditionRelationType;
  inverseRelation: ConditionRelationType;
  required?: boolean;
  description?: string;
}

export interface SemanticLinkSelection {
  statement: string; // The full fact to be vectorized
  sourceTerm: string;
  targetTerm: string;
  forwardRelation: RelationshipType;
  inverseRelation: RelationshipType;
  conditions?: LinkCondition[];
  confidence?: number;
  qualityFlag?: QualityFlag;
  position?: {
    start: number;
    end: number;
  };
}

export interface ExtractionResult {
  inserted: number;
  errors: string[];
}

export interface SemanticLink {
  id: string;
  statement: string;
  sourceTerm: string;
  targetTerm: string;
  forwardRelation: RelationshipType;
  inverseRelation: RelationshipType;
  conditions?: LinkCondition[];
  confidence: number;
  domain: string;
  createdAt: string;
  metadata?: Record<string, any>;
  qualityFlag?: QualityFlag;
}

export interface RAGContext {
  semanticContext: string[];
  summaryContext: string[];
  documentContext: string[];
}

export interface SearchResult {
  semanticContext: string[];
  summaryContext: string[];
  documents: string[];
  answer: string;
}
