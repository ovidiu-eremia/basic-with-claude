package parser

// Test helper functions for readable AST construction
// Shared across all parser test files to avoid duplication

func program(lines ...*Line) *Program { return &Program{Lines: lines} }

func line(num int, _ int, stmts ...Statement) *Line { return &Line{Number: num, Statements: stmts} }

func printStmt(expr Expression, _ int) *PrintStatement { return &PrintStatement{Expression: expr} }

func letStmt(variable string, expr Expression, _ int) *LetStatement {
	return &LetStatement{Variable: variable, Expression: expr}
}

func endStmt(_ int) *EndStatement { return &EndStatement{} }

func runStmt(_ int) *RunStatement { return &RunStatement{} }

func stopStmt(_ int) *StopStatement { return &StopStatement{} }

func gotoStmt(targetLine int, _ int) *GotoStatement { return &GotoStatement{TargetLine: targetLine} }

func gosubStmt(targetLine int, _ int) *GosubStatement { return &GosubStatement{TargetLine: targetLine} }

func returnStmt(_ int) *ReturnStatement { return &ReturnStatement{} }

func ifStmt(condition Expression, thenStmt Statement, _ int) *IfStatement {
	return &IfStatement{Condition: condition, ThenStmt: thenStmt}
}

func inputStmt(prompt string, variable string, _ int) *InputStatement {
	return &InputStatement{Prompt: prompt, Variable: variable}
}

func str(value string, _ int) *StringLiteral { return &StringLiteral{Value: value} }

func num(value string, _ int) *NumberLiteral { return &NumberLiteral{Value: value} }

func varRef(name string, _ int) *VariableReference { return &VariableReference{Name: name} }

func binaryOp(left Expression, operator string, right Expression, _ int) *BinaryOperation {
	return &BinaryOperation{Left: left, Operator: operator, Right: right}
}

func remStmt(_ int) *RemStatement { return &RemStatement{} }

func funcCall(name string, args []Expression, _ int) *FunctionCall {
	return &FunctionCall{FunctionName: name, Arguments: args}
}
