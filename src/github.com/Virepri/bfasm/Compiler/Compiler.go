package Compiler

import (
	"github.com/Virepri/bfasm/Lexer"
	"github.com/Virepri/bfasm/VarLexer"
	"fmt"
	"strings"
	"strconv"
	"github.com/Virepri/bfasm/SyntaxUtil"
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

var ptrloc uint

func Compile(lcon []Lexer.Token) (string,bool) {
	output := ""
	ptrloc = uint(0)
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
			ref, s, a := getRefPtr(lcon[k+1].Dat,line)

			if s {
				switch a {
				case 1:
					fmt.Println("error: Cannot use array as argument for WHILE. line",line)
				default:
					depthpointers = append(depthpointers,ref)
					depthpointerstype = append(depthpointerstype,Lexer.WHILE)
					output += getMoveOp(ref)
				}
			} else {
				return "",false
			}

			output += "["
			depthpointerstype = append(depthpointerstype,Lexer.WHILE)
		case Lexer.IF:
			allocpoint := bindTempAlloc()

			ifallocs = append(ifallocs,len(tempalloc)-1)

			output += getMoveOp(allocpoint)
			output += "[-]"

			ref,s,a := getRefPtr(lcon[k+1].Dat,line)

			if s {
				switch a {
				case 1:
					fmt.Println("error: Cannot use array as argument for IF. line",line)
					return "",false
				default:
					depthpointers = append(depthpointers,ref)
					depthpointerstype = append(depthpointerstype,Lexer.IF)
					output += getMoveOp(ref)
				}
			} else {
				return "",false
			}

			output += "[" + getMoveOp(allocpoint) + "-]" + getMoveOp(allocpoint) + "[" + getMoveOp(ref)

		case Lexer.UNTIL:
			//TODO: this. I'm gonna give it a bit because I'm sick of programming loops
		case Lexer.END:
			switch depthpointerstype[len(depthpointerstype)-1] {
			case Lexer.WHILE:
				output += getMoveOp(depthpointers[len(depthpointers)-1])
				output += "]"
				depthpointers = depthpointers[:len(depthpointers)-1]
				depthpointerstype = depthpointerstype[:len(depthpointerstype)-1]
			case Lexer.IF:
				output += getMoveOp(tempalloc[ifallocs[len(ifallocs)-1]].start)
				output += "[-]]"
				output += getMoveOp(depthpointers[len(depthpointers)-1])
				depthpointers = depthpointers[:len(depthpointers)-1]
				depthpointerstype = depthpointerstype[:len(depthpointerstype)-1]
				ifallocs = ifallocs[:len(ifallocs)-1]
			case Lexer.UNTIL:
			}

		case Lexer.SET:
			ref,s,a := getRefPtr(lcon[k+1].Dat,line)

			if s {
				switch a {
				case 1:
					//array ref
					if SyntaxUtil.GetValType(lcon[k+2].Dat) == 2 {
						if len(lcon[k+2].Dat) - 2 <= VarLexer.Variables[lcon[k+1].Dat].Arrlen {
							output += getMoveOp(ref)
							output += strings.Repeat("[-]>", len(lcon[k+1].Dat))
							output = output[:len(output)-1]
							ptrloc += uint(len(lcon[k+1].Dat))

							output += getMoveOp(ref)
							for k, v := range []uint8(lcon[k+2].Dat[1:len(lcon[k+2].Dat)-1]) {
								output += getMoveOp(ref + uint(k))
								output += strings.Repeat("+", int(v))
							}
						} else {
							fmt.Println("error: Cannot assign a string larger than the array's size. line",line)
							return "",false
						}
					} else {
						fmt.Println("error: Cannot assign non-string to array. line",line)
						return "",false
					}
				default:
					//var ref
					output += getMoveOp(ref)
					output += "[-]"
					vt := SyntaxUtil.GetValType(lcon[k+2].Dat)
					switch vt {
					case 0:
						//hex
						hexout, _ := strconv.ParseInt(lcon[k+2].Dat,16,16);
						output += strings.Repeat("+", int(hexout))
					case 1:
						//int
						num, _ := strconv.Atoi(lcon[k+2].Dat)
						output += strings.Repeat("+", num)
					case 2:
						//string
						info := uint8(lcon[k+2].Dat[1])
						output += strings.Repeat("+", int(info))
					}
				}
			} else {
				return "",false
			}
		case Lexer.CPY:
			//You might as well kill me as I write this.
			fref, fs, fa := getRefPtr(lcon[k+1].Dat,line)
			tref, ts, ta := getRefPtr(lcon[k+2].Dat,line)

			if fs && ts {
				if fa == 1 {
					//array ref
					if ta == 1 {
						//array ref
						fromname := lcon[k+1].Dat[:strings.Index(lcon[k+1].Dat,"[")]
						toname := lcon[k+2].Dat[:strings.Index(lcon[k+2].Dat,"[")]

						if VarLexer.Variables[fromname].Arrlen <= VarLexer.Variables[toname].Arrlen {
							tempref := bindTempArrayAlloc(uint(VarLexer.Variables[fromname].Arrlen))
							output += getMoveOp(tempref)
							output += strings.Repeat("[-]>", VarLexer.Variables[fromname].Arrlen)
							output += getMoveOp(tref)
							output += strings.Repeat("[-]>", VarLexer.Variables[fromname].Arrlen)
							output += getMoveOp(fref)
							for k,_ := range make([]bool,VarLexer.Variables[fromname].Arrlen) {
								output += getMoveOp(fref+uint(k))
								output += "["
								output += getMoveOp(tempref+uint(k))
								output += "+"
								output += getMoveOp(tref+uint(k))
								output += "+"
								output += getMoveOp(fref+uint(k))
								output += "-]"
							} //copy to tempref and toref
							for k,_ := range make([]bool,VarLexer.Variables[fromname].Arrlen) {
								output += getMoveOp(tempref+uint(k))
								output += "["
								output += getMoveOp(fref+uint(k))
								output += "+"
								output += getMoveOp(tempref+uint(k))
								output += "-]"
							} //destroy tempref and move to fromref
							unbindTempAlloc()
						} else {
							fmt.Println("error: Cannot copy an array larger than the destination to the destination. line",line)
						}
					} else {
						//this won't work.
						fmt.Println("error: Cannot copy an array to a simple variable. line",line)
						return "",false
					}
				} else {
					//typical element/var ref
					if ta == 1 {
						fmt.Println("warning: Copying a simple variable to an array without an element reference will only overwrite the first index. line",line)
					}
					tempref := bindTempAlloc()
					output += getMoveOp(tempref) + "[-]" + getMoveOp(tref) + "[-]" + getMoveOp(fref)          //set tempref and toref to 0
					output += "[" + getMoveOp(tempref) + "+" + getMoveOp(tref) + "+" + getMoveOp(fref) + "-]" //copy fromref to tempref and toref
					output += getMoveOp(tempref) + "[" + getMoveOp(fref) + "+" + getMoveOp(tempref) + "-]"    //set fromref to tempref destructively
					unbindTempAlloc()
				}
			} else {
				return "",false
			}
		case Lexer.ADD:
			switch lcon[k+2].Lcon {
			case Lexer.VAR:
				//do a copy operation but don't destroy arg1
				arg0loc,a0suc,a0r := getRefPtr(lcon[k+1].Dat,line)
				arg1loc,a1suc,a1r := getRefPtr(lcon[k+2].Dat,line)

				if a0suc && a1suc {
					if a0r == 1 {
						if a1r == 1 {
							//adding an array to an array
							if VarLexer.Variables[lcon[k+1].Dat].Arrlen >= VarLexer.Variables[lcon[k+2].Dat].Arrlen {
								//possible
								tloc := bindTempArrayAlloc(uint(VarLexer.Variables[lcon[k+2].Dat].Arrlen))
								output += getMoveOp(tloc) + strings.Repeat("[-]>",VarLexer.Variables[lcon[k+2].Dat].Arrlen)
								for k,_ := range make([]bool,VarLexer.Variables[lcon[k+2].Dat].Arrlen) {
									output += getMoveOp(arg1loc+uint(k))
									output += "[" + getMoveOp(tloc+uint(k)) + "+" + getMoveOp(arg0loc+uint(k)) + "+" + getMoveOp(arg1loc+uint(k)) + "-]"
									output += getMoveOp(tloc+uint(k))
									output += "[" + getMoveOp(arg1loc+uint(k)) + "+" + getMoveOp(tloc+uint(k)) + "-]"
								}
								unbindTempAlloc()
							} else {
								//not possible
								fmt.Println("error: Cannot add array", lcon[k+2].Dat, "to", lcon[k+1].Dat, "as", lcon[k+2].Dat, "is larger than", lcon[k+1].Dat, ". line", line)
								return "", false
							}
						} else {
							//adding a variable to an array
							fmt.Println("warning: Adding a variable to an array only adds to the first element. line",line)
							tloc := bindTempAlloc()
							output += getMoveOp(tloc) + "[-]"
							output += getMoveOp(arg1loc) + "[" + getMoveOp(tloc) + "+" + getMoveOp(arg0loc) + "+" + getMoveOp(arg1loc) + "-]"
							output += getMoveOp(tloc) + "[" + getMoveOp(arg1loc) + "+" + getMoveOp(tloc) + "-]"
							unbindTempAlloc()
						}
					} else {
						if a1r == 1 {
							fmt.Println("error: Cannot add an array to a simple variable (or array element). line",line)
							return "",false
						}
						//adding a variable to a variable
						tloc := bindTempAlloc()
						output += getMoveOp(tloc) + "[-]"
						output += getMoveOp(arg1loc) + "[" + getMoveOp(tloc) + "+" + getMoveOp(arg0loc) + "+" + getMoveOp(arg1loc) + "-]"
						output += getMoveOp(tloc) + "[" + getMoveOp(arg1loc) + "+" + getMoveOp(tloc) + "-]"
						unbindTempAlloc()
					}
				} else {
					return "",false
				}
 			case Lexer.VAL:
				//just kinda do a set operation but don't destroy arg1
				if vt := SyntaxUtil.GetValType(lcon[k+2].Dat); vt != 2 {
					arg0loc, a0suc, _ := getRefPtr(lcon[k+1].Dat, line)
					if a0suc {
						output += getMoveOp(arg0loc)
						num := int64(0)
						if vt == 0 {
							num,_ = strconv.ParseInt(lcon[k+2].Dat,16,16)
						} else {
							tnum,_ := strconv.Atoi(lcon[k+2].Dat)
							num = int64(tnum)
						}
						output += strings.Repeat("+",int(num))
					} else {
						return "", false
					}
				} else {
					fmt.Println("error: Cannot add a string. Like, how does that even work? line",line)
					fmt.Println("If you were wanting to add the ASCII values of the string over an array")
					fmt.Println("Consider \nSET <array1> <string>\nADD <array0> <array1>")
					return "",false
				}
			}
		case Lexer.SUB:
			//same thing as ADD, just with -
			switch lcon[k+2].Lcon {
			case Lexer.VAR:
				//do a copy operation but don't destroy arg1
				arg0loc,a0suc,a0r := getRefPtr(lcon[k+1].Dat,line)
				arg1loc,a1suc,a1r := getRefPtr(lcon[k+2].Dat,line)

				if a0suc && a1suc {
					if a0r == 1 {
						if a1r == 1 {
							//adding an array to an array
							if VarLexer.Variables[lcon[k+1].Dat].Arrlen >= VarLexer.Variables[lcon[k+2].Dat].Arrlen {
								//possible
								tloc := bindTempArrayAlloc(uint(VarLexer.Variables[lcon[k+2].Dat].Arrlen))
								output += getMoveOp(tloc) + strings.Repeat("[-]>",VarLexer.Variables[lcon[k+2].Dat].Arrlen)
								for k,_ := range make([]bool,VarLexer.Variables[lcon[k+2].Dat].Arrlen) {
									output += getMoveOp(arg1loc+uint(k))
									output += "[" + getMoveOp(tloc+uint(k)) + "+" + getMoveOp(arg0loc+uint(k)) + "-" + getMoveOp(arg1loc+uint(k)) + "-]"
									output += getMoveOp(tloc+uint(k))
									output += "[" + getMoveOp(arg1loc+uint(k)) + "+" + getMoveOp(tloc+uint(k)) + "-]"
								}
								unbindTempAlloc()
							} else {
								//not possible
								fmt.Println("error: Cannot subtract array", lcon[k+2].Dat, "from", lcon[k+1].Dat, "as", lcon[k+2].Dat, "is larger than", lcon[k+1].Dat, ". line", line)
								return "", false
							}
						} else {
							//adding a variable to an array
							fmt.Println("warning: Subtracting a variable from an array only subtracts to the first element. line",line)
							tloc := bindTempAlloc()
							output += getMoveOp(tloc) + "[-]"
							output += getMoveOp(arg1loc) + "[" + getMoveOp(tloc) + "+" + getMoveOp(arg0loc) + "-" + getMoveOp(arg1loc) + "-]"
							output += getMoveOp(tloc) + "[" + getMoveOp(arg1loc) + "+" + getMoveOp(tloc) + "-]"
							unbindTempAlloc()
						}
					} else {
						if a1r == 1 {
							fmt.Println("error: Cannot subtract an array from a simple variable (or array element). line",line)
							return "",false
						}
						//adding a variable to a variable
						tloc := bindTempAlloc()
						output += getMoveOp(tloc) + "[-]"
						output += getMoveOp(arg1loc) + "[" + getMoveOp(tloc) + "+" + getMoveOp(arg0loc) + "-" + getMoveOp(arg1loc) + "-]"
						output += getMoveOp(tloc) + "[" + getMoveOp(arg1loc) + "+" + getMoveOp(tloc) + "-]"
						unbindTempAlloc()
					}
				} else {
					return "",false
				}
			case Lexer.VAL:
				//just kinda do a set operation but don't destroy arg1
				if vt := SyntaxUtil.GetValType(lcon[k+2].Dat); vt != 2 {
					arg0loc, a0suc, _ := getRefPtr(lcon[k+1].Dat, line)
					if a0suc {
						output += getMoveOp(arg0loc)
						num := int64(0)
						if vt == 0 {
							num,_ = strconv.ParseInt(lcon[k+2].Dat,16,16)
						} else {
							tnum,_ := strconv.Atoi(lcon[k+2].Dat)
							num = int64(tnum)
						}
						output += strings.Repeat("-",int(num))
					} else {
						return "", false
					}
				} else {
					fmt.Println("error: Cannot subtract a string. Like, how does that even work? line",line)
					fmt.Println("If you were wanting to subtract the ASCII values of the string over an array")
					fmt.Println("Consider \nSET <array1> <string>\nSUB <array0> <array1>")
					return "",false
				}
			}
		case Lexer.MUL:
		case Lexer.DIV:

		case Lexer.READ:
		case Lexer.PRINT:

		case Lexer.BF:
		}
	}

	return output,true
}

func getMoveOp(endptr uint) string {
	o := ""

	if ptrloc > endptr {
		//go back
		o += strings.Repeat("<",int(ptrloc-endptr))
	} else if ptrloc < endptr {
		//go forward
		o += strings.Repeat(">",int(endptr-ptrloc))
	}
	ptrloc = endptr

	return o
}

func getRefPtr(dat string,line int) (uint, bool, uint){
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

func bindTempArrayAlloc(l uint) (uint) {
	o := uint(0)

	tempalloc = append(tempalloc,Allocation{
		varname:"temp"+strconv.Itoa(len(tempalloc)),
		start:uint(endofallocs+1),
		end:uint(endofallocs)+l,
	})
	endofallocs += int(l)
	o = uint(endofallocs)

	return o
}

func unbindTempAlloc() {
	talloc := tempalloc[len(tempalloc)-1]
	if talloc.start == talloc.end {
		//unbinding simple variable
		tempalloc = tempalloc[:len(tempalloc)-1]
		endofallocs--
	} else {
		//unbinding array
		arrlen := (talloc.end - talloc.start) + 1
		tempalloc = tempalloc[:len(tempalloc)-1]
		endofallocs -= int(arrlen)
	}
}