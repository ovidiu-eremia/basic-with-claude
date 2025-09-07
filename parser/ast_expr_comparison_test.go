package parser

import (
	"errors"
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestComparisonExpression_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     types.Value
		operator string
		right    types.Value
		expected float64
	}{
		// Numeric comparisons
		{"equal true", types.NewNumberValue(5), "=", types.NewNumberValue(5), 1},
		{"equal false", types.NewNumberValue(5), "=", types.NewNumberValue(3), 0},
		{"not equal true", types.NewNumberValue(5), "<>", types.NewNumberValue(3), 1},
		{"not equal false", types.NewNumberValue(5), "<>", types.NewNumberValue(5), 0},
		{"less than true", types.NewNumberValue(3), "<", types.NewNumberValue(5), 1},
		{"less than false", types.NewNumberValue(5), "<", types.NewNumberValue(3), 0},
		{"greater than true", types.NewNumberValue(5), ">", types.NewNumberValue(3), 1},
		{"greater than false", types.NewNumberValue(3), ">", types.NewNumberValue(5), 0},
		{"less equal true", types.NewNumberValue(3), "<=", types.NewNumberValue(5), 1},
		{"less equal equal", types.NewNumberValue(5), "<=", types.NewNumberValue(5), 1},
		{"less equal false", types.NewNumberValue(5), "<=", types.NewNumberValue(3), 0},
		{"greater equal true", types.NewNumberValue(5), ">=", types.NewNumberValue(3), 1},
		{"greater equal equal", types.NewNumberValue(5), ">=", types.NewNumberValue(5), 1},
		{"greater equal false", types.NewNumberValue(3), ">=", types.NewNumberValue(5), 0},

		// String comparisons
		{"string equal true", types.NewStringValue("HELLO"), "=", types.NewStringValue("HELLO"), 1},
		{"string equal false", types.NewStringValue("HELLO"), "=", types.NewStringValue("WORLD"), 0},
		{"string not equal true", types.NewStringValue("HELLO"), "<>", types.NewStringValue("WORLD"), 1},
		{"string not equal false", types.NewStringValue("HELLO"), "<>", types.NewStringValue("HELLO"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", tt.left)
			mock.setVariable("RIGHT", tt.right)

			expr := &ComparisonExpression{
				Left:     &VariableReference{Name: "LEFT", Line: 1},
				Operator: tt.operator,
				Right:    &VariableReference{Name: "RIGHT", Line: 1},
				Line:     1,
			}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestComparisonExpression_Evaluate_ErrorCases(t *testing.T) {
	t.Run("left evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		expr := &ComparisonExpression{
			Left:     &VariableReference{Name: "A", Line: 1},
			Operator: "=",
			Right:    &NumberLiteral{Value: "5", Line: 1},
			Line:     1,
		}

		_, err := expr.Evaluate(mock)

		assert.Error(t, err)
	})

	t.Run("comparison error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariable("LEFT", types.NewNumberValue(5))
		mock.setVariable("RIGHT", types.NewStringValue("HELLO"))

		expr := &ComparisonExpression{
			Left:     &VariableReference{Name: "LEFT", Line: 1},
			Operator: "<",
			Right:    &VariableReference{Name: "RIGHT", Line: 1},
			Line:     1,
		}

		_, err := expr.Evaluate(mock)

		assert.Error(t, err)
	})
}
