package parser

import (
	"errors"
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestBinaryOperation_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     types.Value
		operator string
		right    types.Value
		expected float64
	}{
		{"addition", types.NewNumberValue(5), "+", types.NewNumberValue(3), 8},
		{"subtraction", types.NewNumberValue(10), "-", types.NewNumberValue(4), 6},
		{"multiplication", types.NewNumberValue(6), "*", types.NewNumberValue(7), 42},
		{"division", types.NewNumberValue(15), "/", types.NewNumberValue(3), 5},
		{"power", types.NewNumberValue(2), "^", types.NewNumberValue(3), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", tt.left)
			mock.setVariable("RIGHT", tt.right)

			expr := &BinaryOperation{
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

func TestBinaryOperation_Evaluate_StringConcatenation(t *testing.T) {
	mock := newMockOps()
	mock.setVariable("LEFT", types.NewStringValue("HELLO"))
	mock.setVariable("RIGHT", types.NewStringValue("WORLD"))

	expr := &BinaryOperation{
		Left:     &VariableReference{Name: "LEFT", Line: 1},
		Operator: "+",
		Right:    &VariableReference{Name: "RIGHT", Line: 1},
		Line:     1,
	}

	result, err := expr.Evaluate(mock)

	assert.NoError(t, err)
	assert.Equal(t, types.StringType, result.Type)
	assert.Equal(t, "HELLOWORLD", result.String)
}

func TestBinaryOperation_Evaluate_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		operator string
	}{
		{"unknown operator", "**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", types.NewNumberValue(5))
			mock.setVariable("RIGHT", types.NewNumberValue(3))

			expr := &BinaryOperation{
				Left:     &VariableReference{Name: "LEFT", Line: 1},
				Operator: tt.operator,
				Right:    &VariableReference{Name: "RIGHT", Line: 1},
				Line:     1,
			}

			_, err := expr.Evaluate(mock)

			assert.Error(t, err)
		})
	}
}

func TestBinaryOperation_Evaluate_LeftEvaluationError(t *testing.T) {
	mock := newMockOps()
	mock.getVariableError = errors.New("variable error")

	expr := &BinaryOperation{
		Left:     &VariableReference{Name: "A", Line: 1},
		Operator: "+",
		Right:    &NumberLiteral{Value: "3", Line: 1},
		Line:     1,
	}

	_, err := expr.Evaluate(mock)

	assert.Error(t, err)
}
