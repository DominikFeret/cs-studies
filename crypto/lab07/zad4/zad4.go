package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	num, exp, mod := getInput()

	r := fastModExp(num, exp, mod)

	fmt.Print(r)
}

func fastModExp(num, exp, mod int) int {
	if exp == 0 {
		return 1
	}

	if exp%2 == 0 {
		return fastModExp(num*num%mod, exp/2, mod) % mod
	}

	return num * fastModExp(num, exp-1, mod) % mod
}

func getInput() (int, int, int) {
	inputData := readInput()
	s := strings.Split(inputData[0], " ")

	a, err := strconv.Atoi(s[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(s[1])
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(s[2])
	if err != nil {
		panic(err)
	}

	return a, b, n
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
