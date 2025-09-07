package parser

import (
	"errors"
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestLetStatement_Execute(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		expression   Expression
		expectedType types.ValueType
		expectedVal  interface{}
	}{
		{
			name:         "assign number to variable",
			variable:     "A",
			expression:   &NumberLiteral{Value: "42", Line: 1},
			expectedType: types.NumberType,
			expectedVal:  42.0,
		},
		{
			name:         "assign string to variable",
			variable:     "A$",
			expression:   &StringLiteral{Value: "HELLO", Line: 1},
			expectedType: types.StringType,
			expectedVal:  "HELLO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			stmt := &LetStatement{
				Variable:   tt.variable,
				Expression: tt.expression,
				Line:       1,
			}

			err := stmt.Execute(mock)
			assert.NoError(t, err)

			value, exists := mock.variables[tt.variable]
			assert.True(t, exists)
			assert.Equal(t, tt.expectedType, value.Type)

			if tt.expectedType == types.NumberType {
				assert.Equal(t, tt.expectedVal, value.Number)
			} else {
				assert.Equal(t, tt.expectedVal, value.String)
			}
		})
	}
}

func TestLetStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("expression evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		stmt := &LetStatement{
			Variable:   "A",
			Expression: &VariableReference{Name: "B", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("set variable error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariableError = errors.New("set error")

		stmt := &LetStatement{
			Variable:   "A",
			Expression: &NumberLiteral{Value: "42", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}
