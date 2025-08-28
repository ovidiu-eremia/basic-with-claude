// ABOUTME: Value type system for BASIC variables and expressions
// ABOUTME: Handles numeric and string values with proper type distinctions and conversions

package interpreter

import (
	"fmt"
	"math"
	"strconv"
)

// ValueType represents the type of a BASIC value
type ValueType int

const (
	NumberType ValueType = iota
	StringType
)

// Value represents a BASIC value with type information
type Value struct {
	Type   ValueType
	Number float64
	String string
}

// NewNumberValue creates a numeric value
func NewNumberValue(n float64) Value {
	return Value{Type: NumberType, Number: n}
}

// NewStringValue creates a string value
func NewStringValue(s string) Value {
	return Value{Type: StringType, String: s}
}

// ParseValue creates a Value from a string representation
func ParseValue(s string) (Value, error) {
	// Try to parse as number first
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return NewNumberValue(num), nil
	}
	// Otherwise treat as string
	return NewStringValue(s), nil
}

// ToString converts the value to its string representation
func (v Value) ToString() string {
	switch v.Type {
	case NumberType:
		// Format numbers as integers if they are whole numbers
		if v.Number == float64(int64(v.Number)) {
			return strconv.FormatInt(int64(v.Number), 10)
		}
		return strconv.FormatFloat(v.Number, 'g', -1, 64)
	case StringType:
		return v.String
	default:
		return ""
	}
}

// ToNumber converts the value to a numeric value
func (v Value) ToNumber() (float64, error) {
	switch v.Type {
	case NumberType:
		return v.Number, nil
	case StringType:
		num, err := strconv.ParseFloat(v.String, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to number", v.String)
		}
		return num, nil
	default:
		return 0, fmt.Errorf("invalid value type")
	}
}

// IsNumber returns true if the value is numeric
func (v Value) IsNumber() bool {
	return v.Type == NumberType
}

// IsString returns true if the value is a string
func (v Value) IsString() bool {
	return v.Type == StringType
}

// Add performs addition on two values
func (v Value) Add(other Value) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(left + right), nil
}

// Subtract performs subtraction on two values
func (v Value) Subtract(other Value) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(left - right), nil
}

// Multiply performs multiplication on two values
func (v Value) Multiply(other Value) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(left * right), nil
}

// Divide performs division on two values
func (v Value) Divide(other Value) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	if right == 0 {
		return Value{}, fmt.Errorf("division by zero")
	}
	return NewNumberValue(left / right), nil
}

// Power performs exponentiation on two values
func (v Value) Power(other Value) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(math.Pow(left, right)), nil
}