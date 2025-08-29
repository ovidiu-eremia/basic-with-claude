# Bugs and Edge Cases Found During Exploratory Testing

This document lists bugs, edge cases, and areas for improvement found during exploratory testing of the BASIC interpreter.

## 1. Division by Zero Error Message

*   **Description:** When a division by zero occurs, the interpreter exits with a generic "Runtime error: division by zero".
*   **Expected Behavior:** The error message should be in the C64 BASIC style, e.g., `?DIVISION BY ZERO ERROR IN 10`.
*   **To Reproduce:** Run a program with a division by zero, e.g., `10 PRINT 1/0`.
*   **Severity:** Low. The error is caught, but the message format is inconsistent.

## 2. String Concatenation with '+' Operator

*   **Description:** The `+` operator is not implemented for string concatenation. The interpreter attempts to convert strings to numbers, resulting in a runtime error.
*   **Expected Behavior:** The `+` operator should concatenate strings. For example, `PRINT "HELLO" + " " + "WORLD"` should output `HELLO WORLD`.
*   **To Reproduce:** Run a program that uses `+` to concatenate strings, e.g., `10 PRINT "HELLO" + " " + "WORLD"`.
*   **Severity:** High. This is a missing implementation of a standard BASIC feature.

## 3. Infinite Loop Handling

*   **Description:** The interpreter enters an infinite loop when a program like `10 GOTO 10` is executed. While it doesn't crash or consume excessive CPU, there is no mechanism to detect or break out of the loop.
*   **Expected Behavior:** Ideally, the interpreter should have a mechanism to detect and handle infinite loops, perhaps by providing a way to interrupt execution (e.g., with Ctrl+C) or by implementing a timeout.
*   **To Reproduce:** Run a program with an infinite loop, e.g., `10 GOTO 10`.
*   **Severity:** Medium. The interpreter becomes unresponsive and must be killed manually.

## 4. Case Sensitivity of Keywords

*   **Description:** The interpreter's parser treats keywords as case-sensitive. For example, `let` and `print` are not recognized as valid keywords.
*   **Expected Behavior:** BASIC is a case-insensitive language. Keywords should be recognized regardless of their case (e.g., `PRINT`, `print`, `Print`).
*   **To Reproduce:** Run a program with lowercase keywords, e.g., `10 let a = 5
20 print a`.
*   **Severity:** High. This violates a fundamental aspect of the BASIC language and breaks compatibility.

## 5. `PRINT` Statement without Arguments

*   **Description:** A `PRINT` statement without any arguments causes a parsing error.
*   **Expected Behavior:** In most BASIC dialects, a `PRINT` statement on its own is valid and should print a newline character.
*   **To Reproduce:** Run a program with a `PRINT` statement on a line by itself, e.g., `10 PRINT`.
*   **Severity:** Medium. This is a missing implementation of a standard BASIC feature.

## 6. Variable Name Length

*   **Description:** The interpreter does not correctly truncate variable names to 2 significant characters as per the language specification. Variables with names longer than 2 characters are treated as distinct variables.
*   **Expected Behavior:** Variable names should be significant to . For example, `VA` and `VARA` should be treated as the same variable.
*   **To Reproduce:** Run the following program:
    ```basic
    10 VA = 5
    20 VARA = 10
    30 PRINT VA
    ```
    The output should be `10`, but it is `5`.
*   **Severity:** High. This is a violation of the language specification.

## 7. Uninitialized String Variables

*   **Description:** Uninitialized string variables are treated as the number 0 instead of an empty string.
*   **Expected Behavior:** Uninitialized string variables should be treated as an empty string (`""`). Printing an uninitialized string variable should produce no visible output.
*   **To Reproduce:** Run the following program:
    ```basic
    10 PRINT A$
    ```
    The output should be a blank line, but it is `0`.
*   **Severity:** High. This can lead to unexpected behavior and type mismatch errors.

## 8. No Type Mismatch Error on String Assignment to Numeric Variable

*   **Description:** The interpreter does not raise a `TYPE MISMATCH` error when a string is assigned to a numeric variable. The assignment fails silently.
*   **Expected Behavior:** A `TYPE MISMATCH` error should be reported.
*   **To Reproduce:** Run the following program:
    ```basic
    10 A = "hello"
    ```
*   **Severity:** Critical. This can lead to subtle bugs and unexpected program behavior.

## 9. No Type Mismatch Error on Mixed-Type Addition

*   **Description:** The interpreter does not raise a `TYPE MISMATCH` error when a string and a number are added. It appears to perform an implicit conversion of the string to a number.
*   **Expected Behavior:** A `TYPE MISMATCH` error should be reported. The `spec.md` does not specify any implicit type conversion rules for arithmetic operations.
*   **To Reproduce:** Run the following program:
    ```basic
    10 A$ = "10"
    20 B = 5
    30 C = A$ + B
    ```
*   **Severity:** Critical. This can lead to incorrect calculations and violates the principle of strong typing that is expected in this context.

## 10. Parser Errors with Arithmetic Expressions

*   **Description:** The parser fails to correctly parse arithmetic expressions involving unary minus signs, floating-point numbers, and multiple operators.
*   **Expected Behavior:** The parser should correctly handle these common arithmetic expressions.
*   **To Reproduce:** Run a program with expressions like:
    ```basic
    10 PRINT -10 + 5
    20 PRINT 2.5 * 2
    30 PRINT 2 * 3 + 4 * 5
    ```
*   **Severity:** Critical. This prevents the use of basic arithmetic operations, a core feature of the language.