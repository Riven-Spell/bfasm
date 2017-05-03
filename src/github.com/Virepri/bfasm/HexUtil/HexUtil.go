package HexUtil

import "math"

var data map[rune]int = map[rune]int{
	'0':0,
	'1':1,
	'2':2,
	'3':3,
	'4':4,
	'5':5,
	'6':6,
	'7':7,
	'8':8,
	'9':9,
	'A':10,
	'B':11,
	'C':12,
	'D':13,
	'E':14,
	'F':15,
	'a':10,
	'b':11,
	'c':12,
	'd':13,
	'e':14,
	'f':15,
}

func HexToInt(hex string) int {
	o := 0

	for k,v := range hex {
		o += int(math.Pow(16,float64(k))) * data[v]
	}

	return o
}