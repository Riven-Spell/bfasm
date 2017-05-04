package SyntaxUtil

import (
	"strings"
	"strconv"
)

func GetValType(dat string) uint {
	if _,err := strconv.ParseInt(dat,16,16); err == nil {
		//hex
		return 0
	} else if _,err := strconv.Atoi(dat); err == nil {
		//int
		return 1
	} else if strings.Index(dat,"\"") != strings.LastIndex(dat,"\"") {
		//string
		return 2
	} else if strings.Index(dat,"'") != strings.LastIndex(dat,"'"){
		//char, treat it as a string.
		return 2
	} else {
		return 3
	}
}