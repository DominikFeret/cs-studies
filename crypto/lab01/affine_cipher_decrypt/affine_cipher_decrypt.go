package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var dict [26]rune
var dictRev map[rune]int

const dictLen = 26

func main() {
	dictGen()
	str, a, b := getInput()
	a = modularInverse(a, dictLen)
	if a == -1 {
		fmt.Print("BŁĄD")
		return
	}

	enc := affineCipherDec(str, a, b, dictLen)

	fmt.Printf("%s\n%d %d", string(enc), a, b)
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

func getInput() ([]rune, int, int) {
	input := readInput()

	s := []rune{}
	for i, str := range input {
		if i == len(input)-1 {
			break
		}
		s = append(s, []rune(str)...)
		if i != len(input)-2 {
			s = append(s, '\n')
		}
	}

	splitStr := strings.Split(input[len(input)-1], " ")
	a, err := strconv.Atoi(splitStr[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(splitStr[1])
	if err != nil {
		panic(err)
	}

	return s, a, b
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
