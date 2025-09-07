package parser

import (
	"basic-interpreter/types"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockInterpreterOperations implements InterpreterOperations for testing AST nodes
type MockInterpreterOperations struct {
	// Variables storage
	variables map[string]types.Value

	// I/O capture
	printedLines []string
	inputQueue   []string
	inputIndex   int

	// Control flow tracking
	gotoRequested   bool
	gotoTarget      int
	endRequested    bool
	stopRequested   bool
	gosubRequested  bool
	gosubTarget     int
	returnRequested bool

	// Error injection for testing
	getVariableError error
	setVariableError error
	printLineError   error
	readInputError   error
}

func newMockOps() *MockInterpreterOperations {
	return &MockInterpreterOperations{
		variables:    make(map[string]types.Value),
		printedLines: make([]string, 0),
		inputQueue:   make([]string, 0),
	}
}

func (m *MockInterpreterOperations) GetVariable(name string) (types.Value, error) {
	if m.getVariableError != nil {
		return types.Value{}, m.getVariableError
	}

	if value, exists := m.variables[name]; exists {
		return value, nil
	}

	// Return zero value for undefined variables (C64 BASIC behavior)
	if name[len(name)-1] == '$' {
		return types.NewStringValue(""), nil
	}
	return types.NewNumberValue(0), nil
}

func (m *MockInterpreterOperations) SetVariable(name string, value types.Value) error {
	if m.setVariableError != nil {
		return m.setVariableError
	}

	m.variables[name] = value
	return nil
}

func (m *MockInterpreterOperations) PrintLine(text string) error {
	if m.printLineError != nil {
		return m.printLineError
	}

	m.printedLines = append(m.printedLines, text)
	return nil
}

func (m *MockInterpreterOperations) ReadInput(prompt string) (string, error) {
	if m.readInputError != nil {
		return "", m.readInputError
	}

	if m.inputIndex >= len(m.inputQueue) {
		return "", errors.New("no input available")
	}

	result := m.inputQueue[m.inputIndex]
	m.inputIndex++
	return result, nil
}

func (m *MockInterpreterOperations) RequestGoto(targetLine int) error {
	m.gotoRequested = true
	m.gotoTarget = targetLine
	return nil
}

func (m *MockInterpreterOperations) RequestEnd() error {
	m.endRequested = true
	return nil
}

func (m *MockInterpreterOperations) RequestStop() error {
	m.stopRequested = true
	return nil
}

func (m *MockInterpreterOperations) RequestGosub(targetLine int) error {
	m.gosubRequested = true
	m.gosubTarget = targetLine
	return nil
}

func (m *MockInterpreterOperations) RequestReturn() error {
	m.returnRequested = true
	return nil
}

func (m *MockInterpreterOperations) NormalizeVariableName(name string) string {
	// Simple implementation for testing - just return as-is
	return name
}

// Loop control no-ops for AST unit testing
func (m *MockInterpreterOperations) BeginFor(variable string, end types.Value, step types.Value) error {
	return nil
}

func (m *MockInterpreterOperations) IterateFor(variable string) error {
	return nil
}

// Data management stub
func (m *MockInterpreterOperations) GetNextData() (types.Value, error) {
	return types.NewNumberValue(0), nil
}

// Helper methods for testing
func (m *MockInterpreterOperations) setInput(inputs []string) {
	m.inputQueue = inputs
	m.inputIndex = 0
}

func (m *MockInterpreterOperations) getOutput() []string {
	return m.printedLines
}

func (m *MockInterpreterOperations) setVariable(name string, value types.Value) {
	m.variables[name] = value
}

func TestPrintStatement_Execute(t *testing.T) {
	tests := []struct {
		name           string
		expression     Expression
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "print string literal",
			expression:     &StringLiteral{Value: "HELLO", Line: 1},
			expectedOutput: "HELLO",
			expectError:    false,
		},
		{
			name:           "print number literal",
			expression:     &NumberLiteral{Value: "42", Line: 1},
			expectedOutput: "42",
			expectError:    false,
		},
		{
			name:           "print empty string",
			expression:     &StringLiteral{Value: "", Line: 1},
			expectedOutput: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			stmt := &PrintStatement{Expression: tt.expression, Line: 1}

			err := stmt.Execute(mock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.Len(t, mock.getOutput(), 1)
				assert.Equal(t, tt.expectedOutput, mock.getOutput()[0])
			}
		})
	}
}

func TestPrintStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("expression evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		stmt := &PrintStatement{
			Expression: &VariableReference{Name: "A", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("print line error", func(t *testing.T) {
		mock := newMockOps()
		mock.printLineError = errors.New("print error")

		stmt := &PrintStatement{
			Expression: &StringLiteral{Value: "TEST", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}

func TestLetStatement_Execute(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		expression   Expression
		expectedType types.ValueType
		expectedVal  interface{}
	}{
		{
			name:         "assign number to variable",
			variable:     "A",
			expression:   &NumberLiteral{Value: "42", Line: 1},
			expectedType: types.NumberType,
			expectedVal:  42.0,
		},
		{
			name:         "assign string to variable",
			variable:     "A$",
			expression:   &StringLiteral{Value: "HELLO", Line: 1},
			expectedType: types.StringType,
			expectedVal:  "HELLO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			stmt := &LetStatement{
				Variable:   tt.variable,
				Expression: tt.expression,
				Line:       1,
			}

			err := stmt.Execute(mock)
			assert.NoError(t, err)

			value, exists := mock.variables[tt.variable]
			assert.True(t, exists)
			assert.Equal(t, tt.expectedType, value.Type)

			if tt.expectedType == types.NumberType {
				assert.Equal(t, tt.expectedVal, value.Number)
			} else {
				assert.Equal(t, tt.expectedVal, value.String)
			}
		})
	}
}

func TestLetStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("expression evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		stmt := &LetStatement{
			Variable:   "A",
			Expression: &VariableReference{Name: "B", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("set variable error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariableError = errors.New("set error")

		stmt := &LetStatement{
			Variable:   "A",
			Expression: &NumberLiteral{Value: "42", Line: 1},
			Line:       1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}

func TestEndStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &EndStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.endRequested)
}

func TestStopStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &StopStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.stopRequested)
}

func TestRunStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &RunStatement{Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
}

func TestGotoStatement_Execute(t *testing.T) {
	mock := newMockOps()
	stmt := &GotoStatement{TargetLine: 50, Line: 1}

	err := stmt.Execute(mock)

	assert.NoError(t, err)
	assert.True(t, mock.gotoRequested)
	assert.Equal(t, 50, mock.gotoTarget)
}

func TestInputStatement_Execute(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		input        string
		expectedType types.ValueType
		expectedVal  interface{}
		expectError  bool
	}{
		{
			name:         "numeric input to numeric variable",
			variable:     "A",
			input:        "42",
			expectedType: types.NumberType,
			expectedVal:  42.0,
			expectError:  false,
		},
		{
			name:         "string input to string variable",
			variable:     "A$",
			input:        "HELLO",
			expectedType: types.StringType,
			expectedVal:  "HELLO",
			expectError:  false,
		},
		{
			name:        "invalid numeric input",
			variable:    "A",
			input:       "ABC",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setInput([]string{tt.input})

			stmt := &InputStatement{Variable: tt.variable, Line: 1}

			err := stmt.Execute(mock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				value, exists := mock.variables[tt.variable]
				assert.True(t, exists)
				assert.Equal(t, tt.expectedType, value.Type)

				if tt.expectedType == types.NumberType {
					assert.Equal(t, tt.expectedVal, value.Number)
				} else {
					assert.Equal(t, tt.expectedVal, value.String)
				}
			}
		})
	}
}

func TestInputStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("read input error", func(t *testing.T) {
		mock := newMockOps()
		mock.readInputError = errors.New("input error")

		stmt := &InputStatement{Variable: "A", Line: 1}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("set variable error", func(t *testing.T) {
		mock := newMockOps()
		mock.setInput([]string{"42"})
		mock.setVariableError = errors.New("set error")

		stmt := &InputStatement{Variable: "A", Line: 1}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}

func TestIfStatement_Execute(t *testing.T) {
	tests := []struct {
		name            string
		conditionValue  types.Value
		expectExecution bool
	}{
		{
			name:            "true condition (non-zero number)",
			conditionValue:  types.NewNumberValue(1),
			expectExecution: true,
		},
		{
			name:            "false condition (zero)",
			conditionValue:  types.NewNumberValue(0),
			expectExecution: false,
		},
		{
			name:            "true condition (negative number)",
			conditionValue:  types.NewNumberValue(-5),
			expectExecution: true,
		},
		{
			name:            "true condition (non-empty string)",
			conditionValue:  types.NewStringValue("HELLO"),
			expectExecution: true,
		},
		{
			name:            "false condition (empty string)",
			conditionValue:  types.NewStringValue(""),
			expectExecution: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("CONDITION", tt.conditionValue)

			condition := &VariableReference{Name: "CONDITION", Line: 1}
			thenStmt := &PrintStatement{
				Expression: &StringLiteral{Value: "EXECUTED", Line: 1},
				Line:       1,
			}

			stmt := &IfStatement{
				Condition: condition,
				ThenStmt:  thenStmt,
				Line:      1,
			}

			err := stmt.Execute(mock)
			assert.NoError(t, err)

			if tt.expectExecution {
				assert.Len(t, mock.getOutput(), 1)
				assert.Equal(t, "EXECUTED", mock.getOutput()[0])
			} else {
				assert.Len(t, mock.getOutput(), 0)
			}
		})
	}
}

func TestIfStatement_Execute_ErrorCases(t *testing.T) {
	t.Run("condition evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		condition := &VariableReference{Name: "A", Line: 1}
		thenStmt := &PrintStatement{
			Expression: &StringLiteral{Value: "TEST", Line: 1},
			Line:       1,
		}

		stmt := &IfStatement{
			Condition: condition,
			ThenStmt:  thenStmt,
			Line:      1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})

	t.Run("then statement execution error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariable("A", types.NewNumberValue(1)) // true condition
		mock.printLineError = errors.New("print error")

		condition := &VariableReference{Name: "A", Line: 1}
		thenStmt := &PrintStatement{
			Expression: &StringLiteral{Value: "TEST", Line: 1},
			Line:       1,
		}

		stmt := &IfStatement{
			Condition: condition,
			ThenStmt:  thenStmt,
			Line:      1,
		}

		err := stmt.Execute(mock)
		assert.Error(t, err)
	})
}

func TestStringLiteral_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"simple string", "HELLO", "HELLO"},
		{"empty string", "", ""},
		{"string with spaces", "HELLO WORLD", "HELLO WORLD"},
		{"string with numbers", "ABC123", "ABC123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			expr := &StringLiteral{Value: tt.value, Line: 1}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.StringType, result.Type)
			assert.Equal(t, tt.expected, result.String)
		})
	}
}

func TestNumberLiteral_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected float64
	}{
		{"integer", "42", 42.0},
		{"float", "42.5", 42.5},
		{"negative", "-123", -123.0},
		{"zero", "0", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			expr := &NumberLiteral{Value: tt.value, Line: 1}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestNumberLiteral_Evaluate_InvalidNumber(t *testing.T) {
	mock := newMockOps()
	expr := &NumberLiteral{Value: "invalid", Line: 1}

	result, err := expr.Evaluate(mock)

	// ParseValue treats invalid numbers as strings, which is correct BASIC behavior
	assert.NoError(t, err)
	assert.Equal(t, types.StringType, result.Type)
	assert.Equal(t, "invalid", result.String)
}

func TestVariableReference_Evaluate(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		setValue     types.Value
		expectedType types.ValueType
		expectedVal  interface{}
	}{
		{
			name:         "numeric variable",
			variable:     "A",
			setValue:     types.NewNumberValue(42),
			expectedType: types.NumberType,
			expectedVal:  42.0,
		},
		{
			name:         "string variable",
			variable:     "A$",
			setValue:     types.NewStringValue("HELLO"),
			expectedType: types.StringType,
			expectedVal:  "HELLO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable(tt.variable, tt.setValue)

			expr := &VariableReference{Name: tt.variable, Line: 1}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, result.Type)

			if tt.expectedType == types.NumberType {
				assert.Equal(t, tt.expectedVal, result.Number)
			} else {
				assert.Equal(t, tt.expectedVal, result.String)
			}
		})
	}
}

func TestVariableReference_Evaluate_UndefinedVariable(t *testing.T) {
	tests := []struct {
		name         string
		variable     string
		expectedType types.ValueType
		expectedVal  interface{}
	}{
		{
			name:         "undefined numeric variable returns 0",
			variable:     "A",
			expectedType: types.NumberType,
			expectedVal:  0.0,
		},
		{
			name:         "undefined string variable returns empty string",
			variable:     "A$",
			expectedType: types.StringType,
			expectedVal:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			expr := &VariableReference{Name: tt.variable, Line: 1}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, result.Type)

			if tt.expectedType == types.NumberType {
				assert.Equal(t, tt.expectedVal, result.Number)
			} else {
				assert.Equal(t, tt.expectedVal, result.String)
			}
		})
	}
}

func TestVariableReference_Evaluate_Error(t *testing.T) {
	mock := newMockOps()
	mock.getVariableError = errors.New("variable error")

	expr := &VariableReference{Name: "A", Line: 1}

	_, err := expr.Evaluate(mock)

	assert.Error(t, err)
}

func TestBinaryOperation_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     types.Value
		operator string
		right    types.Value
		expected float64
	}{
		{"addition", types.NewNumberValue(5), "+", types.NewNumberValue(3), 8},
		{"subtraction", types.NewNumberValue(10), "-", types.NewNumberValue(4), 6},
		{"multiplication", types.NewNumberValue(6), "*", types.NewNumberValue(7), 42},
		{"division", types.NewNumberValue(15), "/", types.NewNumberValue(3), 5},
		{"power", types.NewNumberValue(2), "^", types.NewNumberValue(3), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", tt.left)
			mock.setVariable("RIGHT", tt.right)

			expr := &BinaryOperation{
				Left:     &VariableReference{Name: "LEFT", Line: 1},
				Operator: tt.operator,
				Right:    &VariableReference{Name: "RIGHT", Line: 1},
				Line:     1,
			}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestBinaryOperation_Evaluate_StringConcatenation(t *testing.T) {
	mock := newMockOps()
	mock.setVariable("LEFT", types.NewStringValue("HELLO"))
	mock.setVariable("RIGHT", types.NewStringValue("WORLD"))

	expr := &BinaryOperation{
		Left:     &VariableReference{Name: "LEFT", Line: 1},
		Operator: "+",
		Right:    &VariableReference{Name: "RIGHT", Line: 1},
		Line:     1,
	}

	result, err := expr.Evaluate(mock)

	assert.NoError(t, err)
	assert.Equal(t, types.StringType, result.Type)
	assert.Equal(t, "HELLOWORLD", result.String)
}

func TestBinaryOperation_Evaluate_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		operator string
	}{
		{"unknown operator", "**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", types.NewNumberValue(5))
			mock.setVariable("RIGHT", types.NewNumberValue(3))

			expr := &BinaryOperation{
				Left:     &VariableReference{Name: "LEFT", Line: 1},
				Operator: tt.operator,
				Right:    &VariableReference{Name: "RIGHT", Line: 1},
				Line:     1,
			}

			_, err := expr.Evaluate(mock)

			assert.Error(t, err)
		})
	}
}

func TestBinaryOperation_Evaluate_LeftEvaluationError(t *testing.T) {
	mock := newMockOps()
	mock.getVariableError = errors.New("variable error")

	expr := &BinaryOperation{
		Left:     &VariableReference{Name: "A", Line: 1},
		Operator: "+",
		Right:    &NumberLiteral{Value: "3", Line: 1},
		Line:     1,
	}

	_, err := expr.Evaluate(mock)

	assert.Error(t, err)
}

func TestUnaryOperation_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		operand  types.Value
		expected float64
	}{
		{"negation", "-", types.NewNumberValue(5), -5},
		{"unary plus", "+", types.NewNumberValue(5), 5},
		{"negate negative", "-", types.NewNumberValue(-3), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("OPERAND", tt.operand)

			expr := &UnaryOperation{
				Operator: tt.operator,
				Right:    &VariableReference{Name: "OPERAND", Line: 1},
				Line:     1,
			}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestUnaryOperation_Evaluate_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		operand  types.Value
		errMsg   string
	}{
		{"negate string", "-", types.NewStringValue("HELLO"), "cannot negate non-numeric value"},
		{"unary plus on string", "+", types.NewStringValue("HELLO"), "cannot apply unary plus to non-numeric value"},
		{"unknown operator", "!", types.NewNumberValue(5), "unknown unary operator: !"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("OPERAND", tt.operand)

			expr := &UnaryOperation{
				Operator: tt.operator,
				Right:    &VariableReference{Name: "OPERAND", Line: 1},
				Line:     1,
			}

			_, err := expr.Evaluate(mock)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestComparisonExpression_Evaluate(t *testing.T) {
	tests := []struct {
		name     string
		left     types.Value
		operator string
		right    types.Value
		expected float64
	}{
		// Numeric comparisons
		{"equal true", types.NewNumberValue(5), "=", types.NewNumberValue(5), 1},
		{"equal false", types.NewNumberValue(5), "=", types.NewNumberValue(3), 0},
		{"not equal true", types.NewNumberValue(5), "<>", types.NewNumberValue(3), 1},
		{"not equal false", types.NewNumberValue(5), "<>", types.NewNumberValue(5), 0},
		{"less than true", types.NewNumberValue(3), "<", types.NewNumberValue(5), 1},
		{"less than false", types.NewNumberValue(5), "<", types.NewNumberValue(3), 0},
		{"greater than true", types.NewNumberValue(5), ">", types.NewNumberValue(3), 1},
		{"greater than false", types.NewNumberValue(3), ">", types.NewNumberValue(5), 0},
		{"less equal true", types.NewNumberValue(3), "<=", types.NewNumberValue(5), 1},
		{"less equal equal", types.NewNumberValue(5), "<=", types.NewNumberValue(5), 1},
		{"less equal false", types.NewNumberValue(5), "<=", types.NewNumberValue(3), 0},
		{"greater equal true", types.NewNumberValue(5), ">=", types.NewNumberValue(3), 1},
		{"greater equal equal", types.NewNumberValue(5), ">=", types.NewNumberValue(5), 1},
		{"greater equal false", types.NewNumberValue(3), ">=", types.NewNumberValue(5), 0},

		// String comparisons
		{"string equal true", types.NewStringValue("HELLO"), "=", types.NewStringValue("HELLO"), 1},
		{"string equal false", types.NewStringValue("HELLO"), "=", types.NewStringValue("WORLD"), 0},
		{"string not equal true", types.NewStringValue("HELLO"), "<>", types.NewStringValue("WORLD"), 1},
		{"string not equal false", types.NewStringValue("HELLO"), "<>", types.NewStringValue("HELLO"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockOps()
			mock.setVariable("LEFT", tt.left)
			mock.setVariable("RIGHT", tt.right)

			expr := &ComparisonExpression{
				Left:     &VariableReference{Name: "LEFT", Line: 1},
				Operator: tt.operator,
				Right:    &VariableReference{Name: "RIGHT", Line: 1},
				Line:     1,
			}

			result, err := expr.Evaluate(mock)

			assert.NoError(t, err)
			assert.Equal(t, types.NumberType, result.Type)
			assert.Equal(t, tt.expected, result.Number)
		})
	}
}

func TestComparisonExpression_Evaluate_ErrorCases(t *testing.T) {
	t.Run("left evaluation error", func(t *testing.T) {
		mock := newMockOps()
		mock.getVariableError = errors.New("variable error")

		expr := &ComparisonExpression{
			Left:     &VariableReference{Name: "A", Line: 1},
			Operator: "=",
			Right:    &NumberLiteral{Value: "5", Line: 1},
			Line:     1,
		}

		_, err := expr.Evaluate(mock)

		assert.Error(t, err)
	})

	t.Run("comparison error", func(t *testing.T) {
		mock := newMockOps()
		mock.setVariable("LEFT", types.NewNumberValue(5))
		mock.setVariable("RIGHT", types.NewStringValue("HELLO"))

		expr := &ComparisonExpression{
			Left:     &VariableReference{Name: "LEFT", Line: 1},
			Operator: "<",
			Right:    &VariableReference{Name: "RIGHT", Line: 1},
			Line:     1,
		}

		_, err := expr.Evaluate(mock)

		assert.Error(t, err)
	})
}
