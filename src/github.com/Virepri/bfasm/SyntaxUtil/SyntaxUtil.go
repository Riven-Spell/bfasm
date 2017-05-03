package SyntaxUtil

import (
	"strings"
	"strconv"
)

func GetValType(dat string) uint {
	if strings.Index(dat,"0x") == 0 {
		//hex
		return 0
	} else if _,err := strconv.Atoi(dat); err == nil {
		//int
		return 1
	} else if strings.Index(dat,"\"") != strings.LastIndex(dat,"\"") {
		//string
		return 2
	} else {
		//not a valid value
		return 3
	}
}