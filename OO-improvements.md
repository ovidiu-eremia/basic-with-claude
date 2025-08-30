# Object-Oriented Refactoring Analysis & Future Suggestions

Based on my analysis of the BASIC interpreter codebase, here are the key opportunities for more object-oriented improvements:

## **Completed OO Work (Good Foundation)**
- ✅ **Value type methods** - `IsTrue()`, `Compare()`, arithmetic operations
- ✅ **Error handling** - `RuntimeError` implements error interface
- ✅ **Interface usage** - Runtime environment abstraction

## **High-Impact OO Opportunities**

### 1. **AST Node Self-Execution (Visitor → Command Pattern)**
**Current**: Interpreter has giant switch statements and separate execute methods
**Future**: Each AST node should execute itself
```go
// Instead of: i.executePrintStatement(stmt)
// Have: stmt.Execute(context)
```

### 2. **Variable Management Abstraction** 
**Current**: `map[string]Value` + `normalizeVariableName()` scattered in interpreter
**Future**: `VariableScope` or `VariableManager` type
- Handles C64 name normalization
- Manages variable lookup/storage
- Supports future scoping (functions, subroutines)

### 3. **Program Counter as Object**
**Current**: Primitive `currentLineIndex` with manual jumping
**Future**: `ProgramCounter` type with methods
- `Jump(lineNumber)`, `Next()`, `Current()`
- Encapsulates line index lookups
- Supports call stack for GOSUB/RETURN

### 4. **Expression Evaluation Delegation**
**Current**: Interpreter evaluates all expression types
**Future**: Each expression evaluates itself
```go
// Instead of: i.evaluateExpression(expr)  
// Have: expr.Evaluate(context)
```

### 5. **Statement Factory Pattern**
**Current**: Manual AST construction in parser
**Future**: Statement factories for consistent creation

## **Medium-Impact Improvements**

### 6. **Error Context Management**
**Current**: Manual line number wrapping everywhere
**Future**: `ExecutionContext` with automatic error decoration

### 7. **Utility Functions as Methods**
**Current**: `normalizeVariableName()` as interpreter method
**Future**: Package-level utilities or dedicated types

### 8. **Type-Safe Operators**
**Current**: String-based operator dispatch
**Future**: Operator types with type-safe dispatch

## **Implementation Priority**
1. **AST Self-Execution** (biggest impact)
2. **Variable Management** (frequently used)
3. **Program Counter Object** (control flow clarity)
4. **Expression Self-Evaluation** (consistency)
5. **Remaining utilities** (cleanup)

This analysis identifies the most impactful areas where object-oriented principles can improve code organization, reusability, and maintainability in the BASIC interpreter.