# BASIC Interpreter Specification

## Overview
A BASIC language interpreter written in Go that implements a subset of Commodore 64 BASIC V2, focusing on core language features without hardware-specific operations.

## Program Format
- **Input**: Plain text files with one numbered line per line
- **Line Numbers**: Required, range 0-63999
- **Line Format**: `<line_number> <statement(s)>`
- **Multiple Statements**: Supported using colon (`:`) separator
- **Execution Mode**: Program mode only (run saved programs with RUN command)

## Data Types

### Numeric
- **Type**: Floating point numbers only
- **Variables**: Simple variable names (A, B, X1, etc.)
- **Variable Names**: 2 significant characters maximum

### Strings
- **Variables**: String variable names end with `$` (A$, B$, NAME$, etc.)
- **Maximum Length**: 255 characters
- **Concatenation**: Supported with `+` operator

## Arrays
- **Declaration**: `DIM` statement required before use
- **Syntax**: `DIM A(10)` declares array A with indices 0-10 (11 elements)
- **Types**: Both numeric and string arrays supported
- **Indexing**: 0-based but DIM specifies highest index (C64 convention)

## Commands and Statements

### Program Control
- `RUN` - Execute program from beginning
- `END` - End program execution
- `STOP` - Stop program execution

### Flow Control
- `GOTO <line_number>` - Jump to specified line
- `GOSUB <line_number>` - Call subroutine
- `RETURN` - Return from subroutine
- `IF <condition> THEN <statement>` - Conditional execution

### Loops
- `FOR <var> = <start> TO <end> [STEP <increment>]` - Begin for loop
- `NEXT [<var>]` - End for loop

### Input/Output
- `PRINT [<expression>][;|,]...` - Output to screen
- `INPUT [<prompt>;] <variable>` - Get user input

### Data Handling
- `READ <variable_list>` - Read from DATA statements
- `DATA <constant_list>` - Define data values
- `RESTORE [<line_number>]` - Reset DATA pointer
- `LET <variable> = <expression>` - Variable assignment (LET is optional)

### Other
- `REM <comment>` - Comment line (preserved in listing)
- `DIM <array>(size)[,...]` - Declare arrays

## Operators

### Arithmetic
- `+` - Addition
- `-` - Subtraction
- `*` - Multiplication
- `/` - Division
- `^` - Exponentiation

### Comparison
- `=` - Equal
- `<>` - Not equal
- `<` - Less than
- `>` - Greater than
- `<=` - Less than or equal
- `>=` - Greater than or equal

### Logical
- `AND` - Logical AND
- `OR` - Logical OR
- `NOT` - Logical NOT

## Functions

### String Functions
- `LEN(<string>)` - Return string length
- `LEFT$(<string>, <count>)` - Return leftmost characters
- `RIGHT$(<string>, <count>)` - Return rightmost characters
- `MID$(<string>, <start>, <count>)` - Return substring
- `CHR$(<code>)` - Convert ASCII code to character
- `ASC(<string>)` - Convert first character to ASCII code
- `STR$(<number>)` - Convert number to string
- `VAL(<string>)` - Convert string to number

### Numeric Functions
- `ABS(<number>)` - Absolute value
- `INT(<number>)` - Integer part
- `SQR(<number>)` - Square root
- `SIN(<number>)` - Sine
- `COS(<number>)` - Cosine
- `TAN(<number>)` - Tangent
- `ATN(<number>)` - Arctangent
- `EXP(<number>)` - Exponential
- `LOG(<number>)` - Natural logarithm
- `RND(<number>)` - Random number (0 to 1)

## Error Handling
- Display C64-style error messages (e.g., "?SYNTAX ERROR IN 10")
- Stop execution at error line
- Standard error types:
  - SYNTAX ERROR
  - TYPE MISMATCH
  - OVERFLOW
  - ILLEGAL QUANTITY
  - UNDEFINED STATEMENT
  - OUT OF DATA
  - RETURN WITHOUT GOSUB
  - NEXT WITHOUT FOR

## Constraints (C64 Compatible)
- **Line Numbers**: 0-63999
- **Variable Names**: 2 significant characters
- **String Length**: Maximum 255 characters
- **Array Dimensions**: As per C64 BASIC V2 limits

## Implementation Notes
1. Parser should handle traditional BASIC line number format
2. Tokenization not required (plain text interpretation)
3. Case-insensitive keywords
4. Variables are global scope
5. Implicit variable declaration (no DIM needed for simple variables)
6. Numeric variables initialized to 0, strings to empty string
7. Comments (REM statements) should be preserved during execution