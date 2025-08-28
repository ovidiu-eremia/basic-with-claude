# BASIC Interpreter Implementation Strategy

## Core Philosophy
Build a working BASIC interpreter incrementally, where each milestone delivers a complete, tested interpreter that runs a subset of the BASIC language. Each step builds upon the previous one, gradually expanding the language support.

## Project Structure

### Clean Separation from Start
```
basic-interpreter/
├── cmd/
│   └── basic/
│       └── main.go          # CLI entry point
├── lexer/
│   ├── lexer.go             # Tokenizer
│   └── lexer_test.go
├── parser/
│   ├── parser.go            # Parser
│   ├── ast.go               # AST definitions
│   └── parser_test.go
├── interpreter/
│   ├── interpreter.go       # Core interpreter
│   ├── variables.go         # Variable storage
│   ├── expressions.go       # Expression evaluation
│   └── interpreter_test.go
├── runtime/
│   ├── runtime.go           # Runtime interface
│   ├── standard.go          # Production I/O
│   └── test.go              # Test harness I/O
├── acceptance/
│   ├── harness.go           # Test harness
│   ├── testdata/            # .bas test files
│   └── acceptance_test.go   # Acceptance tests
└── go.mod
```

## Testing Strategy

### Test Harness Architecture
- **Hexagonal Architecture**: Core interpreter is isolated from I/O
- **Runtime Interface**: Abstracts all I/O operations
- **Test Runtime**: Captures output, provides scripted input
- **Acceptance Tests**: Run complete .bas programs with assertions

### Test Harness Implementation
```go
// runtime/runtime.go
type Runtime interface {
    Print(value string) error
    PrintLine(value string) error
    Input(prompt string) (string, error)
    Clear() error
    // Additional methods added as needed
}

// runtime/test.go
type TestRuntime struct {
    OutputBuffer []string
    InputQueue   []string
    InputIndex   int
}

// acceptance/harness.go
type TestHarness struct {
    interpreter *interpreter.Interpreter
    runtime     *runtime.TestRuntime
}

func (h *TestHarness) RunFile(path string) error
func (h *TestHarness) AssertOutput(expected []string) error
func (h *TestHarness) SetInput(inputs []string)
```

### Test Coverage Per Milestone
1. **Acceptance Tests**: Complete .bas programs that verify feature works end-to-end
2. **Unit Tests**: Individual component tests (lexer, parser, interpreter)
3. **Error Tests**: Verify proper error messages and handling
4. **Regression Tests**: All previous milestone tests continue passing

## Incremental Development Approach

### Milestone Structure
Each milestone follows this pattern:
1. Extend lexer to recognize new tokens
2. Extend parser to build AST nodes
3. Extend interpreter to execute new nodes
4. Write acceptance tests first (TDD)
5. Implement until tests pass
6. Add error handling tests

### Feature Growth Strategy
**Breadth-first with complete features**: Add simple versions of features that work completely, rather than partial implementations of complex features.

### Parser Evolution
- Start with minimal parser that only handles current milestone
- Grow organically with each feature
- Maintain clean structure that's easy to extend
- Each AST node type is self-contained

### Error Handling Philosophy
- **Implement proper errors from the start**: No panics or "not implemented"
- **Only parse what we support**: Parser rejects unsupported syntax with clear errors  
- **C64-style error messages**: See `design.md` for complete error handling strategy
- **Graceful degradation**: Unknown keywords produce syntax errors, not crashes

## Implementation Principles

### Code Organization
1. **Single Responsibility**: Each package has one clear purpose
2. **Interface Boundaries**: Clean interfaces between components
3. **Testability First**: Design for testing from the beginning
4. **No Premature Abstraction**: Build what's needed for current milestone

### AST Design
See `design.md` for complete AST structure and interface definitions.

### Value System
See `design.md` for detailed Value types and variable storage implementation.

### Execution Model
1. **Tree-walking interpreter**: Direct AST execution
2. **Line number index**: Built during parsing for GOTO/GOSUB
3. **Execution context**: Tracks program counter, call stack, variables
4. **Error propagation**: Errors bubble up with line number context

## Development Workflow

### Per-Milestone Process
1. **Write acceptance test**: Create .bas file and expected output
2. **Extend lexer**: Add new token types if needed
3. **Extend parser**: Add parsing for new statements/expressions
4. **Extend interpreter**: Implement execution logic
5. **Run acceptance test**: Verify feature works
6. **Add error tests**: Test error conditions
7. **Refactor**: Clean up while tests provide safety net

### Git Strategy
- **Feature branches**: One branch per milestone
- **Atomic commits**: Each commit should pass tests
- **Tag milestones**: Tag each completed milestone

## Quality Metrics

### Definition of Done per Milestone
- [ ] Acceptance tests pass
- [ ] Unit tests for new code
- [ ] Error cases tested
- [ ] Previous tests still pass
- [ ] Code reviewed/refactored
- [ ] Documentation updated

### Testing Metrics
- All acceptance tests from previous milestones pass
- New feature has >80% code coverage
- Error paths are tested
- No regression in existing features

## Technical Decisions

### Why Tree-Walking?
- Simple and sufficient for BASIC performance
- Easy to debug and understand
- Natural mapping to BASIC's structure
- No optimization needed initially

### Why Separate Lexer/Parser?
- Clean separation of concerns
- Easier to test in isolation
- Traditional, well-understood pattern
- Simplifies error reporting

### Why Runtime Interface?
- Enables comprehensive testing
- Allows future extensions (files, graphics)
- Clean separation from core logic
- Supports different execution environments

## Incremental Milestone Philosophy

### Start Minimal
- First milestone: Simplest possible working interpreter
- Each step adds one coherent feature
- Features are complete but simple

### Build Confidence
- Working interpreter at each step
- Always have something that runs
- Visible progress with each milestone

### Maintain Velocity
- Small, achievable milestones
- Quick feedback loops
- Refactor continuously with test safety net

### Avoid Big Bangs
- No milestone requires major refactoring
- Architecture grows organically
- Interfaces remain stable