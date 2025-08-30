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

## TDD Workflow (STRICTLY FOLLOW THIS)
1. Write failing test first
2. Write MINIMAL code to make test compile (but still fail for right reasons)
3. Verify test fails for expected reasons, not compilation errors
4. Implement just enough to make test pass
5. Refactor if needed
6. Repeat

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
- [ ] Documentation updated

## Implementation Rules
- NEVER implement features beyond current step scope
- If demo files need unsupported features, update the demo file, don't add features
- NEVER use infinite loops - always advance tokens/iterators
- Parser infinite loops: ensure `nextToken()` is called in all code paths

## Error Handling Philosophy
- **Implement proper errors from the start**: No panics or "not implemented"
- **Only parse what we support**: Parser rejects unsupported syntax with clear errors
- **Graceful degradation**: Unknown keywords produce syntax errors, not crashes
- **C64-style error messages**: Format as "?ERROR_TYPE ERROR IN LINE_NUMBER"

## Implementation Principles
- **Single Responsibility**: Each package has one clear purpose
- **Interface Boundaries**: Clean interfaces between components
- **Testability First**: Design for testing from the beginning
- **No Premature Abstraction**: Build what's needed for current milestone

## Common Go Pitfalls to Avoid
- Parser loops: always advance position in loops
- Token advancement: check both `currentToken` and `peekToken` progression
- Error recovery: skip to safe tokens (NEWLINE, EOF) when parsing fails

## Code Quality Patterns

### Naming Conventions
- Use descriptive names over abbreviations: `currentChar` not `ch`, `currentToken` not `curToken`
- Named return values for clarity: `(content string, terminated bool)` not `(string, bool)`
- Consistent field naming: avoid conflicts between fields and methods

### Duplication Elimination
- Extract helper methods for repetitive operations (e.g., `createToken()`)
- Consolidate similar parsing methods with shared logic
- Move constants to package level to avoid recreation

### Parser Patterns
- Unified assignment parsing: handle `LET A = 42` and `A = 42` with shared logic
- Token creation helpers reduce boilerplate in lexer
- Clear separation between parsing and validation logic

## Project Overview

This is a BASIC interpreter written in Go implementing Commodore 64 BASIC V2 subset. Uses lexer → parser → AST → tree-walking interpreter architecture.

## Complete Architecture Documentation

This interpreter uses a traditional multi-phase architecture: lexical analysis → parsing to AST → direct tree-walking execution.

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
- **Output**: AST nodes representing program statements
- **Expression Parsing**:
  - Recursive descent with operator precedence
  - Handle arithmetic, comparison, logical operators
  - Support function calls (built-in functions)
  - String concatenation and string expressions
- **Statement Parsing**:
  - One or more statements per line (colon-separated)
  - Each statement type has dedicated parsing logic
  - Maintain line number information in AST nodes

### 3. Abstract Syntax Tree (AST)

#### Core Structure
```go
type ASTNode interface {
    Execute(interpreter *Interpreter) error
    GetLineNumber() int
}
```

#### Node Types
- **Program**: Root node containing all program lines
- **Line**: Container for one or more statements on a line
- **Statement Nodes**:
  - `PrintNode`: PRINT statement
  - `GotoNode`: GOTO statement
  - etc...
- **Expression Nodes**:
  - `BinaryOpNode`: Binary operations (+, -, *, /, ^, etc.)
  - `UnaryOpNode`: Unary operations (-, NOT)
  - etc...

#### Node Properties
- Each node stores its source line number for error reporting and GOTO targets
- Expression nodes return values when evaluated
- Statement nodes perform actions and control flow

### 4. Interpreter

#### Responsibility
- Orchestrate program execution
- Maintain execution state (program counter, stacks)
- Coordinate between AST, runtime environment, and variable storage
- Handle control flow and error propagation

#### Core Components

##### Program Counter
- Points to current AST node being executed
- Can be modified by GOTO, GOSUB, loops

##### Line Number Index
- Map from line numbers to AST nodes
- Enables fast lookup for GOTO/GOSUB targets (efficient GOTO/GOSUB)
- Built during parsing phase

##### Variable Storage
- Must support:
  - Numeric variables (floating point)
  - String variables (max 255 chars)
  - Numeric arrays
  - String arrays
  - 2-character significant variable names

##### Call Stack (GOSUB/RETURN)
- Stack of return addresses (AST nodes)
- Depth limit to prevent stack overflow
- RETURN pops address and jumps back

##### Loop Stack (FOR/NEXT)
- Stores active FOR loop contexts:
  - Loop variable name
  - End value
  - Step value
  - AST node of FOR statement
- NEXT statement updates variable and checks condition
- Nested loops supported via stack (natural for nested structures)

##### Data Management
- Global list of all DATA values collected from program
- Global data pointer (index into data list)
- READ advances pointer
- RESTORE resets pointer (optionally to specific line)

#### Execution Model

##### Main Execution Loop
```
1. Start at first line of program
2. While program counter is valid:
   a. Get current AST node
   b. Execute node
   c. Handle any control flow changes
   d. Advance program counter (unless modified)
3. End when:
   - END/STOP encountered
   - Program counter goes past last line
   - Runtime error occurs
```

##### Expression Evaluation
- Recursive evaluation of expression trees
- Type checking (numeric vs string)
- Automatic numeric-to-string conversion where needed
- String-to-numeric conversion for VAL()

##### Control Flow
- **GOTO**: Update program counter to target line
- **GOSUB**: Push return address, jump to target
- **RETURN**: Pop return address, jump back
- **IF...THEN**: Evaluate condition, conditionally execute
- **FOR**: Initialize loop variable, push loop context
- **NEXT**: Update variable, check condition, loop or exit

### 5. Runtime Environment

#### Responsibility
- Provide abstraction for all I/O operations
- Enable testing by allowing mock implementations (enables testing and I/O abstraction)
- Isolate system dependencies from core interpreter logic

#### Interface
```go
type RuntimeEnvironment interface {
    Print(value string) error
    PrintLine(value string) error
    Input(prompt string) (string, error)
    Clear() error
    // Additional I/O methods as needed
}
```

#### Implementations

##### StandardRuntime
- Real implementation for production use
- Uses os.Stdout for output
- Uses os.Stdin for input
- Direct console I/O

##### TestRuntime
- Mock implementation for testing
- Buffers output for assertion checking
- Provides pre-programmed input sequences
- Enables deterministic testing of interactive programs

#### Benefits
- Complete isolation of I/O from interpreter logic
- Testable PRINT and INPUT statements
- Can simulate complex user interaction scenarios
- Enables testing without actual console I/O

### 6. Error Handling Strategy

#### Approach
- Return errors from all execution methods (idiomatic Go error handling)
- Propagate errors up through call chain
- Format errors in C64 BASIC style at top level
- Include line number in error messages

#### Error Types
- Syntax errors (caught during parsing)
- Runtime errors:
  - Type mismatch
  - Undefined line number (GOTO/GOSUB)
  - Division by zero
  - Out of data (READ)
  - Return without GOSUB
  - Next without FOR
  - Array bounds exceeded
  - String too long (>255 chars)

### 7. Built-in Functions

#### Implementation
- Each function has dedicated evaluation logic
- Type checking on arguments
- Return appropriate type (numeric or string)

#### Function Registry
- Map of function names to implementation functions
- Called during expression evaluation

### Data Flow
```
Source File (.bas)
    ↓
[Lexer] → Tokens (line by line)
    ↓
[Parser] → AST Nodes
    ↓
[Line Index Builder] → Line Number Map
    ↓
[Interpreter]
    ├─ Program Counter
    ├─ Variable Store
    ├─ Call Stack
    ├─ Loop Stack
    └─ Data Pointer
    ↓
[Runtime Environment]
    ↓
Output / Results
```

### Key Design Decisions

1. **Line-by-line lexing**: Matches BASIC's line-oriented nature
2. **AST-based execution**: Clean separation of parsing and runtime
3. **Direct tree walking**: Simple, sufficient for BASIC performance needs
4. **Interface + structs for AST**: Flexible, type-safe Go pattern
5. **Recursive descent parsing**: Well-understood, maintainable
6. **Error returns (not panic)**: Idiomatic Go error handling
7. **Separate indices for line numbers**: Efficient GOTO/GOSUB
8. **Stack-based call/loop management**: Natural for nested structures
9. **Runtime environment interface**: Enables testing and I/O abstraction

### Testing Strategy

#### Unit Tests
- Lexer: Token generation for various inputs
- Parser: AST generation for each statement type
- Expression evaluation: Arithmetic, comparison, logical
- Individual statement execution

#### Integration Tests
- Complete BASIC programs
- Control flow scenarios
- Error conditions
- C64 BASIC compatibility tests

#### Test Programs
- Classic BASIC examples (Hello World, loops, etc.)
- Edge cases (nested loops, deep recursion)
- Error triggering programs

## Quick References

- **Language specification**: See `spec.md` for complete BASIC language features, operators, functions, and error types
- **Complete architecture**: See "Complete Architecture Documentation" section above for all component design details
- **Development approach**: See "Development Philosophy" and "Milestone Structure" sections above

## File Consultation Guide

- **Adding new BASIC statement**: Check `spec.md` for syntax, add AST node type from section above, implement Execute() method
- **Understanding error handling**: Go error returns, format as "?ERROR_TYPE ERROR IN LINE_NUMBER" (see Error Handling Strategy above)
- **Planning development approach**: See "Development Philosophy" and "Milestone Structure" sections above
- **Implementation patterns**: Use patterns below and complete architecture details above

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

## Essential Implementation Patterns

### Quick Reference Patterns
- **AST Node Interface**: All nodes implement Execute(interpreter) + GetLineNumber()
- **Error Handling Rules**: Return Go errors, format as C64 style with line numbers at top level
- **Value System**: Unified string storage for both numeric and string variables
- **Runtime Interface**: All I/O goes through RuntimeEnvironment interface for testability

### Key Constraints
- Variable names: 2 significant characters max (C64 compatible)
- String variables end with $ (A$, NAME$) - lexer handles this automatically
- Line numbers: 0-63999
- I/O must go through Runtime interface for testability
- Variable storage: Unified `map[string]string` works for both numeric and string variables

## Common Pitfalls

- **No partial features**: Each milestone must work completely end-to-end
- **Test error cases**: Verify C64-compatible error messages ("?SYNTAX ERROR IN 10")
- **Preserve line numbers**: Carry line info through lexer → parser → interpreter
- **Variable names**: Remember 2-character limit when implementing symbol table
- **String concatenation**: Uses `+` operator (not `&` like some BASIC variants)
- **Runtime interface**: All I/O must go through interface, never direct console access
