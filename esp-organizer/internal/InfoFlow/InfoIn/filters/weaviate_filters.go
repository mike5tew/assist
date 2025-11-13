package filters

import (
	"esp-organizer/internal/models"
	"log"
	"strconv"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
)

// Convert string operator to Weaviate WhereOperator
func getWhereOperator(op string) filters.WhereOperator {
	switch op {
	case "Equal":
		return filters.Equal
	case "NotEqual":
		return filters.NotEqual
	case "Like":
		return filters.Like
	case "ContainsAny":
		return filters.ContainsAny
	case "ContainsAll":
		return filters.ContainsAll
	case "GreaterThan":
		return filters.GreaterThan
	case "GreaterThanEqual":
		return filters.GreaterThanEqual
	case "LessThan":
		return filters.LessThan
	case "LessThanEqual":
		return filters.LessThanEqual
	default:
		log.Printf("Warning: Unrecognized operator %s, defaulting to Equal", op)
		return filters.Equal
	}
}

// BuildWhereFilterFromAPI creates a Weaviate where filter from the models.Filter struct
func BuildWhereFilterFromAPI(filter *models.Filter) *filters.WhereBuilder {
	if filter == nil || (len(filter.Operands) == 0 && filter.NestedFilter == nil) {
		return nil
	}

	// If we have a nested filter but no operands, process just the nested filter
	if len(filter.Operands) == 0 && filter.NestedFilter != nil {
		return BuildWhereFilterFromAPI(filter.NestedFilter)
	}

	// Build each operand into a WhereBuilder
	var operandBuilders []*filters.WhereBuilder
	for _, op := range filter.Operands {
		// Handle nested filter case
		if op.NestedFilter != nil {
			nestedBuilder := BuildWhereFilterFromAPI(op.NestedFilter)
			if nestedBuilder != nil {
				operandBuilders = append(operandBuilders, nestedBuilder)
			}
			continue
		}

		builder := filters.Where().WithPath(op.Path)
		whereOp := getWhereOperator(op.Operator)

		// Handle numeric values for relevant operators
		switch op.Operator {
		case "Equal", "NotEqual", "Like", "ContainsAny", "ContainsAll":
			// These operators typically work with string values
			builder = builder.WithOperator(whereOp).WithValueString(op.ValueString)

		case "GreaterThan", "GreaterThanEqual", "LessThan", "LessThanEqual":
			// Skip empty values for numeric operators
			if op.ValueString == "" {
				log.Printf("Warning: Skipping numeric operator %s with empty value", op.Operator)
				continue // Skip this operand entirely
			}
			// Try to convert the value to a number if possible
			if num, err := strconv.ParseInt(op.ValueString, 10, 64); err == nil {
				// It's an integer
				builder = builder.WithOperator(whereOp).WithValueInt(num)
			} else if num, err := strconv.ParseFloat(op.ValueString, 64); err == nil {
				// It's a float
				builder = builder.WithOperator(whereOp).WithValueNumber(num)
			} else {
				// Fall back to string for non-numeric values
				log.Printf("Warning: Using string value for numeric operator %s: %s", op.Operator, op.ValueString)
				builder = builder.WithOperator(whereOp).WithValueString(op.ValueString)
			}

		default:
			// Default to string for unknown operators
			builder = builder.WithOperator(filters.Equal).WithValueString(op.ValueString)
		}

		operandBuilders = append(operandBuilders, builder)
	}

	// Combine the operands with the top-level operator
	var topOperator filters.WhereOperator
	switch filter.Operator {
	case "And":
		topOperator = filters.And
	case "Or":
		topOperator = filters.Or
	default:
		// Default to AND for safety
		topOperator = filters.And
	}

	return filters.Where().WithOperator(topOperator).WithOperands(operandBuilders)
}

// BuildWhereFilter creates a Weaviate where filter from a map
func BuildWhereFilter(filter map[string]interface{}) *filters.WhereBuilder {
	path, _ := filter["path"].([]string)
	opStr, _ := filter["op"].(string)
	value := filter["value"]

	where := filters.Where().WithPath(path)

	// Convert string operator to WhereOperator type
	op := getWhereOperator(opStr)

	switch opStr {
	case "Equal":
		switch v := value.(type) {
		case string:
			return where.WithOperator(op).WithValueString(v)
		case int:
			return where.WithOperator(op).WithValueInt(int64(v))
		case float64:
			return where.WithOperator(op).WithValueNumber(v)
		case bool:
			return where.WithOperator(op).WithValueBoolean(v)
		}

	case "ContainsAny":
		if items, ok := value.([]string); ok {
			return where.WithOperator(op).WithValueText(items...)
		}

	case "GreaterThan":
		switch v := value.(type) {
		case int:
			return where.WithOperator(op).WithValueNumber(float64(v))
		case float64:
			return where.WithOperator(op).WithValueNumber(v)
		}
	}
	return nil
}

// BuildSubjectsFilter creates a Weaviate where filter for multiple subjects
func BuildSubjectsFilter(subjects []string) *filters.WhereBuilder {
	if len(subjects) == 0 {
		return nil
	}

	// For a single subject
	if len(subjects) == 1 {
		return filters.Where().
			WithPath([]string{"subjects"}).
			WithOperator(filters.ContainsAny).
			WithValueString(subjects[0])
	}

	// For multiple subjects, create OR operands
	var operands []*filters.WhereBuilder
	for _, subject := range subjects {
		operand := filters.Where().
			WithPath([]string{"subjects"}).
			WithOperator(filters.ContainsAny).
			WithValueString(subject)
		operands = append(operands, operand)
	}

	return filters.Where().
		WithOperator(filters.Or).
		WithOperands(operands)
}
