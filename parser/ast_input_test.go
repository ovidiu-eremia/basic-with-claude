package parser

import (
	"errors"
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestInputStatement_Execute(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		input        string
		expectedType types.ValueType
		expectedVal  interface{}
		expectError  bool
	}{
		{
			name:         "numeric input to numeric variable",
			variable:     "A",
			input:        "42",
			expectedType: types.NumberType,
			expectedVal:  42.0,
			expectError:  false,
		},
		{
			name:         "string input to string variable",
			variable:     "A$",
			input:        "HELLO",
			expectedType: types.StringType,
			expectedVal:  "HELLO",
			expectError:  false,
		},
		{
			name:        "invalid numeric input",
			variable:    "A",
			input:       "ABC",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setInput([]string{tt.input})

			stmt := &InputStatement{Variable: tt.variable, BaseNode: BaseNode{Line: 1}}

			err := stmt.Execute(mock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				value, exists := mock.variables[tt.variable]
				assert.True(t, exists)
				assert.Equal(t, tt.expectedType, value.Type)

				if tt.expectedType == types.NumberType {
					assert.Equal(t, tt.expectedVal, value.Number)
				} else {
					assert.Equal(t, tt.expectedVal, value.String)
				}
			}
		})
	}
}

func TestInputStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("read input error", func(t *testing.T) {
		mock := newMockOps()
		mock.readInputError = errors.New("input error")

		stmt := &InputStatement{Variable: "A", BaseNode: BaseNode{Line: 1}}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("set variable error", func(t *testing.T) {
		mock := newMockOps()
		mock.setInput([]string{"42"})
		mock.setVariableError = errors.New("set error")

		stmt := &InputStatement{Variable: "A", BaseNode: BaseNode{Line: 1}}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}
