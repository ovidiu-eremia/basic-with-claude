package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

func TestInterpreter_ExpFunction(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	// basic values
	cases := []struct{ in, out float64 }{
		{0, 1},
		{1, 2.718281828459045},
	}
	for _, c := range cases {
		got, err := interp.evaluateExpFunction([]types.Value{types.NewNumberValue(c.in)})
		require.NoError(t, err)
		assert.Equal(t, types.NewNumberValue(c.out), got)
	}

	// arity
	_, err := interp.evaluateExpFunction([]types.Value{})
	assert.Error(t, err)
	// type mismatch
	_, err = interp.evaluateExpFunction([]types.Value{types.NewStringValue("A")})
	assert.Error(t, err)
}

func TestInterpreter_LogFunction(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	// natural log
	cases := []struct{ in, out float64 }{
		{1, 0},
		{2.718281828459045, 1},
	}
	for _, c := range cases {
		got, err := interp.evaluateLogFunction([]types.Value{types.NewNumberValue(c.in)})
		require.NoError(t, err)
		assert.Equal(t, types.NewNumberValue(c.out), got)
	}

	// non-positive -> illegal quantity
	_, err := interp.evaluateLogFunction([]types.Value{types.NewNumberValue(0)})
	assert.Error(t, err)
	// type mismatch
	_, err = interp.evaluateLogFunction([]types.Value{types.NewStringValue("A")})
	assert.Error(t, err)
}

func TestInterpreter_RndFunction(t *testing.T) {
	rt := runtime.NewTestRuntime()
	// Use fixed random source in TestRuntime to make predictable? At least validate range
	interp := NewInterpreter(rt)

	// arity
	_, err := interp.evaluateRndFunction([]types.Value{})
	assert.Error(t, err)

	// type mismatch
	_, err = interp.evaluateRndFunction([]types.Value{types.NewStringValue("A")})
	assert.Error(t, err)

	// basic range check
	v, err := interp.evaluateRndFunction([]types.Value{types.NewNumberValue(1)})
	require.NoError(t, err)
	assert.Equal(t, types.NumberType, v.Type)
	assert.GreaterOrEqual(t, v.Number, 0.0)
	assert.Less(t, v.Number, 1.0)
}
