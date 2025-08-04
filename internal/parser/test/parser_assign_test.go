package test

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/parser"
	"testing"
)

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      string
	}{
		{"var a = 5", "a", "5"},
		{"var baba = 9 + 10", "baba", "(9 + 10)"},
		{"var foobar = 1 + 2 * 3", "foobar", "(1 + (2 * 3))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.VarStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not *ast.VarStatement. got=%T", program.Statements[0])
			continue
		}

		if stmt.Name.String() != tt.expectedIdentifier {
			t.Errorf("wrong identifier name. got=%q, expected=%q", stmt.Name.String(), tt.expectedIdentifier)
			continue
		}

		if stmt.Value.String() != tt.expectedValue {
			t.Errorf("wrong value. got=%q, expected=%q", stmt.Value.String(), tt.expectedValue)
			continue
		}
	}
}

func TestAssignStatements(t *testing.T) {
	input := `
x = 5;
y = 10
foo = 838383
x = 4
`

	l := lexer.New(input)
	p := parser.New(l)

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
		{"x", "="},
		{"y", "="},
		{"foo", "="},
		{"x", "="},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testAssignStatement(t, stmt, tt.expectedIdentifier, tt.expectedAssign) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
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
