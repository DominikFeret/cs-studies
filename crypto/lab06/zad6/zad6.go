package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	a, b := getInput()

	x, y, d := extendedEuclid(a, b)

	fmt.Printf("%d %d %d", x, y, d)
}

func extendedEuclid(a, b int) (x, y, d int) {
	x_a, x_b := 1, 0
	y_a, y_b := 0, 1
	c := 0

	for a*b != 0 {
		if a >= b {
			c = a / b
			a %= b
			x_a -= c * x_b
			y_a -= c * y_b
		} else {
			c = b / a
			b %= a
			x_b -= c * x_a
			y_b -= c * y_a
		}
	}

	if a > 0 {
		return x_a, y_a, a
	}

	return x_b, y_b, b
}

func getInput() (int, int) {
	s := strings.Split(readInput()[0], " ")
	a, err := strconv.Atoi(s[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(s[1])
	if err != nil {
		panic(err)
	}

	return a, b
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
