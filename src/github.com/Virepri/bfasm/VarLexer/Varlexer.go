package VarLexer

import (
	"go/types"
	"strings"
	"fmt"
	"strconv"
)

var Variables map[string]Variable = map[string]Variable{}

type Variable struct {
	t types.Type
	array bool
	arrlen int
	v interface{}
}

func LexVars(dat string) {
	for k,v := range strings.Split(dat,"\n") {
		if strings.Index(v, "[") != -1 {
			//Array
			if strings.Index(v,"]") != -1 {
				arlen, _ := strconv.Atoi(v[strings.Index(v,"[")+1:strings.Index(v,"]")])
				Variables[v[:strings.Index(v,"[")]] = Variable{
					array:true,
					arrlen:arlen,
				}
			} else {
				//Error
				fmt.Println("Error on line",k,", Incomplete array")
			}
		} else if len(v) != 0 {
			//Normal variable name
			Variables[v] = Variable{
				array:false,
			}
		}
	}
}