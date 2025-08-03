package test

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/parser"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	input := `
return 5
return 10
return 993322
`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}
