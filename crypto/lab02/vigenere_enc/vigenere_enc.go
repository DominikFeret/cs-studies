package main

import (
	"bufio"
	"fmt"
	"os"
)

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 26

func main() {
	dictGen()
	str, key := getInput()

	cryptogram := vigenereEncrypt(str, key)

	fmt.Printf("%s", string(cryptogram))
}

func vigenereEncrypt(str []rune, key []int) []rune {
	enc := ""
	j := 0
	for _, c := range str {
		if !isLiteral(c) {
			enc += string(c)
			continue
		}

		isUpper := isUppercase(c)
		c = toLower(c)
		encChar := dictRev[mod(int(dict[c])+key[j], dictLen)]
		if isUpper {
			encChar = toUpper(encChar)
		}
		enc += string(encChar)
		j++
		j %= len(key)
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
	dict = make(map[rune]int)
	dictRev = make(map[int]rune)
	for i := 97; i <= 122; i++ {
		dict[rune(i)] = i - 97
		dictRev[i-97] = rune(i)
	}
}

func getInput() ([]rune, []int) {
	input := readInput()

	var str, keyStr []rune
	for i, c := range input {
		if i == len(input)-1 {
			keyStr = []rune(c)
			break
		}
		str = append(str, []rune(c)...)
	}

	keys := make([]int, len(keyStr))
	for i, c := range keyStr {
		keys[i] = dict[toLower(c)]
	}

	return str, keys
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
