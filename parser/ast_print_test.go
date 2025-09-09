package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintStatement_Execute(t *testing.T) {
	tests := []struct {
		name           string
		expression     Expression
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "print string literal",
			expression:     &StringLiteral{Value: "HELLO", BaseNode: BaseNode{Line: 1}},
			expectedOutput: "HELLO",
			expectError:    false,
		},
		{
			name:           "print number literal",
			expression:     &NumberLiteral{Value: "42", BaseNode: BaseNode{Line: 1}},
			expectedOutput: "42",
			expectError:    false,
		},
		{
			name:           "print empty string",
			expression:     &StringLiteral{Value: "", BaseNode: BaseNode{Line: 1}},
			expectedOutput: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			stmt := &PrintStatement{BaseNode: BaseNode{Line: 1}, Expression: tt.expression}

			err := stmt.Execute(mock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.Len(t, mock.getOutput(), 1)
				assert.Equal(t, tt.expectedOutput, mock.getOutput()[0])
			}
		})
	}
}

func TestPrintStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("expression evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		stmt := &PrintStatement{
			Expression: &VariableReference{Name: "A", BaseNode: BaseNode{Line: 1}},
			BaseNode:   BaseNode{Line: 1},
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("print line error", func(t *testing.T) {
		mock := newMockOps()
		mock.printLineError = errors.New("print error")

		stmt := &PrintStatement{
			Expression: &StringLiteral{Value: "TEST", BaseNode: BaseNode{Line: 1}},
			BaseNode:   BaseNode{Line: 1},
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}
