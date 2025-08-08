package analyzer

import (
	"testing"

	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/parser"
)

func TestVarStatements(t *testing.T) {
	input := `
	var a = 5
	var b int = 1000
	var c = true
	var d bool = false
	`

	errors := testAnalyze(input)
	if len(errors) != 0 {
		for _, error := range errors {
			t.Error(error)
		}
	}
}

func testAnalyze(input string) []string {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	env := object.NewEnvironment()
	return AnalyzeProgram(program, env)
}
