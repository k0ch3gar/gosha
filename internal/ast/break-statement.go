package ast

import "kstmc.com/gosha/internal/token"

type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode() {

}

func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BreakStatement) String() string {
	return bs.Token.Literal
}
