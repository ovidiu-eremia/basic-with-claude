package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

func TestInterpreter_LenFunction(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		wantErr  bool
	}{
		{
			name:     "empty string",
			args:     []types.Value{types.NewStringValue("")},
			expected: types.NewNumberValue(0),
			wantErr:  false,
		},
		{
			name:     "single character",
			args:     []types.Value{types.NewStringValue("A")},
			expected: types.NewNumberValue(1),
			wantErr:  false,
		},
		{
			name:     "hello world",
			args:     []types.Value{types.NewStringValue("HELLO WORLD")},
			expected: types.NewNumberValue(11),
			wantErr:  false,
		},
		{
			name:     "string with spaces",
			args:     []types.Value{types.NewStringValue("  HELLO  ")},
			expected: types.NewNumberValue(9),
			wantErr:  false,
		},
		{
			name:    "wrong number of arguments - zero",
			args:    []types.Value{},
			wantErr: true,
		},
		{
			name:    "wrong number of arguments - two",
			args:    []types.Value{types.NewStringValue("TEST"), types.NewNumberValue(1)},
			wantErr: true,
		},
		{
			name:    "wrong argument type - number",
			args:    []types.Value{types.NewNumberValue(42)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateLenFunction(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_LeftFunction(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		wantErr  bool
	}{
		{
			name: "extract first 5 characters",
			args: []types.Value{
				types.NewStringValue("HELLO WORLD"),
				types.NewNumberValue(5),
			},
			expected: types.NewStringValue("HELLO"),
			wantErr:  false,
		},
		{
			name: "extract more characters than available",
			args: []types.Value{
				types.NewStringValue("HI"),
				types.NewNumberValue(10),
			},
			expected: types.NewStringValue("HI"),
			wantErr:  false,
		},
		{
			name: "extract zero characters",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewNumberValue(0),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract negative characters",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewNumberValue(-5),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract from empty string",
			args: []types.Value{
				types.NewStringValue(""),
				types.NewNumberValue(3),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract exact length",
			args: []types.Value{
				types.NewStringValue("TEST"),
				types.NewNumberValue(4),
			},
			expected: types.NewStringValue("TEST"),
			wantErr:  false,
		},
		{
			name:    "wrong number of arguments - one",
			args:    []types.Value{types.NewStringValue("TEST")},
			wantErr: true,
		},
		{
			name:    "wrong number of arguments - three",
			args:    []types.Value{types.NewStringValue("TEST"), types.NewNumberValue(1), types.NewNumberValue(2)},
			wantErr: true,
		},
		{
			name: "wrong first argument type",
			args: []types.Value{
				types.NewNumberValue(123),
				types.NewNumberValue(2),
			},
			wantErr: true,
		},
		{
			name: "wrong second argument type",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewStringValue("2"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateLeftFunction(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_RightFunction(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		wantErr  bool
	}{
		{
			name: "extract last 5 characters",
			args: []types.Value{
				types.NewStringValue("HELLO WORLD"),
				types.NewNumberValue(5),
			},
			expected: types.NewStringValue("WORLD"),
			wantErr:  false,
		},
		{
			name: "extract more characters than available",
			args: []types.Value{
				types.NewStringValue("HI"),
				types.NewNumberValue(10),
			},
			expected: types.NewStringValue("HI"),
			wantErr:  false,
		},
		{
			name: "extract zero characters",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewNumberValue(0),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract negative characters",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewNumberValue(-5),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract from empty string",
			args: []types.Value{
				types.NewStringValue(""),
				types.NewNumberValue(3),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "extract exact length",
			args: []types.Value{
				types.NewStringValue("TEST"),
				types.NewNumberValue(4),
			},
			expected: types.NewStringValue("TEST"),
			wantErr:  false,
		},
		{
			name: "extract single character from end",
			args: []types.Value{
				types.NewStringValue("ABCDE"),
				types.NewNumberValue(1),
			},
			expected: types.NewStringValue("E"),
			wantErr:  false,
		},
		{
			name:    "wrong number of arguments - one",
			args:    []types.Value{types.NewStringValue("TEST")},
			wantErr: true,
		},
		{
			name:    "wrong number of arguments - three",
			args:    []types.Value{types.NewStringValue("TEST"), types.NewNumberValue(1), types.NewNumberValue(2)},
			wantErr: true,
		},
		{
			name: "wrong first argument type",
			args: []types.Value{
				types.NewNumberValue(123),
				types.NewNumberValue(2),
			},
			wantErr: true,
		},
		{
			name: "wrong second argument type",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewStringValue("2"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateRightFunction(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
