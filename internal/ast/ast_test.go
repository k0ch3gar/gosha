package ast

import (
	"kstmc.com/gosha/internal/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignStatement{
				Token: token.Token{
					Type:    token.INITASSIGN,
					Literal: ":=",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "goyda",
					},
					Value: "goyda",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "baba",
					},
					Value: "baba",
				},
			},
		},
	}

	if program.String() != "goyda := baba" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
