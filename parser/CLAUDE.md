# Parser Package Development Guide

This directory contains the parser implementation for the BASIC interpreter, using a double dispatch pattern with polymorphic AST execution.

## Architecture Overview

### Core Interfaces

The parser defines these key types in `ast.go`:

- **`Statement`**: Executes via `Execute(ops InterpreterOperations)`
- **`Expression`**: Evaluates to values via `Evaluate(ops InterpreterOperations)`
- **`InterpreterOperations`**: Enables AST nodes to call back to interpreter without circular dependencies

### Double Dispatch Pattern

This pattern eliminates switch statements and enables clean separation:

1. **First Dispatch**: Interpreter calls `stmt.Execute(interpreter)` - polymorphic method dispatch
2. **Second Dispatch**: AST node calls `ops.GetVariable()`, `ops.PrintLine()`, etc. - interface method dispatch

**Benefits**: Self-executing nodes, no switch statements, clean separation, easy testing with mocks.

## Adding New AST Node Types

### 1. Define the Node Structure

Add to `ast.go`:

```go
// NewStatement represents a NEW statement
type NewStatement struct {
    SomeField Expression // Required fields
}
```

### 2. Implement Execute or Evaluate

**For statements:**
```go
func (ns *NewStatement) Execute(ops InterpreterOperations) error {
    // Use ops interface to interact with interpreter
    value, err := ns.SomeField.Evaluate(ops)
    if err != nil {
        return err
    }
    
    // Perform the statement's action
    return ops.SomeInterpreterOperation(value)
}
```

**For expressions:**
```go
func (ne *NewExpression) Evaluate(ops InterpreterOperations) (types.Value, error) {
    // Evaluate sub-expressions
    result, err := ns.SomeField.Evaluate(ops)
    if err != nil {
        return types.Value{}, err
    }
    
    // Return computed value
    return types.NewNumberValue(computedResult), nil
}
```

### 3. Add Parser Support

In `parser.go`, add parsing logic:

```go
func (p *Parser) parseNewStatement() Statement {
    stmt := &NewStatement{}
    // Parse required fields
    // Always advance tokens to avoid infinite loops
    return stmt
}
```

Update the statement parsing switch in `parseStatement()`.

### 4. Add InterpreterOperations Method (if needed)

If your node needs new interpreter operations, add to the `InterpreterOperations` interface in `ast.go`:

```go
type InterpreterOperations interface {
    // Existing methods...
    
    // New operation for your statement
    SomeInterpreterOperation(value types.Value) error
}
```

## Testing Patterns

### Unit Testing AST Nodes

Use the `MockInterpreterOperations` from `ast_test_helpers.go`:

```go
func TestNewStatement_Execute(t *testing.T) {
    tests := []struct {
        name        string
        field       Expression
        expected    string // or other expected result
        expectError bool
    }{
        {
            name:        "basic case",
            field:       &StringLiteral{Value: "TEST"},
            expected:    "TEST",
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := newMockOps()
            stmt := &NewStatement{SomeField: tt.field}

            err := stmt.Execute(mock)

            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                // Verify mock interactions
                assert.Equal(t, tt.expected, mock.getSomeResult())
            }
        })
    }
}
```

### Error Injection Testing

Test error handling using mock error injection:

```go
func TestNewStatement_Execute_ErrorCases(t *testing.T) {
    t.Run("field evaluation error", func(t *testing.T) {
        mock := newMockOps()
        mock.getVariableError = errors.New("variable error")

        stmt := &NewStatement{SomeField: &VariableReference{Name: "A"}}

        err := stmt.Execute(mock)
        assert.Error(t, err)
    })
}
```

### Test File Organization

Follow these naming patterns:

- `ast_<node_type>_test.go` - Tests for specific AST node types
- `ast_test_helpers.go` - Shared mock and helper functions
- `parser_test.go` - Parser integration tests
- `test_helpers.go` - General parser test utilities

## Testing Guidelines

### Tabular Tests
Use table-driven tests for similar scenarios:

```go
tests := []struct {
    name     string
    input    string
    expected interface{}
    wantErr  bool
}{
    // Test cases...
}
```

### Mock Usage Patterns

**Variable Setup:**
```go
mock := newMockOps()
mock.setVariable("A", types.NewNumberValue(42))
```

**Error Injection:**
```go
mock.getVariableError = errors.New("test error")
mock.printLineError = errors.New("print error")
```

**Output Verification:**
```go
output := mock.getOutput()
assert.Len(t, output, 1)
assert.Equal(t, "expected", output[0])
```

**Control Flow Verification:**
```go
assert.True(t, mock.isGotoRequested())
assert.Equal(t, 100, mock.getGotoTarget())
```

### Behavioral Testing

Test both the polymorphic dispatch and the interface callbacks:

1. **Test the node's Execute/Evaluate method directly** (unit test)
2. **Test that the node calls the correct interface methods** (behavioral verification)
3. **Test error propagation** from interface methods

## Common Pitfalls

### 1. Source Line Tracking
AST nodes do not store source line numbers. The parser tracks source lines during parsing for error reporting; runtime errors use BASIC line numbers from `Line` nodes.

### 2. Infinite Loops in Parsing
Ensure `nextToken()` is called in all parser code paths.

### 3. Type Safety
Use `types.Value` for all values and check types before operations.

### 4. Error Propagation
Always check and return errors from sub-expression evaluation.

### 5. Mock Inconsistency
Keep mock behavior consistent with real interpreter behavior.

## File Structure

```
parser/
├── ast.go                    # Core AST node definitions and interfaces
├── parser.go                 # Main parser implementation
├── precedence.go            # Operator precedence definitions
├── test_helpers.go          # General parser testing utilities
├── ast_test_helpers.go      # AST-specific mocks and helpers
├── parser_test.go           # Parser integration tests
├── ast_<node>_test.go       # Individual AST node unit tests
└── CLAUDE.md               # This guide
```

## Key Principles

1. **Double Dispatch**: AST nodes execute themselves via interpreter operations
2. **Interface Isolation**: Use `InterpreterOperations` to avoid circular dependencies
3. **Comprehensive Testing**: Test both success and error cases with mocks
4. **Type Safety**: Always use `types.Value` and check types
5. **Error Handling**: Return errors, don't panic
6. **Accurate Line Reporting**: Parser maintains source line counters; interpreter uses BASIC line numbers

This architecture enables clean, testable, and extensible parser development while maintaining the polymorphic execution model that eliminates complex switch statements.
