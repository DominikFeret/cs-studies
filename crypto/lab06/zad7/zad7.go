package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	a := getInput()

	num, coprimes := phi(a)

	fmt.Printf("%d\n", num)
	for i, v := range coprimes {
		if i == len(coprimes)-1 {
			fmt.Printf("%d", v)
			break
		}
		fmt.Printf("%d ", v)
	}
}

func gcd(a, b int) int {
	if a == 0 {
		return b
	}

	return gcd(b%a, a)
}

func phi(a int) (int, []int) {
	coprimes := []int{}
	for i := 1; i < a; i++ {
		if gcd(i, a) != 1 {
			continue
		}

		coprimes = append(coprimes, i)
	}

	return len(coprimes), coprimes
}

func getInput() int {
	s := strings.Split(readInput()[0], " ")
	a, err := strconv.Atoi(s[0])
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
