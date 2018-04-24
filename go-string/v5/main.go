package main

import "fmt"

func main() {
	s := "éक्षिaπ汉字"
	for i, rn := range s {
		fmt.Printf("%v: 0x%x %v  %d\n", i, rn, string(rn), len(string(rn)))
	}
}

/*
0: 0x65 e  1
1: 0x301 ́  2
3: 0x915 क  3
6: 0x94d ्  3
9: 0x937 ष  3
12: 0x93f ि  3
15: 0x61 a  1
16: 0x3c0 π  2
18: 0x6c49 汉  3
21: 0x5b57 字  3
*/
// the first character, é, is composed of two runes (3 bytes total)
// the second character, क्षि, is composed of four runes (12 bytes total).
// the English character, a, is composed of one rune (1 byte).
// the character, π, is composed of one rune (2 bytes).
// each of the two Chinese characters, 汉字, is composed of one rune (3 bytes each).
