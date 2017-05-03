package Compiler

import (
	"github.com/Virepri/bfasm/Lexer"
	"github.com/Virepri/bfasm/VarLexer"
	"fmt"
	"strings"
	"strconv"
)

type Allocation struct {
	varname string
	start uint
	end uint
}

func Compile(lcon []Lexer.Token) (string,bool) {
	o := ""
	ptrloc := uint(0)
	depthpointers := []uint {}
	line := 0

	allocref := map[string]*Allocation{}
	allocations := []Allocation{}
	endofallocs := 0

	for k,v := range VarLexer.Variables {
		if v.Array {
			allocations = append(allocations, Allocation{varname:k , start:endofallocs+1 , end:endofallocs+1+v.Arrlen })
			allocref[k] = &allocations[len(allocations)-1]
			endofallocs += v.Arrlen+1
		} else {
			allocations = append(allocations, Allocation{varname:k , start:endofallocs+1 , end: endofallocs+1})
			allocref[k] = &allocations[len(allocations)-1]
			endofallocs++
		}
	}


	for k,v := range lcon {
		if v.Lcon != Lexer.VAR && v.Lcon != Lexer.VAL {
			line++
		}

		switch v.Lcon {
		case Lexer.WHILE:
			if strings.Index(lcon[k+1].Dat,"[") != -1 {
				name := lcon[k+1].Dat[:strings.Index(lcon[k+1].Dat,"[")];
				if !VarLexer.Variables[name].Array {
					fmt.Println("error: Attempted creation of reference pointer to non-array object",name,"on line",line,"failed.")
					return "", false
				}
				//references an array, SyntaxAnalysis makes sure that this is a valid array.
				if point, err := strconv.Atoi(lcon[k+1].Dat[strings.Index(lcon[k+1].Dat,"[")+1:strings.Index(lcon[k+1].Dat,"]")]); err == nil {
					if VarLexer.Variables[name].Arrlen > point && point >= 0{
						depthpointers = append(depthpointers,allocref[name].start+uint(point))
						o += getMoveOp(ptrloc,depthpointers[len(depthpointers)-1])
					} else {
						fmt.Println("error: Reference pointer to",point,"on array",lcon[k+1].Dat[:strings.Index(lcon[k+1].Dat,"[")],"due to being out-of-bounds")
					}
				} else {
					fmt.Println("error: Could not convert reference pointer to int on line",line)
					return "", false
				}
			} else {
				//references a variable.
				if VarLexer.Variables[v.Dat].Array {
					fmt.Println("error: Cannot use an array as a condition of WHILE. line",line)
					return "",false
				}
				depthpointers = append(depthpointers,allocref[v.Dat].start)
				o += getMoveOp(ptrloc,depthpointers[len(depthpointers)-1])
			}
			o += "["
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

func getMoveOp(currptr, endptr uint) string {
	o := ""

	if currptr > endptr {
		//go back
		o += strings.Repeat("<",int(currptr-endptr))
	} else if currptr < endptr {
		//go forward
		o += strings.Repeat(">",int(endptr-currptr))
	}

	return o
}