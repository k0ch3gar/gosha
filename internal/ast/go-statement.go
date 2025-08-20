package ast

import (
	"bytes"

	"kstmc.com/gosha/internal/token"
)

type GoStatement struct {
	Token token.Token
	Expr  Expression
}

func (gs *GoStatement) statementNode() {

}

func (gs *GoStatement) TokenLiteral() string {
	return gs.Token.Literal
}

func (gs *GoStatement) String() string {
	var out bytes.Buffer

	out.WriteString("go ")
	out.WriteString(gs.Expr.String())

	return out.String()
}
