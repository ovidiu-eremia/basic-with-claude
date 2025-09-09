# BASIC Interpreter Architecture

This interpreter uses a traditional multi-phase architecture with polymorphic execution: lexical analysis → parsing to AST → polymorphic tree-walking execution using double dispatch pattern.

## Core Components

### 1. Lexer (Tokenizer)
- Line-by-line tokenization matching BASIC's line-oriented nature
- Case-insensitive keyword recognition
- Token types for numbers, strings, identifiers, operators, keywords
- No per-token line numbers (parser tracks source lines)

### 2. Parser
- Recursive descent parser
- Input: Stream of tokens from lexer (one line at a time)
- Output: AST nodes representing program statements with polymorphic execution methods
- Defines `InterpreterOperations` interface for double dispatch

### 3. Abstract Syntax Tree (AST) with Polymorphic Execution

#### Double Dispatch Pattern
1. **First Dispatch**: Interpreter calls `stmt.Execute(interpreter)` - polymorphic method dispatch
2. **Second Dispatch**: AST node calls back to `interpreter.GetVariable()`, `interpreter.PrintLine()`, etc. - interface method dispatch
3. **Result**: Clean separation without circular dependencies

#### Core Types (defined in `parser/ast.go`)
- `Statement`: AST nodes that execute themselves via `Execute(ops InterpreterOperations)`
- `Expression`: AST nodes that evaluate to values via `Evaluate(ops InterpreterOperations)`
- `InterpreterOperations`: Interface enabling double dispatch from AST nodes back to interpreter

#### Node Types
- **Program**: Root node containing all program lines
- **Line**: Container for one or more statements on a BASIC line (includes the BASIC line number)
- **Statement Nodes**: PrintStatement, GotoStatement, GosubStatement, ReturnStatement, IfStatement, ForStatement, NextStatement, InputStatement, LetStatement, DimStatement, DataStatement, ReadStatement, RestoreStatement, RemStatement, EndStatement, StopStatement, RunStatement
- **Expression Nodes**: BinaryOperation, UnaryOperation, NumberLiteral, StringLiteral, VariableReference, ArrayReference, FunctionCall

### 4. Types Package
- Shared type system for values and operations
- `Value` type with `NumberType` and `StringType` variants
- Arithmetic operations (`Add`, `Subtract`, `Multiply`, `Divide`, `Power`)
- Dependencies: None (foundation package)

### 5. Interpreter
- Orchestrate program execution using polymorphic dispatch
- Implement `InterpreterOperations` interface for AST nodes
- Maintain execution state (program counter, stacks)
- Handle control flow via interpreter-managed state

#### Core Components
- **Program Counter**: Points to current line being executed
- **Line Number Index**: Map from BASIC line numbers to `Line` nodes for fast GOTO/GOSUB lookup
- **Variable Storage**: Uses `types.Value` for unified storage
- **Call Stack**: Stack of return addresses for GOSUB/RETURN
- **Loop Stack**: Stores active FOR loop contexts
- **Data Management**: Global list of DATA values with pointer for READ/RESTORE

#### Execution Model
```
1. Start at first line of program
2. While program counter is valid:
   a. Get current AST node
   b. Call stmt.Execute(interpreter) - polymorphic dispatch
   c. Apply interpreter state changes (jumps, halts) after statement execution
   d. Advance program counter (unless control flow occurred)
3. End when END/STOP requested, program counter goes past last line, or runtime error occurs
```

### 6. Runtime Environment
- Provide abstraction for all I/O operations
- Enable testing by allowing mock implementations
- Interface methods: `Print()`, `PrintLine()`, `Input()`, `Clear()`
- Implementations: StandardRuntime (production), TestRuntime (testing)

### 7. Error Handling
- Return errors from all execution methods (idiomatic Go)
- Format errors in C64 BASIC style: "?ERROR_TYPE ERROR IN LINE_NUMBER"
- Parse errors use source line numbers; runtime errors use BASIC line numbers

## Package Dependencies
```
types/           (foundation - no dependencies)
    ↑
parser/          (AST nodes + InterpreterOperations interface)
    ↑
interpreter/     (implements InterpreterOperations)
    ↑
cmd/basic        (main application)
```

## Key Design Decisions

1. **Polymorphic AST execution**: Each node executes itself using double dispatch
2. **Double dispatch pattern**: Eliminates switch statements while maintaining clean separation
3. **InterpreterOperations interface**: Enables AST nodes to call back to interpreter without circular dependencies
4. **Control flow via interpreter state**: GOTO/END/STOP and FOR/NEXT use ops that adjust pc/flags/stacks
5. **Unified Value type**: `types.Value` handles both numeric and string values with type safety
6. **Runtime environment interface**: Enables testing and I/O abstraction

## Architectural Benefits

**Polymorphic Design:**
- Each AST node contains its own execution logic
- New statement types require no interpreter changes
- AST nodes can be unit tested with mock operations
- Eliminates large switch statements and code duplication

**Double Dispatch Pattern:**
- Clean separation: AST defines behavior, interpreter provides operations
- No circular dependencies: Interface breaks dependency cycles
- Type safety: Compile-time verification of required operations