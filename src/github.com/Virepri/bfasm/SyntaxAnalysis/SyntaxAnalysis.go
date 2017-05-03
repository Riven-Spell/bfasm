package SyntaxAnalysis

import (
	"github.com/Virepri/bfasm/Lexer"
	"github.com/Virepri/bfasm/VarLexer"
	"fmt"
	"strings"
	"strconv"
	"github.com/Virepri/bfasm/SyntaxUtil"
)

//What do we do here?
/*
Plain and simple: Make sure that syntax is correct.
*/

var WantedLexicons map[Lexer.Lexicon][][2]Lexer.Lexicon = map[Lexer.Lexicon][][2]Lexer.Lexicon{
	Lexer.WHILE:{{Lexer.VAR}},
	Lexer.IF:{{Lexer.VAR}},
	Lexer.UNTIL:{{Lexer.VAR}},
	Lexer.END:{}, //Flow control

	Lexer.SET:{{Lexer.VAR},{Lexer.VAR,Lexer.VAL}},
	Lexer.CPY:{{Lexer.VAR},{Lexer.VAR}}, //Memory manipulation

	Lexer.ADD:{{Lexer.VAR},{Lexer.VAR,Lexer.VAL}},
	Lexer.SUB:{{Lexer.VAR},{Lexer.VAR,Lexer.VAL}},
	Lexer.MUL:{{Lexer.VAR},{Lexer.VAR,Lexer.VAL}},
	Lexer.DIV:{{Lexer.VAR},{Lexer.VAR,Lexer.VAL}}, //math ops

	Lexer.READ:{{Lexer.VAR},{Lexer.VAL}},
	Lexer.PRINT:{{Lexer.VAR},{Lexer.VAL}}, //IO ops

	Lexer.BF:{{Lexer.VAL},{Lexer.VAL}}, //Special ops
}

//Some really sexy code
func AnalyzeSyntax(lcons []Lexer.Token, line, errors int) bool {
	if len(lcons) == 0 {
		return errors == 0
	}
	if d,t := WantedLexicons[lcons[0].Lcon]; t == true {
		line++
		for k,v := range d {
			if !testLexicon(lcons[1+k].Lcon,v) {
				//failed the test
				/*
				basically means that argument k was not the expected lexicon
				*/

				//the fact that I'm even putting this here is disappointing
				VVTS := map[[2]Lexer.Lexicon]string{
					{Lexer.VAR}: "VAR",
					{Lexer.VAL}:"VAL",
					{Lexer.VAR,Lexer.VAL}:"VAR or VAL",
				}
				fmt.Println("error",errors,":",lcons[0].Dat,"was expecting a",VVTS[v],"but instead got a",lcons[1+k].Dat,"on line:",line)
				errors++
			}
		}
	} else {
		//must be var or val
		if lcons[0].Lcon == Lexer.VAR {
			if strings.Index(lcons[0].Dat,"[") != -1 {
				if strings.Index(lcons[0].Dat,"]") == -1 {
					fmt.Println("error",errors,": Unfinished array reference on line",line)
					errors++
				}
				if _,t := VarLexer.Variables[lcons[0].Dat[:strings.Index(lcons[0].Dat,"[")]]; !t {
					fmt.Println("error",errors,": Variable",lcons[0].Dat[:strings.Index(lcons[0].Dat,"[")],"does not exist. Line",line)
					errors++
				}
			} else {
				if _,t := VarLexer.Variables[lcons[0].Dat]; !t {
					fmt.Println("error",errors,": Variable",lcons[0].Dat,"does not exist. Line",line)
					errors++
				}
			}
		} else {
			if SyntaxUtil.GetValType(lcons[0].Dat) == 3 {
				fmt.Println("error",errors,": Value on line",line,"is an invalid value.")
				errors++
			}
		}
	}
	return AnalyzeSyntax(lcons[1:],line,errors)
}

func testLexicon(lcon Lexer.Lexicon, lcarr [2]Lexer.Lexicon) bool {
	for _,v := range lcarr {
		if lcon == v {
			return true
		}
	}
	return false
}