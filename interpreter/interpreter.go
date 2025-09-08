// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"fmt"
	"math"
	"strings"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
	"basic-interpreter/types"
)

// Predefined errors for interpreter-level conditions
var (
	ErrIllegalQuantity    = fmt.Errorf("?ILLEGAL QUANTITY ERROR")
	ErrNextWithoutFor     = fmt.Errorf("?NEXT WITHOUT FOR ERROR")
	ErrUndefinedStatement = fmt.Errorf("?UNDEFINED STATEMENT ERROR")
	ErrReturnWithoutGosub = fmt.Errorf("?RETURN WITHOUT GOSUB ERROR")
	ErrStackOverflow      = fmt.Errorf("?OUT OF MEMORY ERROR")
	ErrOutOfData          = fmt.Errorf("?OUT OF DATA ERROR")
)

// ForLoopContext represents an active FOR loop state
type ForLoopContext struct {
	Variable          string      // Normalized loop variable name
	EndValue          types.Value // Target end value
	StepValue         types.Value // Step value (default 1)
	AfterForLineIndex int         // Target line index to jump back to
	AfterForStmtIndex int         // Target statement index within the line (for colon-separated statements)
}

// CallContext represents an active GOSUB call state
type CallContext struct {
	ReturnLineIndex int // Line index to return to after RETURN
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
	runtime      runtime.Runtime
	variables    map[string]types.Value // Variable storage using proper Value types
	lineIndex    map[int]*parser.Line   // Maps line numbers to Line nodes for GOTO
	linePos      map[int]int            // Maps line numbers to their index position
	forStack     *Stack[ForLoopContext] // Stack of active FOR loops for nested loop support
	callStack    *Stack[CallContext]    // Stack of active GOSUB calls for nested subroutine support
	maxSteps     int                    // Maximum number of execution steps before infinite loop protection kicks in
	maxCallDepth int                    // Maximum call stack depth before stack overflow error
	stepCount    int                    // Current step count during execution
	pc           int                    // Program counter: current line index
	stmtIndex    int                    // Current statement index within current line
	jumped       bool                   // Indicates a jump occurred during statement execution
	halted       bool                   // Indicates END/STOP was requested
	stmtJumped   bool                   // Indicates a statement-level jump occurred (for FOR loop completion)

	// DATA/READ state
	dataValues  []types.Value // Collected DATA values
	dataPointer int           // Current READ pointer

	// No RNG here; delegate randomness to runtime
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	maxCallDepth := 100 // Default maximum call depth
	return &Interpreter{
		runtime:      rt,
		variables:    make(map[string]types.Value),
		lineIndex:    make(map[int]*parser.Line),
		linePos:      make(map[int]int),
		forStack:     NewStack[ForLoopContext](maxCallDepth), // Use same limit for FOR loops
		callStack:    NewStack[CallContext](maxCallDepth),
		maxSteps:     1000, // Default maximum steps
		maxCallDepth: maxCallDepth,
		stepCount:    0,
		pc:           0,
		stmtIndex:    0,
		jumped:       false,
		halted:       false,
		stmtJumped:   false,
	}
}

// SetMaxSteps sets the maximum number of execution steps before infinite loop protection
func (i *Interpreter) SetMaxSteps(maxSteps int) {
	i.maxSteps = maxSteps
}

// pushForLoop pushes a new FOR loop context onto the stack
func (i *Interpreter) pushForLoop(variable string, endValue types.Value, stepValue types.Value, afterForLineIndex int, afterForStmtIndex int) error {
	norm := i.NormalizeVariableName(variable)
	forLoop := ForLoopContext{
		Variable:          norm,
		EndValue:          endValue,
		StepValue:         stepValue,
		AfterForLineIndex: afterForLineIndex,
		AfterForStmtIndex: afterForStmtIndex,
	}
	return i.forStack.Push(forLoop)
}

// popForLoop removes the top FOR loop from the stack
func (i *Interpreter) popForLoop() *ForLoopContext {
	return i.forStack.Pop()
}

// peekForLoop returns the top FOR loop without removing it
func (i *Interpreter) peekForLoop() *ForLoopContext {
	return i.forStack.Peek()
}

// findForLoopByVariable finds a FOR loop on the stack by variable name
func (i *Interpreter) findForLoopByVariable(variable string) *ForLoopContext {
	norm := i.NormalizeVariableName(variable)
	return i.forStack.FindByPredicate(func(ctx ForLoopContext) bool {
		return ctx.Variable == norm
	})
}

// pushCallContext pushes a new call context onto the call stack
func (i *Interpreter) pushCallContext(returnLineIndex int) error {
	callContext := CallContext{
		ReturnLineIndex: returnLineIndex,
	}
	return i.callStack.Push(callContext)
}

// popCallContext removes the top call context from the stack
func (i *Interpreter) popCallContext() *CallContext {
	return i.callStack.Pop()
}

// Execute runs a BASIC program
func (i *Interpreter) Execute(program *parser.Program) error {
	// Reset step counter for new execution
	i.stepCount = 0
	i.halted = false
	i.jumped = false

	// Build line number index for GOTO statements
	i.buildLineIndex(program)

	// Collect DATA values before execution
	i.collectData(program)

	// Execute program with program counter for GOTO support
	return i.executeWithProgramCounter(program)
}

// collectData scans the program and collects all DATA values in order
func (i *Interpreter) collectData(program *parser.Program) {
	i.dataValues = i.dataValues[:0]
	i.dataPointer = 0
	for _, line := range program.Lines {
		for _, stmt := range line.Statements {
			if ds, ok := stmt.(*parser.DataStatement); ok {
				for _, expr := range ds.Values {
					val, err := expr.Evaluate(i)
					if err == nil {
						i.dataValues = append(i.dataValues, val)
					}
				}
			}
		}
	}
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
	i.stmtIndex = 0

	for i.pc < len(program.Lines) {
		line := program.Lines[i.pc]

		// Handle statement-level jumps (from FOR loop completion)
		if i.stmtJumped {
			i.stmtJumped = false
			// stmtIndex is already set by the jump, continue from there
		} else {
			i.stmtIndex = 0
		}

		for i.stmtIndex < len(line.Statements) {
			stmt := line.Statements[i.stmtIndex]

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
			if i.stmtJumped {
				goto nextLine // Continue from the jumped-to position
			}

			// Move to next statement
			i.stmtIndex++
		}

		// Move to next line
		i.pc++
	nextLine:
	}

	return nil
}

// wrapErrorWithLine wraps an error with C64 BASIC format including line number
func (i *Interpreter) wrapErrorWithLine(err error, lineNumber int) error {
	msg := err.Error()
	if len(msg) > 0 && msg[0] == '?' {
		// If already C64-style, append line if not present
		if strings.Contains(msg, " IN ") {
			return err
		}
		return fmt.Errorf("%s IN %d", msg, lineNumber)
	}
	return fmt.Errorf("?ERROR IN %d: %s", lineNumber, msg)
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
		return types.ErrTypeMismatch
	}
	if !isStringVariable && value.Type != types.NumberType {
		return types.ErrTypeMismatch
	}

	normalizedName := i.NormalizeVariableName(name)
	i.variables[normalizedName] = value
	return nil
}

// PrintLine outputs text to the runtime environment
func (i *Interpreter) PrintLine(text string) error {
	return i.runtime.PrintLine(text)
}

// Print outputs text without a newline
func (i *Interpreter) Print(text string) error {
	return i.runtime.Print(text)
}

// ReadInput reads input from the runtime environment
func (i *Interpreter) ReadInput(prompt string) (string, error) {
	return i.runtime.Input(prompt)
}

// GetNextData returns the next DATA value, or error if none remain
func (i *Interpreter) GetNextData() (types.Value, error) {
	if i.dataPointer >= len(i.dataValues) {
		return types.Value{}, ErrOutOfData
	}
	v := i.dataValues[i.dataPointer]
	i.dataPointer++
	return v, nil
}

// EvaluateFunction evaluates built-in functions
func (i *Interpreter) EvaluateFunction(functionName string, args []parser.Expression) (types.Value, error) {
	// Evaluate all arguments first
	argValues := make([]types.Value, len(args))
	for idx, arg := range args {
		val, err := arg.Evaluate(i)
		if err != nil {
			return types.Value{}, err
		}
		argValues[idx] = val
	}

	// Dispatch to specific function implementation
	switch strings.ToUpper(functionName) {
	case "LEN":
		return i.evaluateLenFunction(argValues)
	case "LEFT$":
		return i.evaluateLeftFunction(argValues)
	case "RIGHT$":
		return i.evaluateRightFunction(argValues)
	case "MID$":
		return i.evaluateMidFunction(argValues)
	case "CHR$":
		return i.evaluateChrFunction(argValues)
	case "ASC":
		return i.evaluateAscFunction(argValues)
	case "STR$":
		return i.evaluateStrFunction(argValues)
	case "VAL":
		return i.evaluateValFunction(argValues)
	case "RND":
		return i.evaluateRndFunction(argValues)
	case "ABS":
		return i.evaluateAbsFunction(argValues)
	case "INT":
		return i.evaluateIntFunction(argValues)
	case "SQR":
		return i.evaluateSqrFunction(argValues)
	default:
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: unknown function %s", functionName)
	}
}

// RequestGoto requests a GOTO control flow change
func (i *Interpreter) RequestGoto(targetLine int) error {
	// Resolve target line to index and set jump state
	targetLineIndex, found := i.linePos[targetLine]
	if !found {
		// We don't have the source line number here; the caller's line will wrap this error
		return ErrUndefinedStatement
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

// RequestGosub requests a GOSUB jump to a target line
func (i *Interpreter) RequestGosub(targetLine int) error {
	// First, push current position + 1 to call stack for RETURN
	if err := i.pushCallContext(i.pc + 1); err != nil {
		return err
	}

	// Then request jump to target line
	return i.RequestGoto(targetLine)
}

// RequestReturn requests a RETURN from current subroutine
func (i *Interpreter) RequestReturn() error {
	// Pop the top call context
	callContext := i.popCallContext()
	if callContext == nil {
		return ErrReturnWithoutGosub
	}

	// Jump back to the return address
	i.pc = callContext.ReturnLineIndex
	i.jumped = true
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
		return ErrIllegalQuantity
	}
	// Jump back target is the next statement after the FOR statement on the same line
	return i.pushForLoop(variable, end, step, i.pc, i.stmtIndex+1)
}

// IterateFor performs a NEXT iteration; variable may be empty to use the most recent loop
func (i *Interpreter) IterateFor(variableName string) error {
	// Find the appropriate FOR loop context
	var forLoop *ForLoopContext
	if variableName != "" {
		// NEXT with variable name - find specific loop
		forLoop = i.findForLoopByVariable(variableName)
		if forLoop == nil {
			return ErrNextWithoutFor
		}
	} else {
		// NEXT without variable name - use most recent loop
		forLoop = i.peekForLoop()
		if forLoop == nil {
			return ErrNextWithoutFor
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
		// Update loop variable and jump back to the statement after FOR
		err = i.SetVariable(forLoop.Variable, newValue)
		if err != nil {
			return err
		}
		// Signal statement-level jump to AfterForLineIndex:AfterForStmtIndex
		i.pc = forLoop.AfterForLineIndex
		i.stmtIndex = forLoop.AfterForStmtIndex
		i.stmtJumped = true
		return nil
	}

	// Loop finished - pop the loop from stack and continue normally
	i.popForLoop()
	return nil
}

// Built-in function implementations

// evaluateLenFunction implements the LEN function
func (i *Interpreter) evaluateLenFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: LEN requires exactly 1 argument")
	}

	arg := args[0]
	if arg.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: LEN requires string argument")
	}

	return types.NewNumberValue(float64(len(arg.String))), nil
}

// evaluateLeftFunction implements the LEFT$ function
func (i *Interpreter) evaluateLeftFunction(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: LEFT$ requires exactly 2 arguments")
	}

	str := args[0]
	count := args[1]

	if str.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: LEFT$ first argument must be string")
	}
	if count.Type != types.NumberType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: LEFT$ second argument must be number")
	}

	// Handle negative or zero count
	if count.Number <= 0 {
		return types.NewStringValue(""), nil
	}

	// Convert count to integer and handle bounds
	countInt := int(count.Number)
	if countInt >= len(str.String) {
		return str, nil // Return entire string if count exceeds length
	}

	return types.NewStringValue(str.String[:countInt]), nil
}

// evaluateRightFunction implements the RIGHT$ function
func (i *Interpreter) evaluateRightFunction(args []types.Value) (types.Value, error) {
	if len(args) != 2 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: RIGHT$ requires exactly 2 arguments")
	}

	str := args[0]
	count := args[1]

	if str.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: RIGHT$ first argument must be string")
	}
	if count.Type != types.NumberType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: RIGHT$ second argument must be number")
	}

	// Handle negative or zero count
	if count.Number <= 0 {
		return types.NewStringValue(""), nil
	}

	// Convert count to integer and handle bounds
	countInt := int(count.Number)
	if countInt >= len(str.String) {
		return str, nil // Return entire string if count exceeds length
	}

	// Extract rightmost characters
	startPos := len(str.String) - countInt
	return types.NewStringValue(str.String[startPos:]), nil
}

// evaluateMidFunction implements the MID$ function
func (i *Interpreter) evaluateMidFunction(args []types.Value) (types.Value, error) {
	if len(args) != 3 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: MID$ requires exactly 3 arguments")
	}

	src := args[0]
	start := args[1]
	length := args[2]

	if src.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: MID$ first argument must be string")
	}
	if start.Type != types.NumberType || length.Type != types.NumberType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: MID$ second and third arguments must be numbers")
	}

	if len(src.String) == 0 {
		return types.NewStringValue(""), nil
	}

	// 1-based start position
	startInt := int(start.Number)
	countInt := int(length.Number)

	if countInt <= 0 {
		return types.NewStringValue(""), nil
	}

	// Convert to 0-based index
	idx := startInt - 1
	if idx < 0 {
		// Be forgiving: treat before-start as empty result
		return types.NewStringValue(""), nil
	}
	if idx >= len(src.String) {
		return types.NewStringValue(""), nil
	}

	end := idx + countInt
	if end > len(src.String) {
		end = len(src.String)
	}
	if end <= idx {
		return types.NewStringValue(""), nil
	}
	return types.NewStringValue(src.String[idx:end]), nil
}

// evaluateChrFunction implements the CHR$ function
func (i *Interpreter) evaluateChrFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: CHR$ requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.NumberType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: CHR$ requires numeric argument")
	}
	code := int(arg.Number)
	// Normalize to 0..255 range
	if code < 0 {
		code = 256 - ((-code) % 256)
	}
	code = code % 256
	b := byte(code)
	return types.NewStringValue(string([]byte{b})), nil
}

// evaluateAscFunction implements the ASC function
func (i *Interpreter) evaluateAscFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: ASC requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: ASC requires string argument")
	}
	if len(arg.String) == 0 {
		return types.NewNumberValue(0), nil
	}
	// Use first byte of UTF-8 representation for compatibility with simple ASCII
	c := arg.String[0]
	return types.NewNumberValue(float64(int(c))), nil
}

// evaluateStrFunction implements the STR$ function
func (i *Interpreter) evaluateStrFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: STR$ requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.NumberType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: STR$ requires numeric argument")
	}
	return types.NewStringValue(arg.ToString()), nil
}

// evaluateValFunction implements the VAL function
func (i *Interpreter) evaluateValFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: VAL requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.StringType {
		return types.Value{}, fmt.Errorf("?TYPE MISMATCH ERROR: VAL requires string argument")
	}
	s := strings.TrimSpace(arg.String)
	if s == "" {
		return types.NewNumberValue(0), nil
	}
	// Try to parse as float; if it fails, return 0 as C64 VAL behavior
	if v, err := types.ParseValue(s); err == nil && v.Type == types.NumberType {
		return v, nil
	}
	return types.NewNumberValue(0), nil
}

// evaluateRndFunction implements the RND function
// For now, it returns a pseudo-random number in [0,1).
// The argument is required (C64 style) but only used for compatibility.
func (i *Interpreter) evaluateRndFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: RND requires exactly 1 argument")
	}
	if args[0].Type != types.NumberType {
		return types.Value{}, types.ErrTypeMismatch
	}
	return types.NewNumberValue(i.runtime.Random()), nil
}

// evaluateAbsFunction implements the ABS function
func (i *Interpreter) evaluateAbsFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: ABS requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.NumberType {
		return types.Value{}, types.ErrTypeMismatch
	}
	v := arg.Number
	if v < 0 {
		v = -v
	}
	return types.NewNumberValue(v), nil
}

// evaluateIntFunction implements the INT function (floor)
func (i *Interpreter) evaluateIntFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: INT requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.NumberType {
		return types.Value{}, types.ErrTypeMismatch
	}
	// INT returns the greatest integer less than or equal to arg
	// Use math.Floor semantics
	return types.NewNumberValue(math.Floor(arg.Number)), nil
}

// evaluateSqrFunction implements the SQR function
func (i *Interpreter) evaluateSqrFunction(args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return types.Value{}, fmt.Errorf("?SYNTAX ERROR: SQR requires exactly 1 argument")
	}
	arg := args[0]
	if arg.Type != types.NumberType {
		return types.Value{}, types.ErrTypeMismatch
	}
	if arg.Number < 0 {
		return types.Value{}, ErrIllegalQuantity
	}
	return types.NewNumberValue(math.Sqrt(arg.Number)), nil
}
