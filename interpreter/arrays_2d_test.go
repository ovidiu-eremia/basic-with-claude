package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

func TestInterpreter_DeclareArray2D_AndGetSet(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	// Declare 2D numeric array S(2,3)
	err := interp.DeclareArray("S", []int{2, 3}, false)
	require.NoError(t, err)

	// Default should be 0
	v, err := interp.GetArrayElement("S", []int{1, 2})
	require.NoError(t, err)
	assert.Equal(t, types.NumberType, v.Type)
	assert.Equal(t, 0.0, v.Number)

	// Set and get
	err = interp.SetArrayElement("S", []int{1, 2}, types.NewNumberValue(7))
	require.NoError(t, err)
	v, err = interp.GetArrayElement("S", []int{1, 2})
	require.NoError(t, err)
	assert.Equal(t, 7.0, v.Number)

	// Bounds errors
	_, err = interp.GetArrayElement("S", []int{2, 4})
	assert.Error(t, err)
}

func TestInterpreter_DeclareArray2D_String(t *testing.T) {
	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)

	err := interp.DeclareArray("N$", []int{1, 1}, true)
	require.NoError(t, err)

	// Set string value
	err = interp.SetArrayElement("N$", []int{0, 1}, types.NewStringValue("HI"))
	require.NoError(t, err)
	v, err := interp.GetArrayElement("N$", []int{0, 1})
	require.NoError(t, err)
	assert.Equal(t, types.StringType, v.Type)
	assert.Equal(t, "HI", v.String)
}
