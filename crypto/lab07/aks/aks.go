package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	num := getInput()

	if num == 2 {
		fmt.Print("Liczba pierwsza")
		return
	}

	res := aks(num)

	if !res {
		fmt.Print("Liczba złożona")
	} else {
		fmt.Print("Liczba pierwsza")
	}
}

func aks(num uint64) bool {
	if num%2 == 0 {
		return false
	}

	coef := num
	i := num + 1
	for j := uint64(2); j < i; j++ {
		if coef%num != 0 {
			return false
		}
		coef = coef * (i - j) / j
	}

	return true
}

func getInput() uint64 {
	inputData := readInput()
	s := strings.Split(inputData[0], " ")

	a, err := strconv.ParseUint(s[0], 10, 64)
	if err != nil {
		panic(err)
	}

	return a
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
