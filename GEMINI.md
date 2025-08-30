# Gemini Project: BASIC Interpreter in Go

This document provides a comprehensive overview of the BASIC interpreter project, its structure, and development conventions to guide AI-assisted development.

## Project Overview

This project is a BASIC interpreter written in Go, designed to implement a subset of the Commodore 64 BASIC V2 language. The interpreter follows a classic architecture:

1.  **Lexer**: Converts the source code into a stream of tokens.
2.  **Parser**: Builds an Abstract Syntax Tree (AST) from the tokens.
3.  **Interpreter**: Executes the program by walking the AST.

The project is being developed incrementally, with each milestone adding new features and ensuring the interpreter remains in a working state. The development process is strictly test-driven.

## Building and Running

### Commands

-   **Run all tests:**
    ```bash
    go test ./...
    ```
-   **Build the interpreter:**
    ```bash
    go build ./cmd/basic
    ```
-   **Run a BASIC program:**
    ```bash
    go run ./cmd/basic <filename.bas>
    ```
-   **Run tests for a specific package:**
    ```bash
    go test ./lexer
    go test ./parser
    go test ./interpreter
    go test ./acceptance
    ```

## Development Conventions

### Workflow

The project follows a strict Test-Driven Development (TDD) workflow:

1.  Write a failing test first.
2.  Write the minimal code to make the test compile.
3.  Verify the test fails for the expected reason.
4.  Implement just enough code to make the test pass.
5.  Refactor if necessary.
6.  Repeat.

### Error Handling

-   Errors are returned from functions, not panicked.
-   Error messages are formatted in the C64 BASIC style (e.g., `?SYNTAX ERROR IN 10`).
-   Errors include the line number where they occurred.

### Naming Conventions

-   Use descriptive names (e.g., `currentToken` instead of `curToken`).
-   Follow standard Go naming conventions.

## Project Structure

The project is organized into the following packages:

-   `cmd/basic`: The main entry point for the command-line interpreter.
-   `lexer`: The tokenizer that converts source code into tokens.
-   `parser`: The parser that builds the AST. Contains `ast.go` for node definitions.
-   `interpreter`: The core interpreter that executes the AST.
-   `runtime`: An abstraction for I/O operations, allowing for testable `PRINT` and `INPUT` statements.
-   `acceptance`: End-to-end acceptance tests for the interpreter.

## Key Files for Context

-   `spec.md`: The complete language specification for the BASIC dialect being implemented.
-   `CLAUDE.md`: The complete architectural design, development philosophy, and project guidance for the interpreter.
-   `todo.md`: The current status of the project and the list of implemented and pending features.
