package parser

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/lexer"
	"testing"
)

func TestAssignStatements(t *testing.T) {
	input := `
x  5;
y := 10
foo := 838383
x = 4
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedAssign     string
	}{
		{"x", ":="},
		{"y", ":="},
		{"foo", ":="},
		{"x", "="},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testAssignStatement(t, stmt, tt.expectedIdentifier, tt.expectedAssign) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testAssignStatement(t *testing.T, s ast.Statement, name string, assign string) bool {
	if s.TokenLiteral() != assign {
		t.Errorf("s.TokenLiteral() not ':='. got=%q", s.TokenLiteral())
		return false
	}

	assignStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if assignStmt.Name.Value != name {
		t.Errorf("assignStmt.Name.Value not '%s'. got=%s", name, assignStmt.Name.Value)
		return false
	}

	if assignStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, assignStmt.Name)
		return false
	}

	return true
}
