package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_EndStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		{
			name:     "END statement",
			input:    "10 END",
			expected: program(line(10, 1, endStmt(1))),
		},
		{
			name:  "program with END",
			input: "10 PRINT \"START\"\n20 END\n30 PRINT \"NEVER REACHED\"",
			expected: program(
				line(10, 1, printStmt(str("START", 1), 1)),
				line(20, 2, endStmt(2)),
				line(30, 3, printStmt(str("NEVER REACHED", 3), 3)),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()

			require.Empty(t, p.Errors(), "Parser errors: %v", p.Errors())
			assert.Equal(t, tt.expected, program)
		})
	}
}
