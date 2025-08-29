package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"basic-interpreter/lexer"
)

func TestIfStatementParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		{
			name:  "simple IF THEN",
			input: "10 IF 1 THEN PRINT \"TRUE\"",
			expected: &Program{
				Lines: []*Line{
					{
						Number: 10,
						Statements: []Statement{
							&IfStatement{
								Condition: &NumberLiteral{
									Value: "1",
									Line:  1,
								},
								ThenStmt: &PrintStatement{
									Expression: &StringLiteral{
										Value: "TRUE",
										Line:  1,
									},
									Line: 1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
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

			if len(p.Errors()) > 0 {
				t.Errorf("Parser errors: %v", p.Errors())
				return
			}

			assert.Equal(t, tt.expected, program)
		})
	}
}
