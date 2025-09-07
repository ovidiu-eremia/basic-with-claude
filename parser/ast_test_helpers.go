package parser

import (
	"errors"

	"basic-interpreter/types"
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
