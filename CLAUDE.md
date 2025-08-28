# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# MANDATORY: Read These Files FIRST
Before implementing ANY step, you MUST read these files in this order:
1. `spec.md` - Complete BASIC language specification
2. `design.md` - Architecture and component design
3. `implementation-strategy.md` - Development approach and milestones
4. `todo.md` - Current task list

## TDD Workflow (STRICTLY FOLLOW THIS)
1. Write failing test first
2. Write MINIMAL code to make test compile (but still fail for right reasons)
3. Verify test fails for expected reasons, not compilation errors
4. Implement just enough to make test pass
5. Refactor if needed
6. Repeat


## Implementation Rules
- NEVER implement features beyond current step scope
- If demo files need unsupported features, update the demo file, don't add features
- NEVER use infinite loops - always advance tokens/iterators
- Parser infinite loops: ensure `nextToken()` is called in all code paths

## Common Go Pitfalls to Avoid
- Parser loops: always advance position in loops
- Token advancement: check both `currentToken` and `peekToken` progression
- Error recovery: skip to safe tokens (NEWLINE, EOF) when parsing fails

## Code Quality Patterns (Learned from Critic Analysis)

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

## Current Implementation Status

**Project Status**: Step 5 complete - String variables fully implemented and tested
**Current Features**: Line parsing, PRINT statements, numeric variables, string variables, LET assignments
**Next Milestone**: Step 6 - Implement arithmetic expressions with proper operator precedence
**Code Quality**: Recently refactored based on critic analysis - reduced duplication, improved naming
**Testing**: Comprehensive test suite with unit, integration, and acceptance tests (all passing)
**Git Hooks**: Pre-commit hook runs full test suite automatically

## Essential Implementation Patterns

### Quick Reference Patterns
- **AST Node Interface**: See `design.md` for complete interface definition
- **Error Handling Rules**: See `design.md` for comprehensive error strategy  
- **Value System**: See `design.md` for Value types and variable storage

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