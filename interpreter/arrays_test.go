package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
)

func TestInterpreter_DeclareArray(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	// Declare numeric array
	err := interp.DeclareArray("A", 5, false)
	require.NoError(t, err)

	// Redimension should error
	err = interp.DeclareArray("A", 6, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "REDIM'D")

	// Negative size illegal
	err = interp.DeclareArray("B", -1, false)
	assert.Error(t, err)
}
