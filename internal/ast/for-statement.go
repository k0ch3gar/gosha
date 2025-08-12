package ast

import (
	"bytes"

	"kstmc.com/gosha/internal/token"
)

type ForStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (fs *ForStatement) statementNode() {

}

func (fs *ForStatement) TokenLiteral() string {
	return fs.Token.Literal
}

func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.Token.Literal)
	out.WriteString(" ")
	out.WriteString(fs.Condition.String())
	out.WriteString(" {")
	out.WriteString(fs.Consequence.String())
	out.WriteString("}")

	return out.String()
}
