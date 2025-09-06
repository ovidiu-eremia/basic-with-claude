# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Philosophy

### Core Philosophy
Build a working BASIC interpreter incrementally, where each milestone delivers a complete, tested interpreter that runs a subset of the BASIC language. Each step builds upon the previous one, gradually expanding the language support.

### Incremental Development Approach
- **Breadth-first with complete features**: Add simple versions of features that work completely, rather than partial implementations of complex features
- **Start minimal**: First milestone is simplest possible working interpreter  
- **Build confidence**: Working interpreter at each step, always have something that runs
- **Maintain velocity**: Small, achievable milestones with quick feedback loops
- **Avoid big bangs**: No milestone requires major refactoring, architecture grows organically

## TDD Workflow

### Implementation Process (STRICTLY FOLLOW THIS)
1. Write a failing test that defines a desired function or improvement
2. Write just enough code to make the code compile
3. Run the test to confirm it fails as expected (not due to compilation errors)
4. Write the simplest code to make the test pass
5. Run the test to confirm success
6. **Use test failures as learning opportunities**: When tests fail unexpectedly, investigate the root cause rather than just fixing the assertion
7. Refactor code to improve design while keeping tests green
8. **Consider test organization**: As test suites grow, refactor test code for maintainability too
9. Repeat the cycle for each new feature or bugfix

## Milestone Structure
Each milestone follows this pattern:
1. Extend lexer to recognize new tokens
2. Extend parser to build AST nodes  
3. Extend interpreter to execute new nodes
4. Write acceptance tests first (TDD)
5. Implement until tests pass
6. Add error handling tests

## Definition of Done per Milestone
- [ ] Acceptance tests pass
- [ ] Unit tests for new code
- [ ] Error cases tested
- [ ] Previous tests still pass
- [ ] Code reviewed/refactored
- [ ] Methods placed with their data (no interpreter bloat)
- [ ] Documentation updated

## Testing Strategy

### Test Types
- **Unit Tests**: Lexer token generation, parser AST creation, expression evaluation, statement execution
- **Integration Tests**: Complete BASIC programs, control flow scenarios, error conditions, C64 compatibility
- **Test Programs**: Classic BASIC examples (Hello World, loops), edge cases (nested loops, deep recursion), error triggering programs

### Test Organization Patterns
- **Tabular tests for similar scenarios**: Group related test cases in table-driven tests with clear sections by functionality
- **Separate concerns**: Keep parsing tests, error tests, and execution tests in separate functions with distinct purposes
- **Mock interfaces for isolation**: Create focused unit tests using mock implementations of key interfaces (e.g., `InterpreterOperations`)
- **Consolidate when beneficial**: Refactor multiple similar test functions into comprehensive tabular tests for maintainability

### Polymorphic Testing
- **Test both dispatch levels**: For double dispatch patterns, test both the polymorphic method call and the callback operations
- **Mock interfaces for isolation**: Use mock implementations (e.g., `InterpreterOperations`) to isolate components under test
- **Error injection testing**: Verify error handling by injecting failures in mock dependencies to test edge cases
- **Behavioral verification**: Test that AST nodes call the correct interface methods with expected parameters

## Development Guidelines

### Core Principles
- **Single Responsibility**: Each package has one clear purpose
- **Interface Boundaries**: Clean interfaces between components  
- **Testability First**: Design for testing from the beginning
- **No Premature Abstraction**: Build what's needed for current milestone

### Implementation Rules
- **NEVER implement features beyond current step scope**
- **If demo files need unsupported features, update the demo file, don't add features**
- **NEVER use infinite loops** - always advance tokens/iterators
- **Parser infinite loops**: ensure `nextToken()` is called in all code paths
- **Avoid "manager bloat"** - Large switch statements often indicate missing method dispatch

### Go-Specific Patterns & Pitfalls

#### Naming Conventions
- Use descriptive names over abbreviations: `currentChar` not `ch`, `currentToken` not `curToken`
- Named return values for clarity: `(content string, terminated bool)` not `(string, bool)`
- Consistent field naming: avoid conflicts between fields and methods

#### Parser Safety
- **Parser loops**: always advance position in loops
- **Token advancement**: check both `currentToken` and `peekToken` progression
- **Error recovery**: skip to safe tokens (NEWLINE, EOF) when parsing fails

#### Code Quality
- **Extract helper methods** for repetitive operations (e.g., `createToken()`)
- **Consolidate similar parsing methods** with shared logic
- **Move constants to package level** to avoid recreation
- **Unified assignment parsing**: handle `LET A = 42` and `A = 42` with shared logic
- **Clear separation** between parsing and validation logic

#### Method Placement & Architecture
- **Operations on data** → Methods on types that own the data (`left.Compare(right, op)` not `compareValues(left, right, op)`)
- **Pure utilities** → Package-level functions (`NormalizeVariableName()` not receiver method)
- **Polymorphic behavior** → AST nodes execute themselves via double dispatch pattern
- **Refactoring signal**: Methods taking specific types as main parameters should move to those types

### Key Constraints & Patterns
- **Variable names**: 2 significant characters max (C64 compatible)
- **String variables**: End with $ (A$, NAME$) - lexer handles automatically
- **Line numbers**: 0-63999 range
- **Runtime interface**: All I/O must go through Runtime interface for testability
- **Variable storage**: Unified `types.Value` handles both numeric and string variables with type safety
- **String concatenation**: Uses `+` operator (not `&` like some BASIC variants)

### Critical Pitfalls to Avoid
- **No partial features**: Each milestone must work completely end-to-end
- **Test error cases**: Verify C64-compatible error messages ("?SYNTAX ERROR IN 10")
- **Preserve line numbers**: Carry line info through lexer → parser → interpreter
- **Never direct console I/O**: Always use Runtime interface

## Code Quality Workflow

### Pre-commit Integration
- **Embrace tool feedback**: Treat linter/formatter suggestions as learning opportunities
- **Fix root causes**: When tools flag issues, understand why rather than just fixing the immediate problem
- **Iterative improvement**: Use tool feedback to incrementally improve code quality
- **Never bypass quality checks**: Avoid using `--no-verify` when committing code

### Refactoring Guidelines
- **Apply to all code**: Refactor test code for maintainability, not just implementation code
- **Preserve coverage**: When refactoring tests, ensure identical coverage is maintained
- **Consistent patterns**: Use consistent patterns across similar test scenarios
- **Clean up proactively**: Remove unused functions, variables, and imports as flagged by linters

## Architecture Overview

This interpreter uses a traditional multi-phase architecture with polymorphic execution: lexical analysis → parsing to AST → polymorphic tree-walking execution using double dispatch pattern.

### 1. Lexer (Tokenizer)
- **Processing Model**: Line-by-line tokenization matching BASIC's line-oriented nature
- **Responsibility**: Convert raw text input into tokens
- **Key Features**:
  - Case-insensitive keyword recognition
  - Token types for numbers, strings, identifiers, operators, keywords
  - Preserve line number information with each token
  - Handle string literals with quotes
  - Recognize BASIC keywords (PRINT, GOTO, IF, etc.)
  - Distinguish between numeric and string variable names ($ suffix)

### 2. Parser
- **Type**: Recursive descent parser (well-understood, maintainable)
- **Input**: Stream of tokens from lexer (one line at a time)
- **Output**: AST nodes representing program statements with polymorphic execution methods
- **Expression Parsing**:
  - Recursive descent with operator precedence
  - Handle arithmetic, comparison, logical operators
  - Support function calls (built-in functions)
  - String concatenation and string expressions
- **Statement Parsing**:
  - One or more statements per line (colon-separated)
  - Each statement type has dedicated parsing logic
  - Maintain line number information in AST nodes
- **Polymorphic Interface**:
  - Defines `InterpreterOperations` interface for double dispatch
- Control flow for GOTO/END/STOP is stateful in the interpreter (no error types)
  - Enables AST nodes to execute themselves without interpreter switch statements

### 3. Abstract Syntax Tree (AST) with Polymorphic Execution

#### Core Structure (Double Dispatch Pattern)
```go
type Node interface {
    GetLineNumber() int
}

type Statement interface {
    Node
    Execute(ops InterpreterOperations) error
}

type Expression interface {
    Node
    Evaluate(ops InterpreterOperations) (types.Value, error)
}

type InterpreterOperations interface {
    GetVariable(name string) (types.Value, error)
    SetVariable(name string, value types.Value) error
    PrintLine(text string) error
    ReadInput(prompt string) (string, error)
    RequestGoto(targetLine int) error
    RequestEnd() error
    RequestStop() error
    NormalizeVariableName(name string) string
}
```

#### Node Types
- **Program**: Root node containing all program lines
- **Line**: Container for one or more statements on a line
- **Statement Nodes**: PrintStatement, GotoStatement, GosubStatement, ReturnStatement, IfStatement, ForStatement, NextStatement, InputStatement, LetStatement, DimStatement, DataStatement, ReadStatement, RestoreStatement, RemStatement, EndStatement, StopStatement, RunStatement
- **Expression Nodes**: BinaryOperation, UnaryOperation, NumberLiteral, StringLiteral, VariableReference, ArrayReference, FunctionCall

#### Polymorphic Execution Model (Double Dispatch Pattern)

**Double Dispatch Implementation:**
1. **First Dispatch**: Interpreter calls `stmt.Execute(interpreter)` - polymorphic method dispatch
2. **Second Dispatch**: AST node calls back to `interpreter.GetVariable()`, `interpreter.PrintLine()`, etc. - interface method dispatch
3. **Result**: Clean separation without circular dependencies

**Benefits:**
- **Self-Executing Nodes**: Each AST node knows how to execute itself
- **No Switch Statements**: Eliminates large switch statements in interpreter
- **Type Safety**: Interface ensures all required operations are available
- **Clean Separation**: AST defines behavior, interpreter provides operations
- **Testability**: Can mock `InterpreterOperations` for unit testing AST nodes
- **Extensibility**: Adding new statement types requires no interpreter changes

**Example Flow:**
```go
// Interpreter (first dispatch)
err := stmt.Execute(interpreter)

// PrintStatement.Execute (second dispatch)
func (ps *PrintStatement) Execute(ops InterpreterOperations) error {
    value, err := ps.Expression.Evaluate(ops)  // Recursive polymorphism
    if err != nil {
        return err
    }
    return ops.PrintLine(value.ToString())     // Callback to interpreter
}
```

#### Node Properties
- Each node stores its source line number for error reporting and GOTO targets
- Expression nodes evaluate themselves and return `types.Value`
- Statement nodes execute themselves using interpreter operations
- Control flow handled via interpreter state for GOTO/END/STOP; FOR/NEXT uses loop stack

### 4. Types Package
- **Responsibility**: Shared type system for values and operations
- **Key Components**:
  - `Value` type with `NumberType` and `StringType` variants
  - Arithmetic operations (`Add`, `Subtract`, `Multiply`, `Divide`, `Power`)
  - Comparison operations and type conversions
  - String/numeric value parsing and formatting
- **Dependencies**: None (foundation package)

### 5. Interpreter

#### Responsibility
- Orchestrate program execution using polymorphic dispatch
- Implement `InterpreterOperations` interface for AST nodes
- Maintain execution state (program counter, stacks)
- Coordinate between AST, runtime environment, and variable storage
- Handle control flow via interpreter-managed state (pc, flags, stacks)

#### Core Components

##### InterpreterOperations Implementation
- Implements interface defined in parser package
- Provides variable access, I/O operations, and control flow requests
- Enables double dispatch from AST nodes back to interpreter

##### Program Counter
- Points to current line being executed
- Modified directly by GOTO, loops, and returns via interpreter ops

##### Line Number Index
- Map from line numbers to AST nodes
- Enables fast lookup for GOTO/GOSUB targets (efficient GOTO/GOSUB)
- Built during parsing phase

##### Variable Storage
- Uses `types.Value` for unified storage
- Must support:
  - Numeric variables (floating point)
  - String variables (max 255 chars)
  - Numeric arrays
  - String arrays
  - 2-character significant variable names

##### Polymorphic Execution Loop
- Simple loop: `stmt.Execute(interpreter)` for each statement
- No switch statements - each AST node executes itself
- Control flow applied by interpreter state after each statement (jumps, halts)
- Unified error handling and line number reporting

##### Call Stack (GOSUB/RETURN)
- Stack of return addresses (AST nodes)
- Depth limit to prevent stack overflow
- RETURN pops address and jumps back

##### Loop Stack (FOR/NEXT)
- Stores active FOR loop contexts:
  - Normalized loop variable name
  - End value
  - Step value
  - After-FOR line index (jump target)
- NEXT updates variable and checks condition via ops
- Nested loops supported via stack (natural for nested structures)

##### Data Management
- Global list of all DATA values collected from program
- Global data pointer (index into data list)
- READ advances pointer
- RESTORE resets pointer (optionally to specific line)

#### Execution Model

##### Polymorphic Execution Loop
```
1. Start at first line of program
2. While program counter is valid:
   a. Get current AST node
   b. Call stmt.Execute(interpreter) - polymorphic dispatch
   c. Apply interpreter state changes (jumps, halts) after statement execution
   d. Advance program counter (unless control flow occurred)
3. End when:
   - END/STOP requested (halt flag set)
   - Program counter goes past last line
   - Runtime error occurs
```

##### Expression Evaluation (Polymorphic)
- Each expression node evaluates itself: `expr.Evaluate(interpreter)`
- Recursive evaluation of expression trees through double dispatch
- Type checking (numeric vs string) using `types.Value`
- Automatic numeric-to-string conversion where needed
- String-to-numeric conversion for VAL()

##### Control Flow (Interpreter State)
- **GOTO**: `ops.RequestGoto(line)` sets interpreter `pc`/jump state
- **GOSUB**: Push return address, jump to target (similar to GOTO)
- **RETURN**: Pop return address, jump back
- **IF...THEN**: Evaluate condition, conditionally execute THEN statement
- **FOR**: Initializes variable; calls `BeginFor(var, end, step)` to push loop context
- **NEXT**: Calls `IterateFor([var])` which updates variable and either jumps or exits loop
- **END/STOP**: Interpreter sets a halted flag via ops (no errors)

### 6. Runtime Environment

#### Responsibility
- Provide abstraction for all I/O operations
- Enable testing by allowing mock implementations (enables testing and I/O abstraction)
- Isolate system dependencies from core interpreter logic

#### Interface
```go
type Runtime interface {
    Print(value string) error
    PrintLine(value string) error
    Input(prompt string) (string, error)
    Clear() error
    // Additional I/O methods as needed
}
```

#### Implementations
- **StandardRuntime**: Production implementation using os.Stdout/Stdin
- **TestRuntime**: Mock implementation that buffers output and provides scripted input for testing

### 7. Error Handling

#### Philosophy & Approach
- **Implement proper errors from the start**: No panics or "not implemented"
- **Only parse what we support**: Parser rejects unsupported syntax with clear errors
- **Graceful degradation**: Unknown keywords produce syntax errors, not crashes
- **Return errors from all execution methods** (idiomatic Go error handling)
- **Propagate errors up through call chain**
- **Format errors in C64 BASIC style at top level**: "?ERROR_TYPE ERROR IN LINE_NUMBER"
- **Include line number in error messages**

#### Error Types
- **Syntax errors** (caught during parsing)
- **Runtime errors**:
  - Type mismatch
  - Undefined line number (GOTO/GOSUB)  
  - Division by zero
  - Out of data (READ)
  - Return without GOSUB
  - Next without FOR
  - Array bounds exceeded
  - String too long (>255 chars)

### 8. Built-in Functions
- **Function Registry**: Map of function names to implementation functions
- **Type Safety**: Each function performs type checking on arguments and returns appropriate type

### Data Flow (Polymorphic Architecture)
```
Source File (.bas)
    ↓
[Lexer] → Tokens (line by line)
    ↓
[Parser] → Self-Executing AST Nodes + InterpreterOperations Interface
    ↓
[Line Index Builder] → Line Number Map
    ↓
[Polymorphic Execution Loop]
    ↓
stmt.Execute(interpreter) ←→ [Interpreter as InterpreterOperations]
    ↑                              ├─ Program Counter
    │                              ├─ Variable Store (types.Value)
    │                              ├─ Call Stack
    │                              ├─ Loop Stack
    │                              └─ Data Pointer
    ↓                              ↓
[Control Flow State]          [Runtime Environment]
(GOTO/END/STOP via ops; FOR/NEXT via ops)
    ↓                         Output / Results
[Program Counter Updates]
```

### Package Dependencies
```
types/           (foundation - no dependencies)
    ↑
parser/          (AST nodes + InterpreterOperations interface)
    ↑
interpreter/     (implements InterpreterOperations)
    ↑
cmd/basic        (main application)
```

### Key Design Decisions

1. **Line-by-line lexing**: Matches BASIC's line-oriented nature
2. **Polymorphic AST execution**: Each node executes itself using double dispatch
3. **Double dispatch pattern**: Eliminates switch statements while maintaining clean separation
4. **InterpreterOperations interface**: Enables AST nodes to call back to interpreter without circular dependencies
5. **Control flow via interpreter state**: GOTO/END/STOP and FOR/NEXT use ops that adjust pc/flags/stacks
6. **Unified Value type**: `types.Value` handles both numeric and string values with type safety
7. **Simplified package structure**: Three packages (types, parser, interpreter) with linear dependencies
8. **Interface + structs for AST**: Flexible, type-safe Go pattern with polymorphic methods
9. **Recursive descent parsing**: Well-understood, maintainable
10. **Error returns (not panic)**: Idiomatic Go error handling
11. **Separate indices for line numbers**: Efficient GOTO/GOSUB
12. **Stack-based call/loop management**: Natural for nested structures
13. **Runtime environment interface**: Enables testing and I/O abstraction

### Architectural Benefits

**Polymorphic Design:**
- **Maintainability**: Each AST node contains its own execution logic
- **Extensibility**: New statement types require no interpreter changes
- **Testability**: AST nodes can be unit tested with mock operations
- **Code Quality**: Eliminates large switch statements and code duplication

**Double Dispatch Pattern:**
- **Clean Separation**: AST defines behavior, interpreter provides operations
- **No Circular Dependencies**: Interface breaks dependency cycles
- **Type Safety**: Compile-time verification of required operations
- **Performance**: No runtime type checking or reflection needed



## Development Commands

```bash
# Run tests and build
go test ./...
go build ./cmd/basic

# Run specific tests  
go test ./lexer
go test ./parser
go test ./interpreter
go test ./acceptance

# Run BASIC program
go run ./cmd/basic program.bas
```

## File Consultation Guide

- **Adding new BASIC statement**: Check `spec.md` for syntax, add AST node type from architecture section, implement Execute() method
- **Understanding error handling**: See "Error Handling" section in architecture above
- **Development approach**: See "Development Philosophy" and "Milestone Structure" sections above
- **Implementation guidance**: See "Development Guidelines" section above
- find test basic files in the testdata directory in the top folder
- feel free to write small throwaway test files for testing the interpreter
