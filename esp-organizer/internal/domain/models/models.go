package models

import (
	"context"
	"time"
)

// Filter represents a structured query filter for search operations
type Filter struct {
	Operator     string    `json:"operator,omitempty"`     // "And" or "Or"
	Operands     []Operand `json:"operands,omitempty"`     // List of operands in this filter
	NestedFilter *Filter   `json:"nestedFilter,omitempty"` // Support for nested filters
}

// QueryResponse represents the response to a query
type QueryResponse struct {
	Query           string                   `json:"query"`
	Timestamp       time.Time                `json:"timestamp"`
	MongoResults    []Skill                  `json:"mongo_results,omitempty"`
	WeaviateResults []map[string]interface{} `json:"weaviate_results,omitempty"`
	Synthesis       string                   `json:"synthesis,omitempty"`
	RelatedSkills   []Skill                  `json:"related_skills,omitempty"`
}

type QueryProcessor interface {
	ProcessQuery(ctx context.Context, query string, filters *Filter, weaviateClass string) (*QueryResponse, error)
}

type Operand struct {
	Path         []string `json:"path,omitempty"`         // Property path to filter on
	Operator     string   `json:"operator,omitempty"`     // Operator like "Equal", "GreaterThan", etc.
	ValueString  string   `json:"valueString,omitempty"`  // String value for the filter
	ValueNumber  float64  `json:"valueNumber,omitempty"`  // Numeric value for the filter
	ValueBoolean bool     `json:"valueBoolean,omitempty"` // Boolean value for the filter
	NestedFilter *Filter  `json:"nestedFilter,omitempty"` // Support for nested filters in operands
}

// SemanticQueryRequest represents a query to the semantic search API
type SemanticQueryRequest struct {
	Query    string   `json:"query"`
	Filters  *Filter  `json:"filters,omitempty"`
	Subjects []string `json:"subjects,omitempty"` // Multiple subjects to filter by
}
