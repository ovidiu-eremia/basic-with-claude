# Acceptance Tests Guide

This directory contains end-to-end acceptance tests for the BASIC interpreter using a YAML-based test format.

## Test Structure

Tests are organized in YAML files in the `testdata/` directory. Each file contains multiple test cases that verify complete functionality from parsing through execution.

## Adding New Tests

### 1. Choose or Create a Test File

Test files are numbered and grouped by feature:
- `001_basic_functionality.yaml` - Core language features
- `010_input_output.yaml` - INPUT/OUTPUT operations
- `020_arithmetic.yaml` - Mathematical operations
- `050_if_then_basic.yaml` - Conditional statements
- `080_for_next_loops.yaml` - Loop constructs
- `999_error_cases.yaml` - Error handling

### 2. YAML Test Format

Each test has the following structure:

```yaml
tests:
  - name: "TestName"
    program: |
      10 PRINT "HELLO"
      20 END
    inputs:           # Optional: simulated user input
      - "42"
      - "ALICE"
    expected:         # Expected output lines
      - "HELLO\n"
    wantErr: false    # Optional: expect error (default false)
    errContains: ""   # Optional: error message substring
    maxSteps: 1000    # Optional: execution step limit
```

### 3. Key Fields

- **name**: Descriptive test name (used in test output)
- **program**: Multi-line BASIC program using `|` syntax
- **inputs**: Array of strings simulating user INPUT responses
- **expected**: Array of expected output strings (including `\n` for newlines)
- **wantErr**: Set to `true` if test should produce an error
- **errContains**: Substring that must appear in error message
- **maxSteps**: Custom execution limit (default: 1000 steps)

### 4. Output Format Rules

- Each PRINT statement produces one output line ending with `\n`
- Empty PRINT statements produce `"\n"`
- INPUT prompts appear as separate output lines (no newline)
- INPUT values appear as separate output lines with `\n`

### 5. Example Test Cases

**Basic functionality:**
```yaml
- name: "HelloWorld"
  program: |
    10 PRINT "HELLO WORLD"
    20 END
  expected:
    - "HELLO WORLD\n"
```

**With user input:**
```yaml
- name: "Input_With_Prompt"
  program: |
    10 INPUT "ENTER NUMBER"; N
    20 PRINT N
  inputs:
    - "42"
  expected:
    - "ENTER NUMBER"
    - "42\n"
```

**Error testing:**
```yaml
- name: "DivisionByZero"
  program: |
    10 PRINT 1/0
  wantErr: true
  errContains: "DIVISION BY ZERO ERROR"
```

### 6. Running Tests

```bash
# Run all acceptance tests
go test ./acceptance

# Run specific test file pattern
go test ./acceptance -run TestAcceptance
```

## Best Practices

1. **Descriptive names**: Use clear, specific test names
2. **Small focused tests**: One feature/scenario per test
3. **Complete programs**: Include line numbers and END statements
4. **Test both success and error cases**: Use separate tests for error conditions
5. **Verify exact output**: Include all expected output including newlines
6. **Group related tests**: Put similar functionality in the same file
7. **Number files logically**: Use numeric prefixes for ordering (001, 010, 020, etc.)

## File Naming Convention

- Use 3-digit prefixes: `001_`, `010_`, `020_`
- Descriptive suffixes: `_basic_functionality`, `_input_output`
- Always use `.yaml` extension

This system provides comprehensive end-to-end testing ensuring the BASIC interpreter works correctly from parsing through execution.