package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

func TestInterpreter_TrigFunctions(t *testing.T) {
	t.Run("SIN", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		// Exact values for selected inputs
		cases := []struct{ in, out float64 }{
			{0, 0},
			{1.5707963267948966, 1}, // pi/2
		}
		for _, c := range cases {
			got, err := interp.evaluateSinFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		// arity and type checks
		_, err := interp.evaluateSinFunction([]types.Value{})
		assert.Error(t, err)
		_, err = interp.evaluateSinFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})

	t.Run("COS", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct{ in, out float64 }{
			{0, 1},
		}
		for _, c := range cases {
			got, err := interp.evaluateCosFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		_, err := interp.evaluateCosFunction([]types.Value{})
		assert.Error(t, err)
		_, err = interp.evaluateCosFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})

	t.Run("TAN", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct{ in, out float64 }{
			{0, 0},
		}
		for _, c := range cases {
			got, err := interp.evaluateTanFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		_, err := interp.evaluateTanFunction([]types.Value{})
		assert.Error(t, err)
		_, err = interp.evaluateTanFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})

	t.Run("ATN", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct{ in, out float64 }{
			{0, 0},
			{1, 0.7853981633974483}, // pi/4
		}
		for _, c := range cases {
			got, err := interp.evaluateAtnFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		_, err := interp.evaluateAtnFunction([]types.Value{})
		assert.Error(t, err)
		_, err = interp.evaluateAtnFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})
}
