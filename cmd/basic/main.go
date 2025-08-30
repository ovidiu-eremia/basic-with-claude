// ABOUTME: Main CLI entry point for the BASIC interpreter
// ABOUTME: Handles command-line arguments and file loading with proper error handling

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"basic-interpreter/interpreter"
	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

func main() {
	// Define command-line flags
	maxSteps := flag.Int("max-steps", 1000, "Maximum number of execution steps before infinite loop protection triggers")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <filename.bas>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Check for required filename argument
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)

	content, err := readBasicFile(filename)
	if err != nil {
		exitWithError("Error reading file %s: %v", filename, err)
	}

	// Parse the BASIC program
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing error
	if e := p.ParseError(); e != nil {
		// Prepare source lines for context printing (1-based indexing)
		// Normalize newlines in case of Windows files
		normalized := strings.ReplaceAll(content, "\r\n", "\n")
		lines := strings.Split(normalized, "\n")

		// Print offending source line if available (line numbers are 1-based)
		if e.Position.Line >= 1 && e.Position.Line <= len(lines) {
			offending := lines[e.Position.Line-1]
			fmt.Fprintf(os.Stderr, "%s\n", offending)
		}
		fmt.Fprintf(os.Stderr, "line %d: %s\n", e.Position.Line, e.Message)
		os.Exit(1)
	}

	// Execute the program
	fmt.Printf("Program loaded: %s\n", filename)
	fmt.Println("Executing program:")
	fmt.Println()

	// Create runtime and interpreter
	stdRuntime := runtime.NewStandardRuntime()
	interp := interpreter.NewInterpreter(stdRuntime)

	// Configure infinite loop protection
	if *maxSteps > 0 {
		interp.SetMaxSteps(*maxSteps)
	}

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
