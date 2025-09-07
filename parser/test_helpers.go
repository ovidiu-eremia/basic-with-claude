package parser

// Test helper functions for readable AST construction
// Shared across all parser test files to avoid duplication

func program(lines ...*Line) *Program {
	return &Program{Lines: lines}
}

func line(num int, sourceLine int, stmts ...Statement) *Line {
	return &Line{Number: num, SourceLine: sourceLine, Statements: stmts}
}

func printStmt(expr Expression, line int) *PrintStatement {
	return &PrintStatement{Expression: expr, Line: line}
}

func letStmt(variable string, expr Expression, line int) *LetStatement {
	return &LetStatement{Variable: variable, Expression: expr, Line: line}
}

func endStmt(line int) *EndStatement {
	return &EndStatement{Line: line}
}

func runStmt(line int) *RunStatement {
	return &RunStatement{Line: line}
}

func stopStmt(line int) *StopStatement {
	return &StopStatement{Line: line}
}

func gotoStmt(targetLine int, line int) *GotoStatement {
	return &GotoStatement{TargetLine: targetLine, Line: line}
}

func gosubStmt(targetLine int, line int) *GosubStatement {
	return &GosubStatement{TargetLine: targetLine, Line: line}
}

func returnStmt(line int) *ReturnStatement {
	return &ReturnStatement{Line: line}
}

func ifStmt(condition Expression, thenStmt Statement, line int) *IfStatement {
	return &IfStatement{Condition: condition, ThenStmt: thenStmt, Line: line}
}

func inputStmt(prompt string, variable string, line int) *InputStatement {
	return &InputStatement{Prompt: prompt, Variable: variable, Line: line}
}

func str(value string, line int) *StringLiteral {
	return &StringLiteral{Value: value, Line: line}
}

func num(value string, line int) *NumberLiteral {
	return &NumberLiteral{Value: value, Line: line}
}

func varRef(name string, line int) *VariableReference {
	return &VariableReference{Name: name, Line: line}
}

func binaryOp(left Expression, operator string, right Expression, line int) *BinaryOperation {
	return &BinaryOperation{Left: left, Operator: operator, Right: right, Line: line}
}

func remStmt(line int) *RemStatement {
	return &RemStatement{Line: line}
}

func funcCall(name string, args []Expression, line int) *FunctionCall {
	return &FunctionCall{FunctionName: name, Arguments: args, Line: line}
}
