# BASIC Interpreter Implementation Plan

## Project Overview
Build a Commodore 64 BASIC V2 interpreter in Go using incremental development. Each step produces a working, testable interpreter that expands functionality progressively.

## Development Philosophy
- **No big bangs**: Each step builds incrementally on the previous
- **Always working**: Every step produces a runnable interpreter
- **UI testable**: Every feature can be verified through the command line interface
- **TDD approach**: Write tests first, implement to make them pass
- **Complete features**: No partial implementations - each step works end-to-end

## Implementation Steps

### Phase 1: Foundation (Steps 1-3)
Establish project structure and minimal working interpreter.

#### Step 1: Project Bootstrap
- [x] **Goal**: Set up Go project structure with basic CLI that can read files
- [x] **Demo**: `./basic hello.bas` reads file and shows "Program loaded" message
- [x] **Tests**: File reading, error handling for missing files

#### Step 2: Minimal Lexer + Parser
- [x] **Goal**: Parse simplest BASIC program (line numbers + PRINT literals)
- [x] **Demo**: Parse and display `10 PRINT "HELLO"` as structured data
- [x] **Tests**: Token recognition, line number parsing, string literal parsing

#### Step 3: Minimal Interpreter + Runtime
- [ ] **Goal**: Execute PRINT statements with string literals
- [ ] **Demo**: `10 PRINT "HELLO WORLD"` outputs "HELLO WORLD"
- [ ] **Tests**: AST execution, runtime interface, end-to-end program execution

### Phase 2: Variables and Assignment (Steps 4-6)
Add variable support and basic assignment.

#### Step 4: Numeric Variables
- [ ] **Goal**: Support numeric variable assignment and PRINT
- [ ] **Demo**: `10 LET A = 42` then `20 PRINT A` outputs "42"
- [ ] **Tests**: Variable storage, assignment parsing, variable evaluation

#### Step 5: String Variables
- [ ] **Goal**: Support string variable assignment and PRINT
- [ ] **Demo**: `10 LET A$ = "HELLO"` then `20 PRINT A$` outputs "HELLO"
- [ ] **Tests**: String variable handling, type distinction, C64 naming rules

#### Step 6: Arithmetic Expressions
- [ ] **Goal**: Support basic arithmetic in assignments and PRINT
- [ ] **Demo**: `10 PRINT 2 + 3 * 4` outputs "14" (precedence working)
- [ ] **Tests**: Expression parsing, operator precedence, arithmetic evaluation

### Phase 3: Control Flow (Steps 7-9)
Add basic program control and flow.

#### Step 7: RUN Command and Program Control
- [ ] **Goal**: Support RUN, END, STOP commands
- [ ] **Demo**: Multi-line program with END statement stops execution properly
- [ ] **Tests**: Program counter management, execution flow control

#### Step 8: GOTO Statement
- [ ] **Goal**: Support unconditional jumps
- [ ] **Demo**: Program with GOTO creates infinite loop (Ctrl+C to stop)
- [ ] **Tests**: Line number lookup, program counter modification

#### Step 9: IF...THEN Statement
- [ ] **Goal**: Support conditional execution
- [ ] **Demo**: `10 IF A > 5 THEN PRINT "BIG"` with different A values
- [ ] **Tests**: Conditional evaluation, comparison operators, conditional jumps

### Phase 4: User Interaction (Steps 10-11)
Add user input capabilities.

#### Step 10: INPUT Statement
- [ ] **Goal**: Support user input for variables
- [ ] **Demo**: `10 INPUT A` prompts user, stores value, can be PRINTed
- [ ] **Tests**: Input parsing, user interaction, variable assignment from input

#### Step 11: INPUT with Prompts
- [ ] **Goal**: Support INPUT with custom prompts
- [ ] **Demo**: `10 INPUT "ENTER YOUR NAME"; A$` shows custom prompt
- [ ] **Tests**: Prompt display, mixed input types (numeric/string)

### Phase 5: Loops (Steps 12-13)
Add iterative control structures.

#### Step 12: FOR...NEXT Loops
- [ ] **Goal**: Support basic FOR loops with default step
- [ ] **Demo**: `10 FOR I = 1 TO 5` with `30 NEXT I` counts 1 to 5
- [ ] **Tests**: Loop initialization, iteration, termination conditions

#### Step 13: FOR...NEXT with STEP
- [ ] **Goal**: Support custom step values in FOR loops
- [ ] **Demo**: `10 FOR I = 10 TO 1 STEP -2` counts down by 2
- [ ] **Tests**: Custom step values, negative steps, nested loops

### Phase 6: Subroutines (Steps 14-15)
Add subroutine capabilities.

#### Step 14: GOSUB...RETURN
- [ ] **Goal**: Support subroutine calls and returns
- [ ] **Demo**: Main program calls subroutine, subroutine returns to correct location
- [ ] **Tests**: Call stack management, return address tracking

#### Step 15: Nested GOSUB
- [ ] **Goal**: Support nested subroutine calls
- [ ] **Demo**: Subroutine calls another subroutine, both return correctly
- [ ] **Tests**: Multiple stack levels, proper stack unwinding

### Phase 7: Data Handling (Steps 16-18)
Add data storage and retrieval.

#### Step 16: DATA and READ Statements
- [ ] **Goal**: Support static data definition and reading
- [ ] **Demo**: DATA statements provide values for READ statements
- [ ] **Tests**: Data collection, sequential reading, end-of-data detection

#### Step 17: RESTORE Statement
- [ ] **Goal**: Support resetting data pointer
- [ ] **Demo**: READ data, RESTORE, READ same data again
- [ ] **Tests**: Data pointer management, RESTORE without line number

#### Step 18: RESTORE with Line Number
- [ ] **Goal**: Support RESTORE to specific line
- [ ] **Demo**: Multiple DATA blocks, RESTORE to specific one
- [ ] **Tests**: Line-specific data pointer reset, error handling

### Phase 8: String Functions (Steps 19-21)
Add string manipulation capabilities.

#### Step 19: Basic String Functions
- [ ] **Goal**: Support LEN, LEFT$, RIGHT$ functions
- [ ] **Demo**: String length and substring operations work correctly
- [ ] **Tests**: Function parsing, string manipulation, boundary conditions

#### Step 20: More String Functions
- [ ] **Goal**: Support MID$, CHR$, ASC functions
- [ ] **Demo**: Advanced string operations and ASCII conversions
- [ ] **Tests**: Character code conversions, substring extraction

#### Step 21: String Conversion Functions
- [ ] **Goal**: Support STR$, VAL functions
- [ ] **Demo**: Convert between numbers and strings
- [ ] **Tests**: Type conversions, format handling, error cases

### Phase 9: Math Functions (Steps 22-24)
Add mathematical capabilities.

#### Step 22: Basic Math Functions
- [ ] **Goal**: Support ABS, INT, SQR functions
- [ ] **Demo**: Basic mathematical operations work correctly
- [ ] **Tests**: Math function evaluation, negative number handling

#### Step 23: Trigonometric Functions
- [ ] **Goal**: Support SIN, COS, TAN, ATN functions
- [ ] **Demo**: Trigonometric calculations produce correct results
- [ ] **Tests**: Angle calculations, precision requirements

#### Step 24: Advanced Math Functions
- [ ] **Goal**: Support EXP, LOG, RND functions
- [ ] **Demo**: Exponential, logarithmic, and random operations
- [ ] **Tests**: Special cases (LOG of negative), random number generation

### Phase 10: Arrays (Steps 25-27)
Add array support.

#### Step 25: Array Declaration (DIM)
- [ ] **Goal**: Support array declaration with DIM
- [ ] **Demo**: `DIM A(10)` creates array, access with A(0) through A(10)
- [ ] **Tests**: Array allocation, index validation, memory management

#### Step 26: Array Assignment and Access
- [ ] **Goal**: Support array element assignment and retrieval
- [ ] **Demo**: Store values in array elements, retrieve and display them
- [ ] **Tests**: Array indexing, value storage, bounds checking

#### Step 27: String Arrays
- [ ] **Goal**: Support string arrays
- [ ] **Demo**: `DIM A$(5)` works with string storage and retrieval
- [ ] **Tests**: String array operations, mixed array types in programs

### Phase 11: Error Handling and Polish (Steps 28-30)
Add comprehensive error handling and final features.

#### Step 28: C64-Style Error Messages
- [ ] **Goal**: Implement all C64 BASIC error messages
- [ ] **Demo**: Various error conditions produce correct C64-style messages
- [ ] **Tests**: All error types, correct line number reporting

#### Step 29: Comments and REM
- [ ] **Goal**: Support REM statements (comments)
- [ ] **Demo**: Programs with REM statements execute correctly, comments preserved
- [ ] **Tests**: Comment parsing, preservation during execution

#### Step 30: Program Listing and Polish
- [ ] **Goal**: Add program listing capability and final polish
- [ ] **Demo**: Can load, list, and run complete BASIC programs
- [ ] **Tests**: Full C64 BASIC compatibility suite

---

## LLM Implementation Prompts

### Step 1 Prompt

```
I need to create the foundation for a BASIC interpreter project in Go. This is step 1 of a 30-step incremental development plan.

GOAL: Set up Go project structure with basic CLI that can read BASIC program files.

REQUIREMENTS:
1. Create proper Go module structure following the design from implementation-strategy.md
2. Implement a CLI in cmd/basic/main.go that accepts a .bas filename as argument
3. Read the file and display "Program loaded: <filename>" or appropriate error
4. Include proper error handling for file not found, permission errors, etc.
5. Set up basic project structure with placeholder packages (lexer, parser, interpreter, runtime)
6. Write tests for the file reading functionality

ACCEPTANCE CRITERIA:
- Running `go run ./cmd/basic hello.bas` reads the file and shows success message
- Running with non-existent file shows clear error message
- All tests pass with `go test ./...`

Please implement this following TDD principles - write tests first, then implement to make them pass. Focus on clean, simple code that follows Go conventions.
```

### Step 2 Prompt

```
Building on Step 1, I need to implement minimal lexer and parser functionality.

GOAL: Parse the simplest BASIC programs (line numbers + PRINT with string literals).

REQUIREMENTS:
1. Implement lexer/lexer.go to tokenize basic BASIC syntax:
   - Line numbers (integers)
   - PRINT keyword
   - String literals in quotes
   - Basic punctuation and whitespace
2. Implement parser/parser.go and parser/ast.go for minimal AST:
   - Program node (collection of lines)
   - Line node (line number + statements)
   - PrintStatement node for PRINT commands
   - StringLiteral node for quoted strings
3. Update CLI to parse loaded file and display parsed structure
4. Implement comprehensive tests for tokenization and parsing

EXAMPLE INPUT: `10 PRINT "HELLO WORLD"`
EXPECTED: Should tokenize and parse into proper AST structure

ACCEPTANCE CRITERIA:
- Can tokenize line numbers, PRINT, and string literals
- Parser builds correct AST from tokens
- CLI displays parsed structure (can be debug output for now)
- All tests pass, including error cases (unterminated strings, invalid syntax)

Continue following TDD - write failing tests first, implement to pass.
```

### Step 3 Prompt

```
Building on Steps 1-2, I need to add execution capability to create a working interpreter.

GOAL: Execute PRINT statements with string literals - first working BASIC interpreter!

REQUIREMENTS:
1. Implement interpreter/interpreter.go with basic execution engine:
   - Execute method that walks the AST
   - Support for PrintStatement execution
   - Program counter and line number tracking
2. Implement runtime/runtime.go interface and runtime/standard.go implementation:
   - Runtime interface with Print() and PrintLine() methods
   - StandardRuntime that outputs to console
   - TestRuntime for testing (captures output)
3. Update CLI to execute parsed programs instead of just displaying structure
4. Implement AST node Execute() methods following the interface pattern
5. Add comprehensive execution tests using TestRuntime

EXAMPLE PROGRAM:
```
10 PRINT "HELLO"
20 PRINT "WORLD"
```

ACCEPTANCE CRITERIA:
- Running `go run ./cmd/basic hello.bas` executes program and prints output
- PRINT statements output correctly with proper newlines
- TestRuntime allows testing output without console interaction
- All tests pass including execution tests
- Error handling for execution failures

This step creates the first working interpreter! Focus on clean separation between parsing and execution.
```

### Step 4 Prompt

```
Building on Steps 1-3, I need to add numeric variable support.

GOAL: Support numeric variable assignment and PRINT of variables.

REQUIREMENTS:
1. Extend lexer to recognize:
   - LET keyword (optional for assignment)
   - Variable names (letters/numbers, max 2 significant chars per C64)
   - Equals sign for assignment
   - Numeric literals (integers and floats)
2. Extend parser/ast.go with new nodes:
   - LetStatement (assignment)
   - VariableReference (for using variables)
   - NumberLiteral (numeric constants)
3. Extend interpreter with:
   - Variable storage (map of variable names to values)
   - Variable assignment execution
   - Variable retrieval for PRINT
4. Update runtime interface if needed for numeric output
5. Add comprehensive tests for variable operations

EXAMPLE PROGRAMS:
```
10 LET A = 42
20 PRINT A
```
```
10 X = 123
20 PRINT X
```

ACCEPTANCE CRITERIA:
- Can assign numeric values to variables (with and without LET)
- PRINT displays variable values correctly
- Variable names follow C64 rules (2 significant characters)
- Numeric literals are parsed and stored properly
- All tests pass including variable storage and retrieval

Focus on clean variable storage design and proper numeric handling.
```

### Step 5 Prompt

```
Building on Steps 1-4, I need to add string variable support.

GOAL: Support string variable assignment and PRINT of string variables.

REQUIREMENTS:
1. Extend lexer to recognize string variable names (ending with $)
2. Extend parser to handle string variable assignments and references
3. Extend interpreter variable storage to handle both numeric and string variables:
   - Separate storage or unified storage with type information
   - Proper type checking (don't assign numbers to string vars)
4. Implement C64-compatible string variable naming (A$, NAME$, etc.)
5. Add string length validation (max 255 characters per C64 spec)
6. Add comprehensive tests for string operations

EXAMPLE PROGRAMS:
```
10 LET A$ = "HELLO"
20 PRINT A$
```
```
10 NAME$ = "JOHN DOE"
20 PRINT NAME$
```

ACCEPTANCE CRITERIA:
- String variables (with $ suffix) work correctly
- Can assign string literals to string variables
- PRINT displays string variable values
- Type safety: cannot assign string to numeric variable or vice versa
- String length limits enforced (255 chars max)
- All tests pass including type checking and edge cases

Ensure clean separation between string and numeric variable handling.
```

### Step 6 Prompt

```
Building on Steps 1-5, I need to add arithmetic expression support.

GOAL: Support basic arithmetic expressions in assignments and PRINT statements.

REQUIREMENTS:
1. Extend lexer to recognize arithmetic operators: + - * / ^
2. Extend parser with expression parsing:
   - Implement operator precedence (^ highest, then */, then +-)
   - Support parentheses for grouping
   - BinaryOperation AST nodes
3. Extend interpreter with expression evaluation:
   - Recursive expression evaluation
   - Proper operator precedence handling
   - Support variables in expressions
4. Add comprehensive tests for arithmetic operations and precedence
5. Handle edge cases (division by zero, etc.)

EXAMPLE PROGRAMS:
```
10 PRINT 2 + 3 * 4
```
Expected output: 14 (precedence: 3*4=12, then 2+12=14)

```
10 A = 5
20 B = 3
30 PRINT A * B + 1
```
Expected output: 16

ACCEPTANCE CRITERIA:
- Arithmetic expressions evaluate with correct precedence
- Parentheses override default precedence
- Variables can be used in expressions
- PRINT can display calculated expression results
- Error handling for division by zero
- All tests pass including complex expressions

This step transforms the interpreter from a simple variable storage system to a real expression evaluator.
```

### Step 7 Prompt

```
Building on Steps 1-6, I need to add program control commands.

GOAL: Support RUN, END, and STOP commands for proper program flow control.

REQUIREMENTS:
1. Extend lexer/parser for RUN, END, STOP keywords
2. Modify interpreter execution model:
   - Program counter that can be controlled
   - RUN command starts execution from first line
   - END command terminates program cleanly
   - STOP command pauses/stops execution
3. Update CLI to handle RUN command (may need interactive mode or automatic RUN)
4. Add proper program termination handling
5. Add tests for program control flow

EXAMPLE PROGRAMS:
```
10 PRINT "START"
20 PRINT "MIDDLE" 
30 END
40 PRINT "NEVER REACHED"
```

```
10 PRINT "HELLO"
20 STOP
30 PRINT "ALSO NEVER REACHED"
```

ACCEPTANCE CRITERIA:
- Programs execute line by line from lowest to highest line number
- END statement stops execution cleanly
- STOP statement stops execution (may show "BREAK" message like C64)
- Lines after END/STOP are not executed
- RUN command works (either automatic on load or explicit)
- All tests pass including control flow verification

This step establishes proper program execution model essential for remaining features.
```

### Step 8 Prompt

```
Building on Steps 1-7, I need to add GOTO statement support.

GOAL: Support unconditional jumps to line numbers.

REQUIREMENTS:
1. Extend lexer/parser for GOTO keyword and line number targets
2. Modify interpreter execution:
   - Build line number index during parsing (map line numbers to AST nodes)
   - GOTO execution changes program counter to target line
   - Error handling for invalid line numbers
3. Add GOTO AST node and execution logic
4. Handle infinite loops gracefully (Ctrl+C should work)
5. Add comprehensive tests for jumps and error cases

EXAMPLE PROGRAMS:
```
10 PRINT "BEFORE JUMP"
20 GOTO 50
30 PRINT "SKIPPED"
40 PRINT "ALSO SKIPPED"
50 PRINT "AFTER JUMP"
```

```
10 PRINT "HELLO"
20 GOTO 10
```
(Creates infinite loop - should be stoppable)

ACCEPTANCE CRITERIA:
- GOTO jumps to correct line numbers
- Program counter updates correctly after GOTO
- Error message for undefined line numbers ("UNDEFINED STATEMENT" error)
- Infinite loops can be interrupted
- All tests pass including jump logic and error handling

This step adds the first non-linear program flow capability.
```

### Step 9 Prompt

```
Building on Steps 1-8, I need to add conditional execution.

GOAL: Support IF...THEN statements for conditional program flow.

REQUIREMENTS:
1. Extend lexer/parser for IF, THEN keywords and comparison operators (=, <>, <, >, <=, >=)
2. Add conditional expression evaluation:
   - Comparison operations between numbers
   - Comparison operations between strings
   - Boolean result evaluation
3. Add IF statement AST node:
   - Condition expression
   - THEN statement (can be any statement including GOTO)
4. Extend interpreter with conditional execution logic
5. Add comprehensive tests for all comparison operators and edge cases

EXAMPLE PROGRAMS:
```
10 A = 5
20 IF A > 3 THEN PRINT "BIG"
30 PRINT "DONE"
```

```
10 INPUT "ENTER NUMBER"; N
20 IF N = 0 THEN GOTO 50
30 PRINT "NOT ZERO"
40 GOTO 60
50 PRINT "IS ZERO"
60 END
```

ACCEPTANCE CRITERIA:
- IF conditions evaluate correctly for all comparison operators
- THEN clause executes only when condition is true
- Can use any statement after THEN (including GOTO)
- String and numeric comparisons work properly
- All tests pass including complex conditional logic

This step enables decision-making in BASIC programs.
```

### Step 10 Prompt

```
Building on Steps 1-9, I need to add user input capability.

GOAL: Support INPUT statement for getting user input into variables.

REQUIREMENTS:
1. Extend lexer/parser for INPUT keyword
2. Extend runtime interface with Input() method:
   - Update StandardRuntime to read from stdin
   - Update TestRuntime to use pre-programmed input queue
3. Add INPUT statement AST node:
   - Target variable for storing input
   - Input parsing and type conversion
4. Handle both numeric and string input appropriately
5. Add comprehensive tests using TestRuntime with programmed inputs

EXAMPLE PROGRAMS:
```
10 INPUT A
20 PRINT "YOU ENTERED: "; A
```

```
10 INPUT NAME$
20 PRINT "HELLO "; NAME$
```

ACCEPTANCE CRITERIA:
- INPUT prompts user and waits for input
- Numeric input is converted and stored in numeric variables
- String input is stored in string variables
- Type mismatches are handled gracefully
- TestRuntime allows automated testing of interactive programs
- All tests pass including input/output verification

This step makes programs interactive and user-responsive.
```

### Step 11 Prompt

```
Building on Steps 1-10, I need to enhance INPUT with custom prompts.

GOAL: Support INPUT with custom prompt strings.

REQUIREMENTS:
1. Extend parser for INPUT with optional prompt syntax: INPUT "prompt"; variable
2. Enhance INPUT statement execution:
   - Display custom prompt before waiting for input
   - Handle both prompted and non-prompted INPUT
3. Add proper formatting for prompts (space handling, question marks, etc.)
4. Ensure compatibility with C64 INPUT behavior
5. Add tests for both prompt styles

EXAMPLE PROGRAMS:
```
10 INPUT "WHAT IS YOUR NAME"; NAME$
20 PRINT "HELLO "; NAME$
```

```
10 INPUT "ENTER A NUMBER"; N
20 PRINT "DOUBLE IT IS: "; N * 2
```

ACCEPTANCE CRITERIA:
- INPUT with prompts displays the prompt correctly
- INPUT without prompts works as before (shows default "?" prompt)
- Prompts are formatted correctly (no extra spaces/punctuation)
- Both numeric and string prompted input work
- All tests pass including prompt formatting verification

This step improves user experience with descriptive input prompts.
```

### Step 12 Prompt

```
Building on Steps 1-11, I need to add basic loop support.

GOAL: Support FOR...NEXT loops with default step of 1.

REQUIREMENTS:
1. Extend lexer/parser for FOR, TO, NEXT keywords
2. Add loop management to interpreter:
   - FOR loop stack to track active loops
   - Loop variable initialization and management
   - Loop termination checking
3. Add FOR and NEXT AST nodes:
   - FOR: loop variable, start value, end value
   - NEXT: optional loop variable specification
4. Handle loop variable updates and condition checking
5. Add comprehensive tests for various loop scenarios

EXAMPLE PROGRAMS:
```
10 FOR I = 1 TO 5
20 PRINT I
30 NEXT I
```

```
10 FOR X = 0 TO 10
20 PRINT "VALUE: "; X
30 NEXT
```

ACCEPTANCE CRITERIA:
- FOR loops execute the correct number of times
- Loop variable is properly initialized and incremented
- Loop terminates when variable exceeds end value
- NEXT with and without variable name both work
- Nested loop support (basic level)
- All tests pass including loop boundary conditions

This step adds iterative capabilities to BASIC programs.
```

### Step 13 Prompt

```
Building on Steps 1-12, I need to add STEP support to FOR loops.

GOAL: Support custom step values in FOR...NEXT loops including negative steps.

REQUIREMENTS:
1. Extend parser for STEP keyword and step expressions
2. Enhance FOR loop execution:
   - Custom step values (positive, negative, fractional)
   - Proper loop termination logic for different step directions
   - Step value evaluation (can be expressions)
3. Update loop stack management for step values
4. Add comprehensive tests for various step scenarios including edge cases
5. Handle step value of 0 (should be error or infinite loop prevention)

EXAMPLE PROGRAMS:
```
10 FOR I = 10 TO 1 STEP -2
20 PRINT I
30 NEXT I
```

```
10 FOR X = 0 TO 100 STEP 25
20 PRINT "QUARTER: "; X
30 NEXT X
```

ACCEPTANCE CRITERIA:
- Positive step values work correctly
- Negative step values count downward properly
- Fractional step values work (if supported in design)
- Step of 0 is handled appropriately (error or prevention)
- Loop termination logic works for all step directions
- All tests pass including complex step scenarios

This step completes the FOR loop implementation with full C64 compatibility.
```

### Step 14 Prompt

```
Building on Steps 1-13, I need to add subroutine support.

GOAL: Support GOSUB...RETURN for subroutine calls.

REQUIREMENTS:
1. Extend lexer/parser for GOSUB and RETURN keywords
2. Add call stack management to interpreter:
   - Stack to track return addresses
   - GOSUB pushes current line + 1 to stack and jumps
   - RETURN pops address and jumps back
3. Add GOSUB and RETURN AST nodes
4. Handle stack overflow protection and errors
5. Add comprehensive tests for subroutine calls and returns

EXAMPLE PROGRAMS:
```
10 PRINT "MAIN START"
20 GOSUB 100
30 PRINT "MAIN END"
40 END
100 PRINT "IN SUBROUTINE"
110 RETURN
```

```
10 FOR I = 1 TO 3
20 GOSUB 200
30 NEXT I
40 END
200 PRINT "CALL #"; I
210 RETURN
```

ACCEPTANCE CRITERIA:
- GOSUB jumps to target line and saves return address
- RETURN jumps back to line after GOSUB call
- Multiple GOSUB calls work correctly
- Call stack prevents stack overflow
- Error handling for RETURN without GOSUB
- All tests pass including nested call scenarios

This step adds structured programming capabilities to BASIC.
```

### Step 15 Prompt

```
Building on Steps 1-14, I need to add nested subroutine support.

GOAL: Support nested GOSUB calls (subroutines calling other subroutines).

REQUIREMENTS:
1. Enhance call stack to support multiple levels of nesting
2. Add proper stack management:
   - Stack depth limits to prevent overflow
   - Proper unwinding on errors
3. Handle complex call patterns:
   - Subroutine calling another subroutine
   - Multiple return paths
   - Error recovery from deep call stacks
4. Add comprehensive tests for nested scenarios
5. Ensure proper error messages for stack issues

EXAMPLE PROGRAMS:
```
10 PRINT "MAIN"
20 GOSUB 100
30 PRINT "BACK IN MAIN"
40 END
100 PRINT "FIRST SUB"
110 GOSUB 200
120 PRINT "BACK IN FIRST"
130 RETURN
200 PRINT "SECOND SUB"
210 RETURN
```

ACCEPTANCE CRITERIA:
- Nested GOSUB calls work correctly (at least 3-4 levels deep)
- Each RETURN goes to the correct calling location
- Stack overflow is detected and reported
- Complex call patterns execute properly
- Error handling preserves program state
- All tests pass including deep nesting scenarios

This step completes the subroutine implementation with full nesting support.
```

### Step 16 Prompt

```
Building on Steps 1-15, I need to add data storage and retrieval.

GOAL: Support DATA and READ statements for static data in programs.

REQUIREMENTS:
1. Extend lexer/parser for DATA and READ keywords
2. Implement data management system:
   - Collect all DATA values during parsing
   - Global data pointer for READ operations
   - Support for both numeric and string data
3. Add DATA and READ AST nodes:
   - DATA: list of constants
   - READ: list of variables to fill
4. Handle end-of-data conditions and errors
5. Add comprehensive tests for data operations

EXAMPLE PROGRAMS:
```
10 READ A, B, C$
20 PRINT A; B; C$
30 DATA 10, 20, "HELLO"
```

```
10 FOR I = 1 TO 3
20 READ X
30 PRINT "VALUE: "; X
40 NEXT I
50 DATA 100, 200, 300
```

ACCEPTANCE CRITERIA:
- DATA statements store values for later retrieval
- READ statements fill variables with DATA values in order
- Mixed numeric and string data work correctly
- OUT OF DATA error when READ exceeds available data
- DATA can appear anywhere in program (usually at end)
- All tests pass including data exhaustion scenarios

This step adds static data handling capabilities to BASIC programs.
```

### Step 17 Prompt

```
Building on Steps 1-16, I need to add RESTORE functionality.

GOAL: Support RESTORE statement to reset data pointer.

REQUIREMENTS:
1. Extend parser for RESTORE keyword (without line number initially)
2. Enhance data management:
   - Reset data pointer to beginning of data list
   - Allow re-reading of DATA from start
3. Add RESTORE AST node and execution
4. Handle RESTORE with no data (should be safe operation)
5. Add tests for data pointer reset scenarios

EXAMPLE PROGRAMS:
```
10 READ A
20 PRINT "FIRST: "; A
30 RESTORE
40 READ B
50 PRINT "AGAIN: "; B
60 DATA 42
```

```
10 FOR I = 1 TO 2
20 FOR J = 1 TO 3
30 READ X
40 PRINT X;
50 NEXT J
60 PRINT
70 RESTORE
80 NEXT I
90 DATA 1, 2, 3
```

ACCEPTANCE CRITERIA:
- RESTORE resets data pointer to beginning
- Can re-read same DATA values after RESTORE
- RESTORE works even if no previous READ operations
- Multiple RESTORE operations work correctly
- All tests pass including repeated data access

This step adds data reusability to BASIC programs.
```

### Step 18 Prompt

```
Building on Steps 1-17, I need to add line-specific RESTORE.

GOAL: Support RESTORE with specific line numbers.

REQUIREMENTS:
1. Extend parser for RESTORE with optional line number
2. Enhance data management:
   - Track which DATA statements are on which lines
   - Allow setting data pointer to specific line's DATA
3. Handle errors for invalid line numbers or lines without DATA
4. Add comprehensive tests for line-specific restore operations
5. Maintain backward compatibility with parameterless RESTORE

EXAMPLE PROGRAMS:
```
10 DATA 1, 2, 3
20 READ A, B, C
30 PRINT A, B, C
40 RESTORE 60
50 READ X, Y
60 DATA 10, 20
70 PRINT X, Y
```

ACCEPTANCE CRITERIA:
- RESTORE without line number resets to beginning (existing behavior)
- RESTORE with line number sets pointer to that line's DATA
- Error handling for lines without DATA statements
- Error handling for non-existent line numbers
- Multiple DATA blocks can be accessed independently
- All tests pass including error conditions

This step completes the DATA/READ/RESTORE functionality with full C64 compatibility.
```

### Step 19 Prompt

```
Building on Steps 1-18, I need to add basic string functions.

GOAL: Support LEN, LEFT$, and RIGHT$ string functions.

REQUIREMENTS:
1. Extend lexer/parser for function names and function call syntax
2. Implement function evaluation system:
   - Function registry/dispatcher
   - Parameter parsing and validation
   - Type checking (string functions require string params)
3. Add string function implementations:
   - LEN(string) returns length
   - LEFT$(string, count) returns leftmost characters
   - RIGHT$(string, count) returns rightmost characters
4. Handle edge cases (empty strings, count > length, negative count)
5. Add comprehensive tests for all string functions

EXAMPLE PROGRAMS:
```
10 A$ = "HELLO WORLD"
20 PRINT LEN(A$)
30 PRINT LEFT$(A$, 5)
40 PRINT RIGHT$(A$, 5)
```

Expected output:
```
11
HELLO
WORLD
```

ACCEPTANCE CRITERIA:
- LEN function returns correct string length
- LEFT$ returns correct number of leftmost characters
- RIGHT$ returns correct number of rightmost characters
- Functions work in expressions and with variables
- Edge cases handled gracefully (empty strings, bounds)
- All tests pass including boundary conditions

This step begins adding BASIC's built-in function library.
```

### Step 20 Prompt

```
Building on Steps 1-19, I need to add more string functions.

GOAL: Support MID$, CHR$, and ASC string functions.

REQUIREMENTS:
1. Add three new string functions to the function registry:
   - MID$(string, start, count) extracts substring
   - CHR$(code) converts ASCII code to character
   - ASC(string) converts first character to ASCII code
2. Enhance function parameter handling:
   - Support functions with multiple parameters
   - Parameter count validation
   - Mixed parameter types (string + numeric)
3. Handle ASCII conversion edge cases (codes 0-255, empty strings)
4. Add comprehensive tests for new functions

EXAMPLE PROGRAMS:
```
10 A$ = "HELLO WORLD"
20 PRINT MID$(A$, 7, 5)
30 PRINT CHR$(65)
40 PRINT ASC("A")
```

Expected output:
```
WORLD
A
65
```

ACCEPTANCE CRITERIA:
- MID$ extracts correct substrings with start position and count
- CHR$ converts ASCII codes to characters (0-255 range)
- ASC converts first character of string to ASCII code
- Parameter validation prevents crashes on invalid inputs
- Functions work in complex expressions
- All tests pass including edge cases and error conditions

This step expands string manipulation capabilities significantly.
```

### Step 21 Prompt

```
Building on Steps 1-20, I need to add string conversion functions.

GOAL: Support STR$ and VAL functions for string/number conversion.

REQUIREMENTS:
1. Add conversion functions to function registry:
   - STR$(number) converts number to string representation
   - VAL(string) converts string to number (0 if invalid)
2. Handle conversion edge cases:
   - STR$ formatting (spaces, decimal places)
   - VAL parsing (leading spaces, partial numbers, invalid strings)
3. Ensure C64-compatible behavior:
   - STR$ includes leading space for positive numbers
   - VAL stops at first non-numeric character
4. Add comprehensive tests for conversion scenarios

EXAMPLE PROGRAMS:
```
10 A = 42
20 B$ = STR$(A)
30 PRINT "NUMBER AS STRING: '"; B$; "'"
40 C$ = "123.45XYZ"
50 D = VAL(C$)
60 PRINT "STRING AS NUMBER: "; D
```

Expected output:
```
NUMBER AS STRING: ' 42'
STRING AS NUMBER: 123.45
```

ACCEPTANCE CRITERIA:
- STR$ converts numbers to properly formatted strings
- VAL converts strings to numbers with C64-compatible parsing
- Leading spaces and signs are handled correctly
- Invalid strings convert to 0 with VAL
- Functions work in all expression contexts
- All tests pass including format verification

This step completes the essential string function library.
```

### Step 22 Prompt

```
Building on Steps 1-21, I need to add basic math functions.

GOAL: Support ABS, INT, and SQR mathematical functions.

REQUIREMENTS:
1. Add mathematical functions to function registry:
   - ABS(number) returns absolute value
   - INT(number) returns integer part (truncates toward zero)
   - SQR(number) returns square root
2. Handle mathematical edge cases:
   - SQR of negative numbers (error condition)
   - INT with very large numbers
   - Proper floating point handling
3. Ensure C64-compatible behavior for edge cases
4. Add comprehensive tests for all mathematical scenarios

EXAMPLE PROGRAMS:
```
10 PRINT ABS(-42)
20 PRINT INT(3.7)
30 PRINT INT(-3.7)
40 PRINT SQR(16)
50 PRINT SQR(2)
```

Expected output:
```
42
3
-3
4
1.41421356
```

ACCEPTANCE CRITERIA:
- ABS returns absolute value for positive and negative numbers
- INT truncates decimal part correctly (toward zero)
- SQR calculates square root with reasonable precision
- Error handling for SQR of negative numbers
- Functions work in complex mathematical expressions
- All tests pass including precision and error cases

This step begins building the mathematical function library.
```

### Step 23 Prompt

```
Building on Steps 1-22, I need to add trigonometric functions.

GOAL: Support SIN, COS, TAN, and ATN trigonometric functions.

REQUIREMENTS:
1. Add trigonometric functions to function registry:
   - SIN(angle) sine function (radians)
   - COS(angle) cosine function (radians)
   - TAN(angle) tangent function (radians)
   - ATN(angle) arctangent function (returns radians)
2. Handle trigonometric edge cases:
   - TAN at π/2, 3π/2, etc. (undefined points)
   - Precision considerations for common angles
   - Range handling for ATN
3. Ensure reasonable precision matching C64 behavior
4. Add comprehensive tests including special angle values

EXAMPLE PROGRAMS:
```
10 PI = 3.14159265
20 PRINT SIN(PI/2)
30 PRINT COS(0)
40 PRINT TAN(PI/4)
50 PRINT ATN(1)
```

Expected output (approximately):
```
1
1
1
0.785398163
```

ACCEPTANCE CRITERIA:
- SIN, COS return correct values for standard angles
- TAN handles undefined points gracefully
- ATN returns correct arctangent values
- Functions work with variables and expressions
- Precision is reasonable for BASIC programs
- All tests pass including special cases

This step adds trigonometric capabilities for mathematical programs.
```

### Step 24 Prompt

```
Building on Steps 1-23, I need to add advanced math functions.

GOAL: Support EXP, LOG, and RND mathematical functions.

REQUIREMENTS:
1. Add advanced mathematical functions:
   - EXP(number) exponential function (e^x)
   - LOG(number) natural logarithm
   - RND(number) random number generator (0 to 1)
2. Handle mathematical edge cases:
   - LOG of negative numbers or zero (error)
   - EXP overflow conditions
   - RND seeding behavior (C64-compatible if possible)
3. Implement proper random number generation:
   - RND(1) returns random 0-1
   - RND(0) returns last random number
   - RND with negative seeds random generator
4. Add comprehensive tests for all functions

EXAMPLE PROGRAMS:
```
10 PRINT EXP(1)
20 PRINT LOG(EXP(1))
30 FOR I = 1 TO 5
40 PRINT RND(1)
50 NEXT I
```

Expected output (EXP/LOG exact, RND will vary):
```
2.71828183
1
(5 random numbers between 0 and 1)
```

ACCEPTANCE CRITERIA:
- EXP calculates exponential with good precision
- LOG calculates natural logarithm correctly
- RND generates random numbers in 0-1 range
- Error handling for LOG of non-positive numbers
- Random number state management works properly
- All tests pass including error conditions

This step completes the mathematical function library.
```

### Step 25 Prompt

```
Building on Steps 1-24, I need to add array declaration support.

GOAL: Support DIM statement for declaring arrays.

REQUIREMENTS:
1. Extend lexer/parser for DIM keyword and array syntax:
   - DIM A(10) syntax
   - Multiple array declarations: DIM A(10), B$(5)
   - Array name validation (including string arrays with $)
2. Implement array storage system:
   - Dynamic array allocation
   - Support for both numeric and string arrays
   - C64-compatible indexing (DIM A(10) creates indices 0-10)
3. Add DIM statement AST node and execution
4. Handle array redimensioning (should be error like C64)
5. Add comprehensive tests for array declaration

EXAMPLE PROGRAMS:
```
10 DIM A(5)
20 DIM B$(10), C(3)
30 PRINT "ARRAYS DECLARED"
```

```
10 DIM NUMBERS(100)
20 PRINT "LARGE ARRAY READY"
```

ACCEPTANCE CRITERIA:
- DIM creates arrays with correct size (0 to specified index)
- Both numeric and string arrays can be declared
- Multiple arrays can be declared in one DIM statement
- Array redimensioning produces appropriate error
- Memory is allocated properly for declared arrays
- All tests pass including declaration verification

This step establishes the foundation for array operations.
```

### Step 26 Prompt

```
Building on Steps 1-25, I need to add array element access.

GOAL: Support array element assignment and retrieval.

REQUIREMENTS:
1. Extend parser for array indexing syntax:
   - Array assignment: A(5) = 42
   - Array access in expressions: PRINT A(5)
   - Array indices can be expressions: A(I+1) = 10
2. Implement array element operations:
   - Bounds checking (0 to declared size)
   - Proper storage and retrieval
   - Type safety (numeric vs string arrays)
3. Handle array access errors:
   - Subscript out of range
   - Using undeclared arrays
   - Type mismatches
4. Add comprehensive tests for array operations

EXAMPLE PROGRAMS:
```
10 DIM A(5)
20 FOR I = 0 TO 5
30 A(I) = I * 10
40 NEXT I
50 FOR I = 0 TO 5
60 PRINT "A("; I; ") = "; A(I)
70 NEXT I
```

```
10 DIM NAMES$(3)
20 NAMES$(0) = "ALICE"
30 NAMES$(1) = "BOB"
40 NAMES$(2) = "CHARLIE"
50 FOR I = 0 TO 2
60 PRINT NAMES$(I)
70 NEXT I
```

ACCEPTANCE CRITERIA:
- Array elements can be assigned and retrieved
- Array indices can be expressions (variables, calculations)
- Bounds checking prevents out-of-range access
- String and numeric arrays work independently
- Error messages match C64 behavior
- All tests pass including bounds and type checking

This step makes arrays fully functional for data storage.
```

### Step 27 Prompt

```
Building on Steps 1-26, I need to enhance array support.

GOAL: Complete array implementation with string arrays and mixed usage.

REQUIREMENTS:
1. Ensure robust string array operations:
   - String storage and retrieval in arrays
   - String length validation in array elements
   - Proper memory management for string arrays
2. Support complex array usage patterns:
   - Arrays in function calls: PRINT LEN(NAMES$(0))
   - Arrays in expressions: IF A(I) > 0 THEN...
   - Array elements as function parameters
3. Add comprehensive array testing:
   - Large arrays
   - Mixed array types in same program
   - Complex indexing expressions
4. Performance considerations for array access

EXAMPLE PROGRAMS:
```
10 DIM MATRIX(2,2): REM Note: may need 2D arrays or simulate with 1D
20 DIM WORDS$(10)
30 FOR I = 0 TO 9
40 WORDS$(I) = "WORD" + STR$(I)
50 NEXT I
60 FOR I = 0 TO 9
70 PRINT I; ": "; WORDS$(I); " ("; LEN(WORDS$(I)); " chars)"
80 NEXT I
```

ACCEPTANCE CRITERIA:
- String arrays handle variable-length strings properly
- Arrays work correctly in all expression contexts
- Memory management is efficient for large arrays
- Complex programs with multiple array types work
- Performance is reasonable for typical array sizes
- All tests pass including stress testing

This step completes the array implementation with full functionality.
```

### Step 28 Prompt

```
Building on Steps 1-27, I need to implement comprehensive error handling.

GOAL: Implement all C64-style error messages with proper formatting.

REQUIREMENTS:
1. Implement complete C64 BASIC error message system:
   - "?SYNTAX ERROR IN <line>" format
   - All error types from spec.md: SYNTAX ERROR, TYPE MISMATCH, OVERFLOW, etc.
   - Proper line number reporting in all error cases
2. Enhance error detection throughout the system:
   - Parser errors with specific line numbers
   - Runtime errors with execution context
   - Mathematical errors (division by zero, etc.)
   - Array and variable errors
3. Add error recovery and program state preservation
4. Comprehensive error testing for all error conditions

EXAMPLE ERROR CASES:
```
10 PRINT "HELLO
```
Should produce: ?SYNTAX ERROR IN 10

```
10 A$ = 42
```
Should produce: ?TYPE MISMATCH ERROR IN 10

ACCEPTANCE CRITERIA:
- All error messages match C64 format exactly
- Line numbers are reported correctly in all error cases  
- Error detection covers all language features implemented
- Programs stop at error line (don't continue execution)
- Error handling doesn't crash the interpreter
- All tests pass including comprehensive error case coverage

This step ensures robust error handling matching C64 behavior.
```

### Step 29 Prompt

```
Building on Steps 1-28, I need to add comment support and final language features.

GOAL: Support REM statements (comments) and any remaining basic features.

REQUIREMENTS:
1. Extend lexer/parser for REM keyword:
   - REM consumes rest of line as comment
   - Comments are preserved but don't execute
   - Handle REM in multi-statement lines (after colon)
2. Add any missing basic operators or syntax:
   - Logical operators (AND, OR, NOT)
   - String concatenation with +
   - Any missing comparison operators
3. Polish existing features:
   - Clean up any remaining parsing edge cases
   - Improve error messages
   - Performance optimizations if needed
4. Add comprehensive tests for comments and remaining features

EXAMPLE PROGRAMS:
```
10 REM This is a comment
20 PRINT "HELLO": REM Comment after statement
30 A = 5: REM Set A to 5
40 IF A > 0 AND A < 10 THEN PRINT "MIDDLE"
50 B$ = "HEL" + "LO"
60 PRINT B$
```

ACCEPTANCE CRITERIA:
- REM statements work correctly and don't execute
- Comments can appear after colons in multi-statement lines
- Logical operators (AND, OR, NOT) work in conditions
- String concatenation with + operator works
- All edge cases in parsing are handled gracefully
- All tests pass including comment preservation and logical operations

This step completes the core BASIC language implementation.
```

### Step 30 Prompt

```
Building on Steps 1-29, I need to add final polish and program listing capability.

GOAL: Complete the BASIC interpreter with listing capability and final polish.

REQUIREMENTS:
1. Add program listing functionality:
   - LIST command shows program lines
   - Optional line number ranges: LIST 10-50
   - Proper formatting matching C64 style
2. Final polish and optimization:
   - Code cleanup and organization
   - Performance improvements where beneficial
   - Memory usage optimization
3. Complete C64 BASIC compatibility testing:
   - Run comprehensive test suite
   - Test classic BASIC programs
   - Verify all features work together properly
4. Documentation and final testing

EXAMPLE USAGE:
```bash
$ ./basic
READY.
> LOAD "test.bas"
PROGRAM LOADED
> LIST
10 PRINT "HELLO WORLD"
20 FOR I = 1 TO 5
30 PRINT I
40 NEXT I
50 END
> RUN
HELLO WORLD
1
2
3
4
5
READY.
```

ACCEPTANCE CRITERIA:
- LIST command displays programs in proper format
- All 29 previous steps continue to work perfectly
- Can run substantial BASIC programs without issues
- Performance is acceptable for typical program sizes
- Memory usage is reasonable
- Complete test suite passes with 100% success
- Documentation is complete and accurate

This final step delivers a complete, production-ready BASIC interpreter that faithfully implements the Commodore 64 BASIC V2 subset as specified.

CONGRATULATIONS! You now have a fully functional BASIC interpreter built through careful, incremental development with comprehensive testing at each step.
```

---

## Implementation Notes

Each step should be implemented following these principles:

1. **Test-Driven Development**: Write failing tests first, implement to make them pass
2. **Incremental Progress**: Each step builds directly on previous steps
3. **Working Software**: Every step produces a runnable interpreter
4. **Comprehensive Testing**: Unit tests, integration tests, and acceptance tests
5. **Clean Code**: Follow Go conventions and maintain code quality
6. **No Orphaned Code**: Everything implemented gets integrated and used

The prompts are designed to be self-contained while building on previous work, making them suitable for code-generation LLMs working incrementally on this project.