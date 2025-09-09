package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableReference_Evaluate_Error(t *testing.T) {
	mock := newMockOps()
	mock.getVariableError = errors.New("variable error")

	expr := &VariableReference{Name: "A"}

	_, err := expr.Evaluate(mock)

	assert.Error(t, err)
}
