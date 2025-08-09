package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"kstmc.com/gosha/internal/analyzer"
	"kstmc.com/gosha/internal/evaluator"
	"kstmc.com/gosha/internal/lexer"
	"kstmc.com/gosha/internal/object"
	"kstmc.com/gosha/internal/parser"
)

const (
	PROMPT = "gosha>> "
)

func Start(in io.Reader, out io.Writer) {
	file, ok := in.(*os.File)

	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	if !(ok && file == os.Stdin) {
		data, err := io.ReadAll(in)
		if err != nil {
			fmt.Errorf("unable to read data from file: %s", err.Error())
			return
		}

		startPos := bytes.IndexByte(data, '\n')

		processInput(out, string(data[startPos:]), env)
		return
	}

	for {
		io.WriteString(out, PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		processInput(out, line, env)
	}
}

func processInput(out io.Writer, input string, env *object.Environment) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	errors := analyzer.AnalyzeProgram(program, env)
	if len(errors) != 0 {
		printParserErrors(out, errors)
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() != parser.NIL {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
