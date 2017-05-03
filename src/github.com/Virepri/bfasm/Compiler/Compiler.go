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

var allocref map[string]*Allocation
var allocations []Allocation
var tempalloc []Allocation
var endofallocs int
var ifallocs []int

func Compile(lcon []Lexer.Token) (string,bool) {
	o := ""
	ptrloc := uint(0)
	depthpointers := []uint {}
	depthpointerstype := []Lexer.Lexicon{}
	line := 0

	allocref = map[string]*Allocation{}
	allocations = []Allocation{}
	tempalloc = []Allocation{}
	endofallocs = 0

	ifallocs = []int{}

	for k,v := range VarLexer.Variables {
		if v.Array {
			allocations = append(allocations, Allocation{varname:k , start:uint(endofallocs+1) , end:uint(endofallocs+1+v.Arrlen) })
			allocref[k] = &allocations[len(allocations)-1]
			endofallocs += v.Arrlen+1
		} else {
			allocations = append(allocations, Allocation{varname:k , start:uint(endofallocs+1) , end:uint(endofallocs+1)})
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
			ref, s, a := getRefPtr(lcon[k+1].Dat,line,allocref)

			if s {
				switch a {
				case 1:
					fmt.Println("error: Cannot use array as argument for WHILE. line",line)
				default:
					depthpointers = append(depthpointers,ref)
					depthpointerstype = append(depthpointerstype,Lexer.WHILE)
					o += getMoveOp(ptrloc,ref)
				}
			} else {
				return "",false
			}

			o += "["
			depthpointerstype = append(depthpointerstype,Lexer.WHILE)
		case Lexer.IF:
			allocpoint := bindTempAlloc()

			ifallocs = append(ifallocs,len(tempalloc)-1)

			o += getMoveOp(ptrloc,allocpoint)
			o += "[-]"

			ref,s,a := getRefPtr(lcon[k+1].Dat,line,allocref)

			if s {
				switch a {
				case 1:
					fmt.Println("error: Cannot use array as argument for IF. line",line)
					return "",false
				default:
					depthpointers = append(depthpointers,ref)
					depthpointerstype = append(depthpointerstype,Lexer.IF)
					o += getMoveOp(ptrloc,ref)
				}
			} else {
				return "",false
			}

			o += "[" + getMoveOp(ptrloc,allocpoint) + "-]" + getMoveOp(ptrloc,allocpoint) + "[" + getMoveOp(ptrloc,ref)

		case Lexer.UNTIL:
		case Lexer.END:
			switch depthpointerstype[len(depthpointerstype)-1] {
			case Lexer.WHILE:
				o += getMoveOp(ptrloc,depthpointers[len(depthpointers)-1])
				o += "]"
				depthpointers = depthpointers[:len(depthpointers)-1]
				depthpointerstype = depthpointerstype[:len(depthpointerstype)-1]
			case Lexer.IF:
				o += getMoveOp(ptrloc,tempalloc[ifallocs[len(ifallocs)-1]].start)
				o += "[-]]"
				o += getMoveOp(ptrloc,depthpointers[len(depthpointers)-1])
				depthpointers = depthpointers[:len(depthpointers)-1]
				depthpointerstype = depthpointerstype[:len(depthpointerstype)-1]
				ifallocs = ifallocs[:len(ifallocs)-1]
			case Lexer.UNTIL:
			}

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

	return o,true
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

func getRefPtr(dat string,line int,allocref map[string]*Allocation) (uint, bool, uint){
	loc, success, arrref := uint(0),true,uint(0)
	//arrref: 0 = not an array pointer
	//1 = pointer to array name
	//2 = pointer to array element

	if strings.Index(dat,"[") < strings.Index(dat,"]") && (strings.Index(dat,"[") != -1 && strings.Index(dat,"]") != -1) {
		//confirmed to be an array reference

		numstr := dat[strings.Index(dat,"[")+1:strings.Index(dat,"]")]
		namestr := dat[:strings.Index(dat,"[")]

		if num,err := strconv.Atoi(numstr); err == nil {
			if VarLexer.Variables[namestr].Array {
				if num < VarLexer.Variables[namestr].Arrlen && num >= 0 {
					loc = allocref[namestr].start + uint(num)
					arrref = 2
				} else {
					fmt.Println("error:",num,"is out-of-bounds on array",namestr,"on line",line,"max:",VarLexer.Variables[namestr].Arrlen)
					success = false
				}
			} else {
				fmt.Println("error: Cannot create array reference to non-array", namestr, "on line", line)
				success = false
			}
		} else {
			fmt.Println("error: Cannot reference array with non-integer",numstr,"on line",line)
			success = false
		}
	} else {
		//not an array reference
		//if it's reached here, it's definitely a variable. Assume so and return the location.
		if VarLexer.Variables[dat].Array {
			arrref = 1
		}
		loc = allocref[dat].start
	}

	return loc,success,arrref
}

func bindTempAlloc() uint {
	o := uint(0)

	tempalloc = append(tempalloc,Allocation{
		varname:"temp"+strconv.Itoa(len(tempalloc)),
		start:uint(endofallocs+1),
		end:uint(endofallocs+1),
	})
	endofallocs++
	o = uint(endofallocs)

	return o
}