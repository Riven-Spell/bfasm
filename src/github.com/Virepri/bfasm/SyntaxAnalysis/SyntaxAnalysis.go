package SyntaxAnalysis

import (
	"github.com/Virepri/bfasm/Lexer"
)

//What do we do here?
/*
Plain and simple: Make sure that syntax is correct.
*/

func AnalyzeSyntax(file string,lcons []Lexer.Token) bool {
	result := true

	for _,v := range lcons {
		switch v.Lcon {
		case Lexer.WHILE:
		case Lexer.IF:
		case Lexer.UNTIL:
		case Lexer.END:
		case Lexer.SET:
		case Lexer.CPY:
		case Lexer.ADD:
		case Lexer.SUB:
		case Lexer.MUL:
		case Lexer.DIV:
		case Lexer.READ:
		case Lexer.PRINT:
		case Lexer.BF:
		case Lexer.VAR:
		case Lexer.VAL:
		}
	}
	return result
}