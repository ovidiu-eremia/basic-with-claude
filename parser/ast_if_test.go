package parser

import (
	"errors"
	"testing"

	"basic-interpreter/types"
	"github.com/stretchr/testify/assert"
)

func TestIfStatement_Execute(t *testing.T) {
	tests := []struct {
		name            string
		conditionValue  types.Value
		expectExecution bool
	}{
		{"true condition (non-zero number)", types.NewNumberValue(1), true},
		{"false condition (zero)", types.NewNumberValue(0), false},
		{"true condition (negative number)", types.NewNumberValue(-5), true},
		{"true condition (non-empty string)", types.NewStringValue("HELLO"), true},
		{"false condition (empty string)", types.NewStringValue(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("CONDITION", tt.conditionValue)

			condition := &VariableReference{Name: "CONDITION", Line: 1}
			thenStmt := &PrintStatement{Expression: &StringLiteral{Value: "EXECUTED", Line: 1}, Line: 1}

			stmt := &IfStatement{Condition: condition, ThenStmt: thenStmt, Line: 1}

			err := stmt.Execute(mock)
			assert.NoError(t, err)

			if tt.expectExecution {
				assert.Len(t, mock.getOutput(), 1)
				assert.Equal(t, "EXECUTED", mock.getOutput()[0])
			} else {
				assert.Len(t, mock.getOutput(), 0)
			}
		})
	}
}

func TestIfStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("condition evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		condition := &VariableReference{Name: "A", Line: 1}
		thenStmt := &PrintStatement{Expression: &StringLiteral{Value: "TEST", Line: 1}, Line: 1}

		stmt := &IfStatement{Condition: condition, ThenStmt: thenStmt, Line: 1}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("then statement execution error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariable("A", types.NewNumberValue(1)) // true condition
		mock.printLineError = errors.New("print error")

		condition := &VariableReference{Name: "A", Line: 1}
		thenStmt := &PrintStatement{Expression: &StringLiteral{Value: "TEST", Line: 1}, Line: 1}

		stmt := &IfStatement{Condition: condition, ThenStmt: thenStmt, Line: 1}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}
