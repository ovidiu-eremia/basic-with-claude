// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"fmt"
	"strings"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

// ForLoopContext represents an active FOR loop state
type ForLoopContext struct {
	Variable     string      // Loop variable name
	EndValue     types.Value // Target end value
	StepValue    types.Value // Step value (default 1)
	ForLineIndex int         // Index of the FOR statement line in program
}

// RuntimeError represents an error that occurred during program execution
type RuntimeError struct {
	Message  string
	Position lexer.Position
}

// Error implements the error interface
func (re *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at line %d, column %d: %s", re.Position.Line, re.Position.Column, re.Message)
}

// Interpreter executes BASIC programs by walking the AST
type Interpreter struct {
	runtime   runtime.Runtime
	variables map[string]types.Value // Variable storage using proper Value types
	lineIndex map[int]*parser.Line   // Maps line numbers to Line nodes for GOTO
	forStack  []ForLoopContext       // Stack of active FOR loops for nested loop support
	maxSteps  int                    // Maximum number of execution steps before infinite loop protection kicks in
	stepCount int                    // Current step count during execution
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime:   rt,
		variables: make(map[string]types.Value),
		lineIndex: make(map[int]*parser.Line),
		forStack:  make([]ForLoopContext, 0),
		maxSteps:  1000, // Default maximum steps
		stepCount: 0,
	}
}

// SetMaxSteps sets the maximum number of execution steps before infinite loop protection
func (i *Interpreter) SetMaxSteps(maxSteps int) {
	i.maxSteps = maxSteps
}

// pushForLoop pushes a new FOR loop context onto the stack
func (i *Interpreter) pushForLoop(variable string, endValue types.Value, stepValue types.Value, forLineIndex int) {
	forLoop := ForLoopContext{
		Variable:     variable,
		EndValue:     endValue,
		StepValue:    stepValue,
		ForLineIndex: forLineIndex,
	}
	i.forStack = append(i.forStack, forLoop)
}

// popForLoop removes the top FOR loop from the stack
func (i *Interpreter) popForLoop() *ForLoopContext {
	if len(i.forStack) == 0 {
		return nil
	}
	top := i.forStack[len(i.forStack)-1]
	i.forStack = i.forStack[:len(i.forStack)-1]
	return &top
}

// peekForLoop returns the top FOR loop without removing it
func (i *Interpreter) peekForLoop() *ForLoopContext {
	if len(i.forStack) == 0 {
		return nil
	}
	return &i.forStack[len(i.forStack)-1]
}

// findForLoopByVariable finds a FOR loop on the stack by variable name
func (i *Interpreter) findForLoopByVariable(variable string) *ForLoopContext {
	// Search from top of stack (most recent) to bottom
	for j := len(i.forStack) - 1; j >= 0; j-- {
		if i.forStack[j].Variable == variable {
			return &i.forStack[j]
		}
	}
	return nil
}

// Execute runs a BASIC program
func (i *Interpreter) Execute(program *parser.Program) error {
	// Reset step counter for new execution
	i.stepCount = 0

	// Build line number index for GOTO statements
	i.buildLineIndex(program)

	// Execute program with program counter for GOTO support
	return i.executeWithProgramCounter(program)
}

// buildLineIndex creates a map from line numbers to Line nodes
func (i *Interpreter) buildLineIndex(program *parser.Program) {
	for _, line := range program.Lines {
		i.lineIndex[line.Number] = line
	}
}

// executeWithProgramCounter executes program with support for GOTO jumps using polymorphic dispatch
func (i *Interpreter) executeWithProgramCounter(program *parser.Program) error {
	if len(program.Lines) == 0 {
		return nil
	}

	// Start execution at the first line
	currentLineIndex := 0

	for currentLineIndex < len(program.Lines) {
		line := program.Lines[currentLineIndex]

		for _, stmt := range line.Statements {
			// Increment step counter and check for infinite loop protection
			i.stepCount++
			if i.maxSteps > 0 && i.stepCount > i.maxSteps {
				return fmt.Errorf("?INFINITE LOOP ERROR")
			}

			// Polymorphic dispatch - AST node executes itself using double dispatch
			err := stmt.Execute(i)
			if err != nil {
				// Handle control flow requests
				if gotoCtrl, ok := err.(*parser.GotoControl); ok {
					targetLineIndex, found := i.findLineIndex(program, gotoCtrl.TargetLine)
					if !found {
						return fmt.Errorf("?UNDEFINED STATEMENT ERROR IN %d", line.Number)
					}
					currentLineIndex = targetLineIndex
					goto nextLine
				}

				if _, ok := err.(*parser.EndControl); ok {
					return nil
				}

				if _, ok := err.(*parser.StopControl); ok {
					return nil
				}

				// Handle FOR loop initialization
				if forCtrl, ok := err.(*parser.ForControl); ok {
					// Push FOR loop onto stack with default step of 1
					stepValue := types.NewNumberValue(1)
					i.pushForLoop(forCtrl.Variable, forCtrl.EndValue, stepValue, currentLineIndex)
					// Continue with next statement (don't jump anywhere)
					continue
				}

				// Handle NEXT statement
				if nextCtrl, ok := err.(*parser.NextControl); ok {
					err := i.handleNextStatement(program, nextCtrl.Variable, &currentLineIndex)
					if err != nil {
						// Check if it's a loop jump
						if loopCtrl, ok := err.(*parser.LoopControl); ok {
							currentLineIndex = loopCtrl.TargetLineIndex
							goto nextLine
						}
						return err
					}
					// Loop finished normally - continue to next line with increment
					break // Exit the statement loop to allow normal line increment
				}

				// Regular error - wrap with line number
				return i.wrapErrorWithLine(err, line.Number)
			}
		}

		// Move to next line
		currentLineIndex++
	nextLine:
	}

	return nil
}

// findLineIndex finds the index of a line with the given line number
func (i *Interpreter) findLineIndex(program *parser.Program, lineNumber int) (int, bool) {
	for index, line := range program.Lines {
		if line.Number == lineNumber {
			return index, true
		}
	}
	return 0, false
}

// wrapErrorWithLine wraps an error with C64 BASIC format including line number
func (i *Interpreter) wrapErrorWithLine(err error, lineNumber int) error {
	// Check if it's already a C64 format error (starts with ?)
	errMsg := err.Error()
	if len(errMsg) > 0 && errMsg[0] == '?' {
		return err // Already formatted
	}

	// Convert common errors to C64 BASIC format
	switch {
	case strings.Contains(errMsg, "division by zero"):
		return fmt.Errorf("?DIVISION BY ZERO ERROR IN %d", lineNumber)
	case strings.Contains(errMsg, "TYPE MISMATCH ERROR"):
		return fmt.Errorf("?TYPE MISMATCH ERROR IN %d", lineNumber)
	default:
		return fmt.Errorf("?ERROR IN %d: %s", lineNumber, errMsg)
	}
}

// InterpreterOperations interface implementation
// These methods enable double dispatch from AST nodes back to interpreter

// GetVariable retrieves a variable value by name
func (i *Interpreter) GetVariable(name string) (types.Value, error) {
	normalizedName := i.NormalizeVariableName(name)
	if value, exists := i.variables[normalizedName]; exists {
		return value, nil
	}

	// Default values
	if strings.HasSuffix(name, "$") {
		return types.NewStringValue(""), nil
	}
	return types.NewNumberValue(0), nil
}

// SetVariable sets a variable value with type checking
func (i *Interpreter) SetVariable(name string, value types.Value) error {
	// Type check: string variables can only hold strings, numeric variables can only hold numbers
	isStringVariable := strings.HasSuffix(name, "$")
	if isStringVariable && value.Type != types.StringType {
		return fmt.Errorf("TYPE MISMATCH ERROR")
	}
	if !isStringVariable && value.Type != types.NumberType {
		return fmt.Errorf("TYPE MISMATCH ERROR")
	}

	normalizedName := i.NormalizeVariableName(name)
	i.variables[normalizedName] = value
	return nil
}

// PrintLine outputs text to the runtime environment
func (i *Interpreter) PrintLine(text string) error {
	return i.runtime.PrintLine(text)
}

// ReadInput reads input from the runtime environment
func (i *Interpreter) ReadInput(prompt string) (string, error) {
	return i.runtime.Input(prompt)
}

// RequestGoto requests a GOTO control flow change
func (i *Interpreter) RequestGoto(targetLine int) error {
	return &parser.GotoControl{TargetLine: targetLine}
}

// RequestEnd requests program termination
func (i *Interpreter) RequestEnd() error {
	return &parser.EndControl{}
}

// RequestStop requests program stop
func (i *Interpreter) RequestStop() error {
	return &parser.StopControl{}
}

// NormalizeVariableName truncates variable name to first 2 characters (C64 BASIC behavior)
func (i *Interpreter) NormalizeVariableName(name string) string {
	if len(name) > 2 {
		return name[:2]
	}
	return name
}

// handleNextStatement handles NEXT statement execution and loop iteration
func (i *Interpreter) handleNextStatement(program *parser.Program, variableName string, currentLineIndex *int) error {
	// Find the appropriate FOR loop context
	var forLoop *ForLoopContext
	if variableName != "" {
		// NEXT with variable name - find specific loop
		forLoop = i.findForLoopByVariable(variableName)
		if forLoop == nil {
			return fmt.Errorf("?NEXT WITHOUT FOR ERROR")
		}
	} else {
		// NEXT without variable name - use most recent loop
		forLoop = i.peekForLoop()
		if forLoop == nil {
			return fmt.Errorf("?NEXT WITHOUT FOR ERROR")
		}
	}

	// Get current value of loop variable
	currentValue, err := i.GetVariable(forLoop.Variable)
	if err != nil {
		return err
	}

	// Increment the loop variable by the step value (default 1)
	newValue, err := currentValue.Add(forLoop.StepValue)
	if err != nil {
		return err
	}

	// Check if loop should continue
	shouldContinue, err := newValue.Compare(forLoop.EndValue, "<=")
	if err != nil {
		return err
	}

	if shouldContinue {
		// Update loop variable and jump back to the line after FOR
		err = i.SetVariable(forLoop.Variable, newValue)
		if err != nil {
			return err
		}
		// Jump back to the line after the FOR statement
		targetIndex := forLoop.ForLineIndex + 1
		// Return LoopControl to indicate loop jump
		return &parser.LoopControl{TargetLineIndex: targetIndex}
	} else {
		// Loop finished - pop the loop from stack
		i.popForLoop()
		// Continue with next statement after NEXT
		return nil
	}
}
