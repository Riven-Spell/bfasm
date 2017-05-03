package Compiler

import (
	"github.com/Virepri/bfasm/Lexer"
	"github.com/Virepri/bfasm/VarLexer"
)

type Allocation struct {
	varname string
	start int
	end int
}

func Compile(lcon []Lexer.Token) string {
	o := ""
	//ptrloc := 0
	line := 0

	allocations := []Allocation{}
	endofallocs := 0

	for k,v := range VarLexer.Variables {
		if v.Array {
			allocations = append(allocations, Allocation{varname:k , start:endofallocs+1 , end:endofallocs+1+v.Arrlen })
			endofallocs += v.Arrlen+1
		} else {
			allocations = append(allocations, Allocation{varname:k , start:endofallocs+1 , end: endofallocs+1})
			endofallocs++
		}
	}


	for _,v := range lcon {
		if v.Lcon != Lexer.VAR && v.Lcon != Lexer.VAL {
			line++
		}

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
		}
	}

	return o
}