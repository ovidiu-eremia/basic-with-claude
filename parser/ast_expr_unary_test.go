package parser

import (
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestUnaryOperation_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		operand  types.Value
		expected float64
	}{
		{"negation", "-", types.NewNumberValue(5), -5},
		{"unary plus", "+", types.NewNumberValue(5), 5},
		{"negate negative", "-", types.NewNumberValue(-3), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("OPERAND", tt.operand)

			expr := &UnaryOperation{
				Operator: tt.operator,
				Right:    &VariableReference{Name: "OPERAND", Line: 1},
				Line:     1,
			}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestUnaryOperation_Evaluate_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		operand  types.Value
		errMsg   string
	}{
		{"negate string", "-", types.NewStringValue("HELLO"), "cannot negate non-numeric value"},
		{"unary plus on string", "+", types.NewStringValue("HELLO"), "cannot apply unary plus to non-numeric value"},
		{"unknown operator", "!", types.NewNumberValue(5), "unknown unary operator: !"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("OPERAND", tt.operand)

			expr := &UnaryOperation{
				Operator: tt.operator,
				Right:    &VariableReference{Name: "OPERAND", Line: 1},
				Line:     1,
			}

			_, err := expr.Evaluate(mock)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
