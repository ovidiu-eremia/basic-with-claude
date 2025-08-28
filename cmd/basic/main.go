// ABOUTME: Main CLI entry point for the BASIC interpreter
// ABOUTME: Handles command-line arguments and file loading with proper error handling

package main

import (
	"fmt"
	"os"

	"basic-interpreter/interpreter"
	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.bas>\n", os.Args[0])
		os.Exit(1)
	}
	
	filename := os.Args[1]
	
	content, err := readBasicFile(filename)
	if err != nil {
		exitWithError("Error reading file %s: %v", filename, err)
	}
	
	// Parse the BASIC program
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()
	
	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Parsing errors:\n")
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "  %s\n", err)
		}
		os.Exit(1)
	}
	
	// Execute the program
	fmt.Printf("Program loaded: %s\n", filename)
	fmt.Println("Executing program:")
	fmt.Println()
	
	// Create runtime and interpreter
	stdRuntime := runtime.NewStandardRuntime()
	interp := interpreter.NewInterpreter(stdRuntime)
	
	// Execute the program
	err = interp.Execute(program)
	if err != nil {
		exitWithError("Runtime error: %v", err)
	}
}

// exitWithError prints an error message and exits with code 1
func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// readBasicFile reads the contents of a BASIC program file
func readBasicFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

