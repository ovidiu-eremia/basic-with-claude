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

func TestInterpreter_MidFunction(t *testing.T) {
	tests := []struct {
		name     string
		args     []types.Value
		expected types.Value
		wantErr  bool
	}{
		{
			name: "basic substring",
			args: []types.Value{
				types.NewStringValue("HELLO"),
				types.NewNumberValue(2),
				types.NewNumberValue(3),
			},
			expected: types.NewStringValue("ELL"),
			wantErr:  false,
		},
		{
			name: "start beyond length",
			args: []types.Value{
				types.NewStringValue("ABC"),
				types.NewNumberValue(5),
				types.NewNumberValue(2),
			},
			expected: types.NewStringValue(""),
			wantErr:  false,
		},
		{
			name: "length overflow to end",
			args: []types.Value{
				types.NewStringValue("ABCDE"),
				types.NewNumberValue(4),
				types.NewNumberValue(99),
			},
			expected: types.NewStringValue("DE"),
			wantErr:  false,
		},
		{ // wrong arity
			name:    "wrong number of arguments",
			args:    []types.Value{types.NewStringValue("A"), types.NewNumberValue(1)},
			wantErr: true,
		},
		{ // wrong types
			name: "wrong types",
			args: []types.Value{
				types.NewNumberValue(123),
				types.NewNumberValue(1),
				types.NewNumberValue(1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateMidFunction(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_ChrFunction(t *testing.T) {
	tests := []struct {
		name     string
		arg      types.Value
		expected types.Value
		wantErr  bool
	}{
		{name: "A", arg: types.NewNumberValue(65), expected: types.NewStringValue("A")},
		{name: "a", arg: types.NewNumberValue(97), expected: types.NewStringValue("a")},
		{name: "zero", arg: types.NewNumberValue(0), expected: types.NewStringValue("\x00")},
		{name: "wrong type", arg: types.NewStringValue("A"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateChrFunction([]types.Value{tt.arg})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_AscFunction(t *testing.T) {
	tests := []struct {
		name     string
		arg      types.Value
		expected types.Value
		wantErr  bool
	}{
		{name: "A", arg: types.NewStringValue("A"), expected: types.NewNumberValue(65)},
		{name: "hello", arg: types.NewStringValue("hello"), expected: types.NewNumberValue(104)},
		{name: "empty", arg: types.NewStringValue(""), expected: types.NewNumberValue(0)},
		{name: "wrong type", arg: types.NewNumberValue(65), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateAscFunction([]types.Value{tt.arg})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_StrFunction(t *testing.T) {
	tests := []struct {
		name     string
		arg      types.Value
		expected types.Value
		wantErr  bool
	}{
		{name: "int", arg: types.NewNumberValue(42), expected: types.NewStringValue("42")},
		{name: "neg", arg: types.NewNumberValue(-3), expected: types.NewStringValue("-3")},
		{name: "float", arg: types.NewNumberValue(12.5), expected: types.NewStringValue("12.5")},
		{name: "wrong type", arg: types.NewStringValue("A"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateStrFunction([]types.Value{tt.arg})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_ValFunction(t *testing.T) {
	tests := []struct {
		name     string
		arg      types.Value
		expected types.Value
		wantErr  bool
	}{
		{name: "int", arg: types.NewStringValue("123"), expected: types.NewNumberValue(123)},
		{name: "float", arg: types.NewStringValue("12.5"), expected: types.NewNumberValue(12.5)},
		{name: "leading spaces", arg: types.NewStringValue("  12.5"), expected: types.NewNumberValue(12.5)},
		{name: "empty", arg: types.NewStringValue(""), expected: types.NewNumberValue(0)},
		{name: "nonnumeric", arg: types.NewStringValue("A"), expected: types.NewNumberValue(0)},
		{name: "wrong type", arg: types.NewNumberValue(1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := runtime.NewTestRuntime()
			interp := NewInterpreter(rt)

			result, err := interp.evaluateValFunction([]types.Value{tt.arg})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpreter_AbsIntSqrFunctions(t *testing.T) {
	t.Run("ABS", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct {
			in  float64
			out float64
		}{
			{5, 5}, {-3, 3}, {0, 0}, {-2.5, 2.5},
		}
		for _, c := range cases {
			got, err := interp.evaluateAbsFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		// type mismatch
		_, err := interp.evaluateAbsFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})

	t.Run("INT", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct {
			in  float64
			out float64
		}{
			{3.7, 3}, {-3.2, -4}, {0, 0}, {-0.1, -1},
		}
		for _, c := range cases {
			got, err := interp.evaluateIntFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		// type mismatch
		_, err := interp.evaluateIntFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})

	t.Run("SQR", func(t *testing.T) {
		rt := runtime.NewTestRuntime()
		interp := NewInterpreter(rt)
		cases := []struct {
			in  float64
			out float64
		}{
			{9, 3}, {2.25, 1.5}, {0, 0},
		}
		for _, c := range cases {
			got, err := interp.evaluateSqrFunction([]types.Value{types.NewNumberValue(c.in)})
			require.NoError(t, err)
			assert.Equal(t, types.NewNumberValue(c.out), got)
		}
		// negative -> illegal quantity
		_, err := interp.evaluateSqrFunction([]types.Value{types.NewNumberValue(-1)})
		assert.Error(t, err)
		// type mismatch
		_, err = interp.evaluateSqrFunction([]types.Value{types.NewStringValue("A")})
		assert.Error(t, err)
	})
}
