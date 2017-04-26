package SyntaxAnalysis

import (
	"github.com/Virepri/bfasm/Lexer"
	"strings"
)

//What do we do here?
/*
Plain and simple: Make sure that syntax is correct.
*/

func AnalyzeSyntax(file string,lcons []Lexer.Lexicon) bool {
	result := true
	fLCons := []string{}
	//prepare an equivalent list of string lexicons
	for _,v := range strings.Split(file,"\n") {
		for _,lcon := range strings.Split(v," ") {
			if lcon != "" {
				fLCons = append(fLCons,lcon)
			}
		}
	}

	for _,v := range lcons {
		switch v {
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