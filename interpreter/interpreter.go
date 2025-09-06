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
	Variable      string      // Normalized loop variable name
	EndValue      types.Value // Target end value
	StepValue     types.Value // Step value (default 1)
	AfterForIndex int         // Target line index to jump back to (line after FOR)
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
	linePos   map[int]int            // Maps line numbers to their index position
	forStack  []ForLoopContext       // Stack of active FOR loops for nested loop support
	maxSteps  int                    // Maximum number of execution steps before infinite loop protection kicks in
	stepCount int                    // Current step count during execution
	pc        int                    // Program counter: current line index
	jumped    bool                   // Indicates a jump occurred during statement execution
	halted    bool                   // Indicates END/STOP was requested
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime:   rt,
		variables: make(map[string]types.Value),
		lineIndex: make(map[int]*parser.Line),
		linePos:   make(map[int]int),
		forStack:  make([]ForLoopContext, 0),
		maxSteps:  1000, // Default maximum steps
		stepCount: 0,
		pc:        0,
		jumped:    false,
		halted:    false,
	}
}

// SetMaxSteps sets the maximum number of execution steps before infinite loop protection
func (i *Interpreter) SetMaxSteps(maxSteps int) {
	i.maxSteps = maxSteps
}

// pushForLoop pushes a new FOR loop context onto the stack
func (i *Interpreter) pushForLoop(variable string, endValue types.Value, stepValue types.Value, afterForIndex int) {
	norm := i.NormalizeVariableName(variable)
	forLoop := ForLoopContext{
		Variable:      norm,
		EndValue:      endValue,
		StepValue:     stepValue,
		AfterForIndex: afterForIndex,
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
	norm := i.NormalizeVariableName(variable)
	for j := len(i.forStack) - 1; j >= 0; j-- {
		if i.forStack[j].Variable == norm {
			return &i.forStack[j]
		}
	}
	return nil
}

// Execute runs a BASIC program
func (i *Interpreter) Execute(program *parser.Program) error {
	// Reset step counter for new execution
	i.stepCount = 0
	i.halted = false
	i.jumped = false

	// Build line number index for GOTO statements
	i.buildLineIndex(program)

	// Execute program with program counter for GOTO support
	return i.executeWithProgramCounter(program)
}

// buildLineIndex creates a map from line numbers to Line nodes
func (i *Interpreter) buildLineIndex(program *parser.Program) {
	i.lineIndex = make(map[int]*parser.Line)
	i.linePos = make(map[int]int)
	for idx, line := range program.Lines {
		i.lineIndex[line.Number] = line
		i.linePos[line.Number] = idx
	}
}

// executeWithProgramCounter executes program with support for GOTO jumps using polymorphic dispatch
func (i *Interpreter) executeWithProgramCounter(program *parser.Program) error {
	if len(program.Lines) == 0 {
		return nil
	}

	// Start execution at the first line
	i.pc = 0

	for i.pc < len(program.Lines) {
		line := program.Lines[i.pc]

		for _, stmt := range line.Statements {
			// Increment step counter and check for infinite loop protection
			i.stepCount++
			if i.maxSteps > 0 && i.stepCount > i.maxSteps {
				return fmt.Errorf("?INFINITE LOOP ERROR")
			}

			// Polymorphic dispatch - AST node executes itself using double dispatch
			err := stmt.Execute(i)
			if err != nil {
				// Regular error - wrap with line number
				return i.wrapErrorWithLine(err, line.Number)
			}

			// After successful execution, check for END/STOP or GOTO performed via ops
			if i.halted {
				return nil
			}
			if i.jumped {
				i.jumped = false
				goto nextLine
			}
		}

		// Move to next line
		i.pc++
	nextLine:
	}

	return nil
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
	case strings.Contains(errMsg, "ILLEGAL QUANTITY ERROR"):
		return fmt.Errorf("?ILLEGAL QUANTITY ERROR IN %d", lineNumber)
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
	// Resolve target line to index and set jump state
	targetLineIndex, found := i.linePos[targetLine]
	if !found {
		// We don't have the source line number here; the caller's line will wrap this error
		return fmt.Errorf("?UNDEFINED STATEMENT ERROR")
	}
	i.pc = targetLineIndex
	i.jumped = true
	return nil
}

// RequestEnd requests program termination
func (i *Interpreter) RequestEnd() error {
	i.halted = true
	return nil
}

// RequestStop requests program stop
func (i *Interpreter) RequestStop() error {
	i.halted = true
	return nil
}

// NormalizeVariableName truncates variable name to first 2 characters (C64 BASIC behavior)
func (i *Interpreter) NormalizeVariableName(name string) string {
	if len(name) > 2 {
		return name[:2]
	}
	return name
}

// BeginFor starts a FOR loop by pushing a loop context
func (i *Interpreter) BeginFor(variable string, end types.Value, step types.Value) error {
	// Validate step (cannot be zero)
	if step.Type != types.NumberType || step.Number == 0 {
		return fmt.Errorf("ILLEGAL QUANTITY ERROR")
	}
	// Jump back target is the line after the FOR statement
	i.pushForLoop(variable, end, step, i.pc+1)
	return nil
}

// IterateFor performs a NEXT iteration; variable may be empty to use the most recent loop
func (i *Interpreter) IterateFor(variableName string) error {
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

	// Increment the loop variable by the step value
	newValue, err := currentValue.Add(forLoop.StepValue)
	if err != nil {
		return err
	}

	// Determine comparison based on step direction
	cmpOp := "<="
	if forLoop.StepValue.Number < 0 {
		cmpOp = ">="
	}
	shouldContinue, err := newValue.Compare(forLoop.EndValue, cmpOp)
	if err != nil {
		return err
	}

	if shouldContinue {
		// Update loop variable and jump back to the line after FOR
		err = i.SetVariable(forLoop.Variable, newValue)
		if err != nil {
			return err
		}
		// Signal jump to AfterForIndex
		i.pc = forLoop.AfterForIndex
		i.jumped = true
		return nil
	}

	// Loop finished - pop the loop from stack
	i.popForLoop()
	return nil
}
