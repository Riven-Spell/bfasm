package VarLexer

import (
	"strings"
	"fmt"
	"strconv"
)

/*
What all goes on here?
*/

var Variables map[string]Variable = map[string]Variable{}

type Variable struct {
	Array bool
	Arrlen int
	Type uint //0 hex 1 int 2 string 3 invalid
}

func LexVars(dat string) {
	for k,v := range strings.Split(dat,"\n") {
		if strings.Index(v, "[") != -1 {
			//Array
			if strings.Index(v,"]") != -1 {
				arlen, _ := strconv.Atoi(v[strings.Index(v,"[")+1:strings.Index(v,"]")])
				Variables[v[:strings.Index(v,"[")]] = Variable{
					Array:true,
					Arrlen:arlen,
				}
			} else {
				//Error
				fmt.Println("Error on line",k,", Incomplete array")
			}
		} else if len(v) != 0 {
			//Normal variable name
			Variables[v] = Variable{
				Array:false,
			}
		}
	}
}