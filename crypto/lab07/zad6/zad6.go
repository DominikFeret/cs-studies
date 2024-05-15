package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	num := getInput()
	num--
	exp := 0
	for num%2 == 0 {
		num /= 2
		exp++
	}

	fmt.Printf("2^%d * %d", exp, num)
}

func getInput() int {
	inputData := readInput()
	s := inputData[0]

	a, err := strconv.Atoi(s)
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
