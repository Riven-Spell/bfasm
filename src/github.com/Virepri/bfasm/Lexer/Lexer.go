package Lexer

import (
	"strings"
	"github.com/Virepri/bfasm/VarLexer"
)

//Not much of a lexical analyzer but yknow it works

type Lexicon uint8
type Token struct {
	Lcon Lexicon
	Dat string
}

const (
	VAR Lexicon = iota
	VAL
	WHILE
	IF
	UNTIL
	END
	SET
	CPY
	ADD
	SUB
	MUL
	DIV
	READ
	PRINT
	BF
)

var Lexicons map[string]Token = map[string]Token{
	"WHILE":{WHILE,"WHILE"},
	"IF":{IF,"IF"},
	"UNTIL":{UNTIL,"UNTIL"},
	"END":{END,"END"},
	"SET":{SET,"SET"},
	"CPY":{CPY,"CPY"},
	"ADD":{ADD,"ADD"},
	"SUB":{SUB,"SUB"},
	"MUL":{MUL,"MUL"},
	"DIV":{DIV,"DIV"},
	"READ":{READ,"READ"},
	"PRINT":{PRINT,"PRINT"},
	"BF":{BF,"BF"},
}

func Lex(dat string) []Token {
	o := make([]Token,0)
	word := ""
	for _, ln := range strings.Split(dat,"\n") {
		for _, v := range strings.Split(ln, " ") {
			lcon, t := Lexicons[v];
			switch {
			case t:
				o = append(o, lcon)
			default:
				if strings.Index(v,"[") != -1 {
					//Definitely a variable
					o = append(o,Token{VAR,v})
				} else if _,t := VarLexer.Variables[v]; t {
					//Definitely a variable
					o = append(o,Token{VAR,v})
				} else if v != "" {
					//Probably not a variable since it's not defined, Lexicon checking (after syntactic analysis) should catch this.
					if strings.Index(v,"\"") != strings.LastIndex(v,"\"") || strings.Index(v,"\"") == -1 {
						o = append(o, Token{VAL, v})
					} else {
						if word != "" {
							word += " "
						}
						word += v
						if strings.LastIndex(word,"\"") == len(word)-1 && strings.Index(word,"\"") == 0 {
							o = append(o,Token{
								VAL,
								word,
							})
							word = ""
						}
					}
				}
			}
		}
	}
	return o
}