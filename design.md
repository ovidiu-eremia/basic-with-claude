# BASIC Interpreter Design Document

## Overview
This document describes the design and architecture for a BASIC interpreter written in Go, implementing the specification defined in spec.md. The interpreter uses a traditional multi-phase architecture with lexical analysis, parsing to AST, and direct tree-walking execution.

## Architecture Components

### 1. Lexer (Tokenizer)
- **Processing Model**: Line-by-line tokenization
- **Responsibility**: Convert raw text input into tokens
- **Key Features**:
  - Case-insensitive keyword recognition
  - Token types for numbers, strings, identifiers, operators, keywords
  - Preserve line number information with each token
  - Handle string literals with quotes
  - Recognize BASIC keywords (PRINT, GOTO, IF, etc.)
  - Distinguish between numeric and string variable names ($ suffix)

### 2. Parser
- **Type**: Recursive descent parser
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

#### Node Types (Implementation Structs)
- **Program**: Root node containing all program lines
- **Line**: Container for one or more statements on a line
- **Statements**:
  - `PrintNode`: PRINT statement
  - `GotoNode`: GOTO statement
  - `GosubNode`: GOSUB statement
  - `ReturnNode`: RETURN statement
  - `IfNode`: IF...THEN statement
  - `ForNode`: FOR loop initialization
  - `NextNode`: NEXT statement
  - `InputNode`: INPUT statement
  - `LetNode`: Variable assignment
  - `DimNode`: Array declaration
  - `DataNode`: DATA statement
  - `ReadNode`: READ statement
  - `RestoreNode`: RESTORE statement
  - `RemNode`: REM (comment) statement
  - `EndNode`: END statement
  - `StopNode`: STOP statement
  - `RunNode`: RUN command
- **Expressions**:
  - `BinaryOpNode`: Binary operations (+, -, *, /, ^, etc.)
  - `UnaryOpNode`: Unary operations (-, NOT)
  - `NumberNode`: Numeric literal
  - `StringNode`: String literal
  - `VariableNode`: Variable reference
  - `ArrayRefNode`: Array element reference
  - `FunctionCallNode`: Built-in function calls

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
- Enables fast lookup for GOTO/GOSUB targets
- Built during parsing phase

##### Variable Storage
- Implementation detail left flexible
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
- Nested loops supported via stack

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

### 5. Error Handling

#### Strategy
- Return errors from all execution methods
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

### 6. Built-in Functions

#### Implementation
- Each function has dedicated evaluation logic
- Type checking on arguments
- Return appropriate type (numeric or string)

#### Function Registry
- Map of function names to implementation functions
- Called during expression evaluation

## Data Flow

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

### 5. Runtime Environment

#### Responsibility
- Provide abstraction for all I/O operations
- Enable testing by allowing mock implementations
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

## Key Design Decisions

1. **Line-by-line lexing**: Matches BASIC's line-oriented nature
2. **AST-based execution**: Clean separation of parsing and runtime
3. **Direct tree walking**: Simple, sufficient for BASIC performance needs
4. **Interface + structs for AST**: Flexible, type-safe Go pattern
5. **Recursive descent parsing**: Well-understood, maintainable
6. **Error returns (not panic)**: Idiomatic Go error handling
7. **Separate indices for line numbers**: Efficient GOTO/GOSUB
8. **Stack-based call/loop management**: Natural for nested structures
9. **Runtime environment interface**: Enables testing and I/O abstraction

## Testing Strategy

### Unit Tests
- Lexer: Token generation for various inputs
- Parser: AST generation for each statement type
- Expression evaluation: Arithmetic, comparison, logical
- Individual statement execution

### Integration Tests
- Complete BASIC programs
- Control flow scenarios
- Error conditions
- C64 BASIC compatibility tests

### Test Programs
- Classic BASIC examples (Hello World, loops, etc.)
- Edge cases (nested loops, deep recursion)
- Error triggering programs