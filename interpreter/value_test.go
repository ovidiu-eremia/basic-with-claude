package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValue_Creation(t *testing.T) {
	t.Run("numeric value", func(t *testing.T) {
		v := NewNumberValue(42.5)
		assert.True(t, v.IsNumber())
		assert.False(t, v.IsString())
		assert.Equal(t, NumberType, v.Type)
		assert.Equal(t, 42.5, v.Number)
	})

	t.Run("string value", func(t *testing.T) {
		v := NewStringValue("hello")
		assert.True(t, v.IsString())
		assert.False(t, v.IsNumber())
		assert.Equal(t, StringType, v.Type)
		assert.Equal(t, "hello", v.String)
	})
}

func TestValue_ParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Value
	}{
		{"integer", "42", NewNumberValue(42)},
		{"float", "42.5", NewNumberValue(42.5)},
		{"negative", "-123", NewNumberValue(-123)},
		{"string", "hello", NewStringValue("hello")},
		{"mixed", "42abc", NewStringValue("42abc")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseValue(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValue_ToString(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"whole number", NewNumberValue(42), "42"},
		{"float", NewNumberValue(42.5), "42.5"},
		{"zero", NewNumberValue(0), "0"},
		{"negative", NewNumberValue(-123), "-123"},
		{"string", NewStringValue("hello"), "hello"},
		{"empty string", NewStringValue(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.ToString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValue_ToNumber(t *testing.T) {
	tests := []struct {
		name        string
		value       Value
		expected    float64
		expectError bool
	}{
		{"number value", NewNumberValue(42.5), 42.5, false},
		{"string number", NewStringValue("123"), 123, false},
		{"string float", NewStringValue("42.5"), 42.5, false},
		{"invalid string", NewStringValue("abc"), 0, true},
		{"mixed string", NewStringValue("42abc"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.value.ToNumber()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValue_ArithmeticOperations(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		v1 := NewNumberValue(5)
		v2 := NewNumberValue(3)
		result, err := v1.Add(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(8), result)
	})

	t.Run("subtraction", func(t *testing.T) {
		v1 := NewNumberValue(5)
		v2 := NewNumberValue(3)
		result, err := v1.Subtract(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(2), result)
	})

	t.Run("multiplication", func(t *testing.T) {
		v1 := NewNumberValue(5)
		v2 := NewNumberValue(3)
		result, err := v1.Multiply(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(15), result)
	})

	t.Run("division", func(t *testing.T) {
		v1 := NewNumberValue(10)
		v2 := NewNumberValue(2)
		result, err := v1.Divide(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(5), result)
	})

	t.Run("division by zero", func(t *testing.T) {
		v1 := NewNumberValue(10)
		v2 := NewNumberValue(0)
		_, err := v1.Divide(v2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "division by zero")
	})

	t.Run("power", func(t *testing.T) {
		v1 := NewNumberValue(2)
		v2 := NewNumberValue(3)
		result, err := v1.Power(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(8), result)
	})

	t.Run("string operands", func(t *testing.T) {
		v1 := NewStringValue("5")
		v2 := NewStringValue("3")
		result, err := v1.Add(v2)
		require.NoError(t, err)
		assert.Equal(t, NewNumberValue(8), result)
	})

	t.Run("invalid string operand", func(t *testing.T) {
		v1 := NewStringValue("abc")
		v2 := NewNumberValue(3)
		_, err := v1.Add(v2)
		assert.Error(t, err)
	})
}
