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
	return &PrintStatement{BaseNode: BaseNode{Line: line}, Expression: expr}
}

func letStmt(variable string, expr Expression, line int) *LetStatement {
	return &LetStatement{BaseNode: BaseNode{Line: line}, Variable: variable, Expression: expr}
}

func endStmt(line int) *EndStatement {
	return &EndStatement{BaseNode: BaseNode{Line: line}}
}

func runStmt(line int) *RunStatement {
	return &RunStatement{BaseNode: BaseNode{Line: line}}
}

func stopStmt(line int) *StopStatement {
	return &StopStatement{BaseNode: BaseNode{Line: line}}
}

func gotoStmt(targetLine int, line int) *GotoStatement {
	return &GotoStatement{BaseNode: BaseNode{Line: line}, TargetLine: targetLine}
}

func gosubStmt(targetLine int, line int) *GosubStatement {
	return &GosubStatement{BaseNode: BaseNode{Line: line}, TargetLine: targetLine}
}

func returnStmt(line int) *ReturnStatement {
	return &ReturnStatement{BaseNode: BaseNode{Line: line}}
}

func ifStmt(condition Expression, thenStmt Statement, line int) *IfStatement {
	return &IfStatement{BaseNode: BaseNode{Line: line}, Condition: condition, ThenStmt: thenStmt}
}

func inputStmt(prompt string, variable string, line int) *InputStatement {
	return &InputStatement{BaseNode: BaseNode{Line: line}, Prompt: prompt, Variable: variable}
}

func str(value string, line int) *StringLiteral {
	return &StringLiteral{BaseNode: BaseNode{Line: line}, Value: value}
}

func num(value string, line int) *NumberLiteral {
	return &NumberLiteral{BaseNode: BaseNode{Line: line}, Value: value}
}

func varRef(name string, line int) *VariableReference {
	return &VariableReference{BaseNode: BaseNode{Line: line}, Name: name}
}

func binaryOp(left Expression, operator string, right Expression, line int) *BinaryOperation {
	return &BinaryOperation{BaseNode: BaseNode{Line: line}, Left: left, Operator: operator, Right: right}
}

func remStmt(line int) *RemStatement {
	return &RemStatement{BaseNode: BaseNode{Line: line}}
}

func funcCall(name string, args []Expression, line int) *FunctionCall {
	return &FunctionCall{BaseNode: BaseNode{Line: line}, FunctionName: name, Arguments: args}
}
