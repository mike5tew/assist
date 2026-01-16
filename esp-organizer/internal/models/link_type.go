package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LinkType represents a reusable relationship type between entities.
// Stores both directions of the relationship for symmetric access.
type LinkType struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// The link label from A's perspective (e.g., "causes reduction in")
	LinkAToB string `bson:"link_a_to_b" json:"link_a_to_b"`

	// The link label from B's perspective (e.g., "is reduced by")
	LinkBToA string `bson:"link_b_to_a" json:"link_b_to_a"`

	// Optional category for grouping (e.g., "causal", "hierarchical", "temporal")
	Category string `bson:"category,omitempty" json:"category,omitempty"`

	// Usage count for sorting by popularity
	UsageCount int `bson:"usage_count" json:"usage_count"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// CommonLinkTypes provides seed data for common relationship types
var CommonLinkTypes = []LinkType{
	{LinkAToB: "causes", LinkBToA: "is caused by", Category: "causal"},
	{LinkAToB: "causes reduction in", LinkBToA: "is reduced by", Category: "causal"},
	{LinkAToB: "causes increase in", LinkBToA: "is increased by", Category: "causal"},
	{LinkAToB: "inhibits", LinkBToA: "is inhibited by", Category: "causal"},
	{LinkAToB: "activates", LinkBToA: "is activated by", Category: "causal"},
	{LinkAToB: "produces", LinkBToA: "is produced by", Category: "causal"},
	{LinkAToB: "requires", LinkBToA: "is required by", Category: "dependency"},
	{LinkAToB: "depends on", LinkBToA: "is depended on by", Category: "dependency"},
	{LinkAToB: "is a type of", LinkBToA: "includes", Category: "hierarchical"},
	{LinkAToB: "is part of", LinkBToA: "contains", Category: "hierarchical"},
	{LinkAToB: "is located in", LinkBToA: "contains", Category: "spatial"},
	{LinkAToB: "precedes", LinkBToA: "follows", Category: "temporal"},
	{LinkAToB: "occurs before", LinkBToA: "occurs after", Category: "temporal"},
	{LinkAToB: "is associated with", LinkBToA: "is associated with", Category: "association"},
	{LinkAToB: "is similar to", LinkBToA: "is similar to", Category: "comparison"},
	{LinkAToB: "is different from", LinkBToA: "is different from", Category: "comparison"},
	{LinkAToB: "leads to", LinkBToA: "results from", Category: "causal"},
	{LinkAToB: "prevents", LinkBToA: "is prevented by", Category: "causal"},
	{LinkAToB: "treats", LinkBToA: "is treated by", Category: "medical"},
	{LinkAToB: "presents with", LinkBToA: "is presentation of", Category: "medical"},
}
