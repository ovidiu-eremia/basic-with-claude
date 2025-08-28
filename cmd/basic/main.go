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
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
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
	rt := runtime.NewStandardRuntime()
	interp := interpreter.New(rt)
	
	// Execute the program
	err = interp.Execute(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}

// readBasicFile reads the contents of a BASIC program file
func readBasicFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// displayParsedProgram displays the parsed AST structure for debugging
func displayParsedProgram(program *parser.Program) {
	for _, line := range program.Lines {
		fmt.Printf("Line %d:\n", line.Number)
		for _, stmt := range line.Statements {
			switch s := stmt.(type) {
			case *parser.PrintStatement:
				fmt.Printf("  PRINT ")
				switch expr := s.Expression.(type) {
				case *parser.StringLiteral:
					fmt.Printf("string: %q\n", expr.Value)
				default:
					fmt.Printf("unknown expression type\n")
				}
			case *parser.EndStatement:
				fmt.Printf("  END\n")
			default:
				fmt.Printf("  unknown statement type\n")
			}
		}
		fmt.Println()
	}
}