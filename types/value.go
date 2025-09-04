// ABOUTME: Value type system for BASIC variables and expressions
// ABOUTME: Handles numeric and string values with proper type distinctions and conversions

package types

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

// binaryArithmeticOp performs a binary arithmetic operation on two values
func (v Value) binaryArithmeticOp(other Value, operation func(float64, float64) float64) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(operation(left, right)), nil
}

// binaryArithmeticOpWithError performs a binary arithmetic operation that can return an error
func (v Value) binaryArithmeticOpWithError(other Value, operation func(float64, float64) (float64, error)) (Value, error) {
	left, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	right, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}
	result, err := operation(left, right)
	if err != nil {
		return Value{}, err
	}
	return NewNumberValue(result), nil
}

// Add performs addition on two values
func (v Value) Add(other Value) (Value, error) {
	// If both values are strings, try numeric conversion first
	if v.Type == StringType && other.Type == StringType {
		// Try to convert both strings to numbers
		leftNum, leftErr := v.ToNumber()
		rightNum, rightErr := other.ToNumber()

		// If both can be converted to numbers, do numeric addition
		if leftErr == nil && rightErr == nil {
			return NewNumberValue(leftNum + rightNum), nil
		}

		// Otherwise, do string concatenation
		return NewStringValue(v.String + other.String), nil
	}

	// If one is string and other is number, this is a type mismatch error in BASIC
	if v.Type == StringType || other.Type == StringType {
		return Value{}, fmt.Errorf("TYPE MISMATCH ERROR")
	}

	// Both are numbers, perform arithmetic addition
	return v.binaryArithmeticOp(other, func(left, right float64) float64 {
		return left + right
	})
}

// Subtract performs subtraction on two values
func (v Value) Subtract(other Value) (Value, error) {
	return v.binaryArithmeticOp(other, func(left, right float64) float64 {
		return left - right
	})
}

// Multiply performs multiplication on two values
func (v Value) Multiply(other Value) (Value, error) {
	return v.binaryArithmeticOp(other, func(left, right float64) float64 {
		return left * right
	})
}

// Divide performs division on two values
func (v Value) Divide(other Value) (Value, error) {
	return v.binaryArithmeticOpWithError(other, func(left, right float64) (float64, error) {
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	})
}

// Power performs exponentiation on two values
func (v Value) Power(other Value) (Value, error) {
	return v.binaryArithmeticOp(other, func(left, right float64) float64 {
		return math.Pow(left, right)
	})
}

// IsTrue determines if a value evaluates to true in BASIC conditional contexts
func (v Value) IsTrue() bool {
	switch v.Type {
	case NumberType:
		return v.Number != 0 // Non-zero numbers are true
	case StringType:
		return v.String != "" // Non-empty strings are true
	default:
		return false
	}
}

// Compare compares this value with another value using the specified operator
func (v Value) Compare(other Value, operator string) (bool, error) {
	// Handle comparison based on types
	if v.Type == NumberType && other.Type == NumberType {
		// Numeric comparison
		return compareNumbers(v.Number, other.Number, operator), nil
	} else if v.Type == StringType && other.Type == StringType {
		// String comparison
		return compareStrings(v.String, other.String, operator), nil
	} else {
		// Type mismatch
		return false, fmt.Errorf("TYPE MISMATCH ERROR")
	}
}

// compareNumbers performs numeric comparison
func compareNumbers(left, right float64, operator string) bool {
	switch operator {
	case "=":
		return left == right
	case "<>":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	default:
		return false // Invalid operator
	}
}

// compareStrings performs string comparison
func compareStrings(left, right string, operator string) bool {
	switch operator {
	case "=":
		return left == right
	case "<>":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	default:
		return false // Invalid operator
	}
}
