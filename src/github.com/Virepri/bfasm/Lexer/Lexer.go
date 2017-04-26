package Lexer

import (
	"strings"
	"github.com/Virepri/bfasm/VarLexer"
)

//Not much of a lexical analyzer but yknow it works

type Lexicon uint8

const (
	VAR Lexicon = iota
	VAL
	WHILE
	IF
	UNTIL
	END
	SET
	CPY
	SUB
	MUL
	DIV
	READ
	PRINT
	BF
)

var Lexicons map[string]Lexicon = map[string]Lexicon{
	"WHILE":WHILE,
	"IF":IF,
	"UNTIL":UNTIL,
	"END":END,
	"SET":SET,
	"CPY":CPY,
	"SUB":SUB,
	"MUL":MUL,
	"DIV":DIV,
	"READ":READ,
	"PRINT":PRINT,
	"BF":BF,
}

func Lex(dat string) []Lexicon {
	o := make([]Lexicon,0)
	for _, ln := range strings.Split(dat,"\n") {
		for _, v := range strings.Split(ln, " ") {
			lcon, t := Lexicons[v];
			switch {
			case t:
				o = append(o, lcon)
			default:
				if strings.Index(v,"[") != -1 {
					//Definitely a variable
					o = append(o,VAR)
				} else if _,t := VarLexer.Variables[v]; t {
					//Definitely a variable
					o = append(o,VAR)
				} else if v != "" {
					//Probably not a variable since it's not defined, Lexicon checking (after syntactic analysis) should catch this.
					o = append(o,VAL)
				}
			}
		}
	}
	return o
}