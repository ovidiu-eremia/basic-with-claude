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
- [ ] Methods placed with their data (no interpreter bloat)
- [ ] Documentation updated

## Testing Strategy

### Test Types
- **Unit Tests**: Lexer token generation, parser AST creation, expression evaluation, statement execution
- **Integration Tests**: Complete BASIC programs, control flow scenarios, error conditions, C64 compatibility
- **Test Programs**: Classic BASIC examples (Hello World, loops), edge cases (nested loops, deep recursion), error triggering programs

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
- **Place methods with their data** - Operations should live on the types they operate on
- **Extract pure functions** - Utility functions with no state dependencies go to package level
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

#### Method Placement
- **Value operations** → Value type methods (`left.Compare(right, op)` not `compareValues(left, right, op)`)
- **Pure utilities** → Package-level functions (`NormalizeVariableName()` not receiver method)
- **Type-specific logic** → Respective types (AST nodes execute themselves)
- **Refactoring signal**: Methods taking specific types as main parameters should move to those types

### Key Constraints & Patterns
- **AST Node Interface**: All nodes implement Execute(interpreter) + GetLineNumber()
- **Variable names**: 2 significant characters max (C64 compatible)
- **String variables**: End with $ (A$, NAME$) - lexer handles automatically
- **Line numbers**: 0-63999 range
- **Runtime interface**: All I/O must go through RuntimeEnvironment interface for testability
- **Variable storage**: Unified `map[string]string` works for both numeric and string variables
- **String concatenation**: Uses `+` operator (not `&` like some BASIC variants)

### Critical Pitfalls to Avoid
- **No partial features**: Each milestone must work completely end-to-end
- **Test error cases**: Verify C64-compatible error messages ("?SYNTAX ERROR IN 10")
- **Preserve line numbers**: Carry line info through lexer → parser → interpreter
- **Never direct console I/O**: Always use Runtime interface

## Architecture Overview

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
- **Statement Nodes**: PrintNode, GotoNode, GosubNode, ReturnNode, IfNode, ForNode, NextNode, InputNode, LetNode, DimNode, DataNode, ReadNode, RestoreNode, RemNode, EndNode, StopNode, RunNode
- **Expression Nodes**: BinaryOpNode, UnaryOpNode, NumberNode, StringNode, VariableNode, ArrayRefNode, FunctionCallNode

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
- **StandardRuntime**: Production implementation using os.Stdout/Stdin
- **TestRuntime**: Mock implementation that buffers output and provides scripted input for testing

### 6. Error Handling

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

### 7. Built-in Functions
- **Function Registry**: Map of function names to implementation functions
- **Type Safety**: Each function performs type checking on arguments and returns appropriate type

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