package ast

import "kstmc.com/gosha/internal/token"

type Identifier struct {
	Token    token.Token
	Value    string
	DataType *DataType
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
