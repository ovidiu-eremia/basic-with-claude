package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

func TestInterpreter_TabFunction(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	t.Run("basic spaces", func(t *testing.T) {
		got, err := interp.evaluateTabFunction([]types.Value{types.NewNumberValue(5)})
		require.NoError(t, err)
		assert.Equal(t, types.NewStringValue("     "), got)
	})

	t.Run("zero and negative yield empty or zero length", func(t *testing.T) {
		got0, err := interp.evaluateTabFunction([]types.Value{types.NewNumberValue(0)})
		require.NoError(t, err)
		assert.Equal(t, types.NewStringValue(""), got0)

		gotNeg, err := interp.evaluateTabFunction([]types.Value{types.NewNumberValue(-3)})
		require.NoError(t, err)
		assert.Equal(t, types.NewStringValue(""), gotNeg)
	})

	t.Run("non-integer counts are floored", func(t *testing.T) {
		got, err := interp.evaluateTabFunction([]types.Value{types.NewNumberValue(3.8)})
		require.NoError(t, err)
		assert.Equal(t, types.NewStringValue("   "), got)
	})

	t.Run("type mismatch", func(t *testing.T) {
		_, err := interp.evaluateTabFunction([]types.Value{types.NewStringValue("X")})
		assert.Error(t, err)
	})
}
