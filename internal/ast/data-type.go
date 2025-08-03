package ast

import "kstmc.com/gosha/internal/token"

type DataType struct {
	Token token.Token
	Name  string
}
