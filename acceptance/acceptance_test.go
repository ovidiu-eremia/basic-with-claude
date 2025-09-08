// ABOUTME: End-to-end acceptance tests for the BASIC interpreter
// ABOUTME: Tests complete pipeline from YAML test files through execution and output verification

package acceptance

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"basic-interpreter/interpreter"
	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

type YamlTest struct {
	Name        string   `yaml:"name"`
	Program     string   `yaml:"program"`
	Inputs      []string `yaml:"inputs,omitempty"`
	Expected    []string `yaml:"expected,omitempty"`
	WantErr     bool     `yaml:"wantErr,omitempty"`
	ErrContains string   `yaml:"errContains,omitempty"`
	MaxSteps    int      `yaml:"maxSteps,omitempty"`
}

type YamlTestFile struct {
	Tests []YamlTest `yaml:"tests"`
}

type AcceptanceTest struct {
	name        string
	program     string
	inputs      []string
	expected    []string
	wantErr     bool
	errContains string
	maxSteps    int // Custom max steps limit, 0 means use default
}

// loadTestsFromYAML loads all YAML test files from testdata directory
func loadTestsFromYAML(t *testing.T) []AcceptanceTest {
	t.Helper()

	testdataDir := "testdata"
	entries, err := os.ReadDir(testdataDir)
	require.NoError(t, err, "Failed to read testdata directory")

	// Sort entries to ensure predictable test order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var allTests []AcceptanceTest
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			filePath := filepath.Join(testdataDir, entry.Name())
			tests := loadTestFile(t, filePath)
			allTests = append(allTests, tests...)
		}
	}

	return allTests
}

// loadTestFile loads tests from a single YAML file
func loadTestFile(t *testing.T, filePath string) []AcceptanceTest {
	t.Helper()

	data, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read test file %s", filePath)

	var yamlFile YamlTestFile
	err = yaml.Unmarshal(data, &yamlFile)
	require.NoError(t, err, "Failed to parse YAML file %s", filePath)

	var tests []AcceptanceTest
	for _, yamlTest := range yamlFile.Tests {
		test := AcceptanceTest{
			name:        yamlTest.Name,
			program:     yamlTest.Program,
			inputs:      yamlTest.Inputs,
			expected:    yamlTest.Expected,
			wantErr:     yamlTest.WantErr,
			errContains: yamlTest.ErrContains,
			maxSteps:    yamlTest.MaxSteps,
		}
		tests = append(tests, test)
	}

	return tests
}

// executeBasicProgramWithMaxSteps parses and executes a BASIC program string with custom max steps
func executeBasicProgramWithMaxSteps(t *testing.T, program string, inputs []string, maxSteps int) ([]string, error) {
	t.Helper()

	// Parse the program
	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.Errors())
	}
	if ast == nil {
		return nil, fmt.Errorf("parsing returned nil AST")
	}

	// Create test runtime and interpreter
	testRuntime := runtime.NewTestRuntime()
	if len(inputs) > 0 {
		testRuntime.SetInput(inputs)
	}
	interp := interpreter.NewInterpreter(testRuntime)

	// Set custom max steps if specified
	if maxSteps > 0 {
		interp.SetMaxSteps(maxSteps)
	}

	// Execute the program
	err := interp.Execute(ast)
	if err != nil {
		return nil, err
	}

	// Return captured output
	return testRuntime.GetOutput(), nil
}

const DEFAULT_MAX_STEPS = 1000

func TestAcceptance(t *testing.T) {
	tests := loadTestsFromYAML(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output []string
			var err error
			if tt.maxSteps == 0 {
				tt.maxSteps = DEFAULT_MAX_STEPS
			}
			output, err = executeBasicProgramWithMaxSteps(t, tt.program, tt.inputs, tt.maxSteps)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, output)
			}
		})
	}
}
