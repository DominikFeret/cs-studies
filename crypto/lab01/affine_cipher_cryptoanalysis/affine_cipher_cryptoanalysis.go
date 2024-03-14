package main

import (
	"bufio"
	"fmt"
	"os"
)

var dict [26]rune
var dictRev map[rune]int

const dictLen = 26

func main() {
	dictGen()
	str := getInput()
	keyArray := getKeys()
	offsetArray := getOffsets()

	strOut := ""
	for i, a := range keyArray {
		for j, b := range offsetArray {
			enc := affineCipherDec(str, modularInverse(a, dictLen), b, dictLen)
			strOut += fmt.Sprintf("A=%d B=%d %s", a, b, string(enc))
			if i != len(keyArray)-1 || j != len(offsetArray)-1 {
				strOut += "\n"
			}
		}
	}

	fmt.Printf("%s", strOut)
}

func affineCipherDec(str []rune, a, b, z int) []rune {
	enc := ""
	for _, c := range str {
		if !isLiteral(c) {
			enc += string(c)
			continue
		}

		isUppercase := isUppercase(c)
		c = toLower(c) - 97

		encChar := dict[mod(a*(int(c)-b), z)]

		if isUppercase {
			encChar = toUpper(encChar)
		}
		enc += string(encChar)
	}

	return []rune(enc)
}

func getKeys() []int {
	keys := []int{}
	for i := 1; i < dictLen; i++ {
		if modularInverse(i, dictLen) != -1 {
			keys = append(keys, i)
		}
	}

	return keys
}

func getOffsets() []int {
	offsets := []int{}
	for i := range dictLen {
		offsets = append(offsets, i)
	}

	return offsets
}

func isLiteral(c rune) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func isUppercase(c rune) bool {
	return (c >= 65 && c <= 90)
}

func toLower(c rune) rune {
	if !isUppercase(c) {
		return c
	}
	return c + 32
}

func toUpper(c rune) rune {
	if isUppercase(c) {
		return c
	}
	return c - 32
}

func modularInverse(num int, mod int) int {
	s := 0
	oldS := 1
	r := mod
	newR := num

	for newR != 0 {
		quotient := r / newR

		s, oldS = oldS, s-quotient*oldS
		r, newR = newR, r-quotient*newR
	}

	if r != 1 {
		return -1
	}

	if s < 0 {
		s += mod
	}

	return s
}

func mod(num int, z int) int {
	num = num % z
	if num < 0 {
		num += z
	}
	return num
}

func dictGen() {
	dictRev = make(map[rune]int)
	for i := 97; i <= 122; i++ {
		dict[i-97] = rune(i)
		dictRev[rune(i)] = i - 97
	}
}

func getInput() []rune {
	input := readInput()

	s := input[0]

	return []rune(s)
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
