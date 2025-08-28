# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a BASIC interpreter written in Go implementing Commodore 64 BASIC V2 subset. Uses lexer → parser → AST → tree-walking interpreter architecture.

## Quick References

- **Language specification**: See `spec.md` for complete BASIC language features, operators, functions, and error types
- **Architecture details**: See `design.md` for detailed component design, AST structure, and execution model  
- **Implementation strategy**: See `implementation-strategy.md` for milestone approach, testing philosophy, and development workflow

## File Consultation Guide

- **Adding new BASIC statement**: Check `spec.md` for syntax, `design.md` for AST node design
- **Understanding error handling**: See `design.md` lines 156-175 for comprehensive strategy
- **Planning development approach**: Refer to `implementation-strategy.md` for milestone philosophy
- **Implementation patterns**: Use quick patterns below, reference `design.md` for full details

## Development Commands

### Basic Go Commands
```bash
# Initialize Go module (if not done)
go mod init basic-interpreter

# Build the interpreter
go build ./cmd/basic

# Run tests
go test ./...

# Run specific package tests
go test ./lexer
go test ./parser
go test ./interpreter

# Run acceptance tests
go test ./acceptance

# Run with race detection
go test -race ./...

# Build and run a BASIC program
go run ./cmd/basic program.bas
```

### TDD Development Workflow
1. Write acceptance test (.bas file with expected output)
2. Extend lexer/parser/interpreter incrementally
3. Run tests until green
4. Refactor with test safety net

## Current Implementation Status

**Project Status**: Planning phase - no Go code implemented yet
**Next Milestone**: Initialize Go module and implement minimal interpreter (PRINT, variables, RUN)

## Essential Implementation Patterns

### AST Node Interface (Quick Reference)
```go
type Node interface {
    Execute(interpreter *Interpreter) error
    GetLineNumber() int  // For error reporting and GOTO targets
}
```

### Error Handling Rules
- Always return errors, never panic
- Include line numbers: "?SYNTAX ERROR IN 10" 
- Only parse implemented syntax (reject unsupported with clear errors)

### Key Constraints
- Variable names: 2 significant characters max (C64 compatible)
- String variables end with $ (A$, NAME$)
- Line numbers: 0-63999
- I/O must go through Runtime interface for testability

## Common Pitfalls

- **No partial features**: Each milestone must work completely end-to-end
- **Test error cases**: Verify C64-compatible error messages ("?SYNTAX ERROR IN 10")
- **Preserve line numbers**: Carry line info through lexer → parser → interpreter
- **Variable names**: Remember 2-character limit when implementing symbol table
- **String concatenation**: Uses `+` operator (not `&` like some BASIC variants)
- **Runtime interface**: All I/O must go through interface, never direct console access
- use stretchr/testify in all tests
- prefer tabular tests
- after writing a test, we want to make just enough production code to see the tests fail.  Let's make sure they fail for the right reason, not for compilation errors