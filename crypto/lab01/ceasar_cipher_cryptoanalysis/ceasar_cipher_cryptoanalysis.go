package main

import (
	"bufio"
	"fmt"
	"os"
)

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 52

func main() {
	initDict()
	encStr, plainText := getInput()

	dec, k := analyseText(encStr, plainText)

	fmt.Printf("%s\n%d", string(dec), k)
}

func analyseText(encStr, decStr []rune) ([]rune, int) {
	i := 0
	for j, c := range decStr {
		if c == ' ' || !isLiteral(c) {
			continue
		}
		i = j
		break
	}

	k := mod(dict[encStr[i]]-dict[decStr[i]], dictLen)

	s := ceasarCipherDecrypt(encStr, k, dictLen)
	return s, k
}

func ceasarCipherDecrypt(str []rune, k int, z int) []rune {
	dec := ""
	for _, c := range str {
		if !isLiteral(c) {
			dec += string(c)
			continue
		}

		decVal := mod(dict[c]-k, z)
		dec += string(dictRev[decVal])
	}

	return []rune(dec)
}

func isLiteral(c rune) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func mod(num int, z int) int {
	num = num % z
	if num < 0 {
		num += z
	}
	return num
}

func initDict() {
	dict = make(map[rune]int)
	dictRev = make(map[int]rune)

	offset := 97
	for i := 0; i < 26; i++ {
		dict[rune(offset+i)] = i
		dictRev[i] = rune(offset + i)
	}
	offset = 65
	for i := 0; i < 26; i++ {
		dict[rune(offset+i)] = i + 26
		dictRev[i+26] = rune(offset + i)
	}
}

func getInput() ([]rune, []rune) {
	input := readInput()

	var s, plainText []rune
	for i, c := range input {
		if i == len(input)-1 {
			plainText = []rune(c)
			break
		}
		s = append(s, []rune(c)...)
	}

	return s, plainText
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
