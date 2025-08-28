package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_ParseProgram(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		{
			name:  "single line with PRINT",
			input: `10 PRINT "HELLO"`,
			expected: &Program{
				Lines: []*Line{
					{
						Number: 10,
						Statements: []Statement{
							&PrintStatement{
								Expression: &StringLiteral{
									Value: "HELLO",
									Line:  1,
								},
								Line: 1,
							},
						},
						Line: 1,
					},
				},
			},
		},
		{
			name:  "multiple lines",
			input: "10 PRINT \"LINE1\"\n20 PRINT \"LINE2\"",
			expected: &Program{
				Lines: []*Line{
					{
						Number: 10,
						Statements: []Statement{
							&PrintStatement{
								Expression: &StringLiteral{
									Value: "LINE1",
									Line:  1,
								},
								Line: 1,
							},
						},
						Line: 1,
					},
					{
						Number: 20,
						Statements: []Statement{
							&PrintStatement{
								Expression: &StringLiteral{
									Value: "LINE2",
									Line:  2,
								},
								Line: 2,
							},
						},
						Line: 2,
					},
				},
			},
		},
		{
			name:  "empty string",
			input: `10 PRINT ""`,
			expected: &Program{
				Lines: []*Line{
					{
						Number: 10,
						Statements: []Statement{
							&PrintStatement{
								Expression: &StringLiteral{
									Value: "",
									Line:  1,
								},
								Line: 1,
							},
						},
						Line: 1,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			
			program := p.ParseProgram()
			
			require.NotNil(t, program, "ParseProgram() returned nil")
			require.Empty(t, p.Errors(), "Parser errors: %v", p.Errors())
			
			assert.Equal(t, tt.expected, program)
		})
	}
}

func TestParser_ParseErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectError   bool
	}{
		{
			name:        "unterminated string",
			input:       `10 PRINT "HELLO`,
			expectError: true,
		},
		{
			name:        "missing line number",
			input:       `PRINT "HELLO"`,
			expectError: true,
		},
		{
			name:        "invalid syntax",
			input:       `10 INVALID "HELLO"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			
			program := p.ParseProgram()
			
			if tt.expectError {
				assert.True(t, len(p.Errors()) > 0, "Expected parsing errors but got none")
			} else {
				assert.Empty(t, p.Errors(), "Expected no parsing errors but got: %v", p.Errors())
				assert.NotNil(t, program)
			}
		})
	}
}