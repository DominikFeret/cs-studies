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
	str, pt := getInput()
	dec, a, b := analyseText(str, pt)

	fmt.Printf("%s\n%d %d", string(dec), a, b)
}

func analyseText(encStr, decStr []rune) ([]rune, int, int) {
	indexes := []int{}
	for j, c := range decStr {
		if !isLiteral(c) {
			continue
		}
		indexes = append(indexes, j)
	}

	encSample := []rune{}
	decSample := []rune{}

	for _, i := range indexes {
		encSample = append(encSample, encStr[i])
		decSample = append(decSample, decStr[i])
	}

	if len(encSample) < 2 {
		return encStr, 1, 0
	}
	if len(encSample) == 2 {
		a, b := affineCipherBreak(encSample, decSample)
		s := affineCipherDec(encStr, modularInverse(a, dictLen), b, dictLen)
		return s, a, b
	}

	possibleAs := []int{}
	possibleBs := []int{}
	a, b := affineCipherBreak(encSample[:2], decSample[:2])
	possibleAs = append(possibleAs, a)
	possibleBs = append(possibleBs, b)
	a, b = affineCipherBreak(encSample[1:], decSample[1:])
	possibleAs = append(possibleAs, a)
	possibleBs = append(possibleBs, b)
	a, b = affineCipherBreak([]rune{encSample[0], encSample[2]}, []rune{decSample[0], decSample[2]})
	possibleAs = append(possibleAs, a)
	possibleBs = append(possibleBs, b)
	if possibleAs[0] == possibleAs[1] || possibleAs[0] == possibleAs[2] {
		a = possibleAs[0]
		b = possibleBs[0]
	} else {
		a = possibleAs[1]
		b = possibleBs[1]
	}

	s := affineCipherDec(encStr, modularInverse(a, dictLen), b, dictLen)

	return s, a, b
}

func affineCipherBreak(str, pt []rune) (a, b int) {
	decChar1 := dictRev[toLower(pt[0])]
	decChar2 := dictRev[toLower(pt[1])]
	encChar1 := dictRev[toLower(str[0])]
	encChar2 := dictRev[toLower(str[1])]

	x := mod(decChar1-decChar2, dictLen)
	y := mod(encChar1-encChar2, dictLen)
	if x == 0 {
		x = 1
	}

	for y%x != 0 {
		y += dictLen
		a++
	}
	a = y / x

	b = mod(encChar1-(a*decChar1), dictLen)

	return a, b
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

func getInput() ([]rune, []rune) {
	input := readInput()

	s := []rune(input[0])
	ps := []rune(input[1])

	return s, ps
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
