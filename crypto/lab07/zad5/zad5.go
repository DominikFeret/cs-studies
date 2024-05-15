package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	num, reps, s := getInput()
	rand.NewSource(int64(s))

	if num == 2 {
		fmt.Print("Liczba prawdopodobnie pierwsza")
		return
	}

	r := mayBePrime(num, reps)

	if r == false {
		fmt.Print("Liczba złożona")
	} else {
		fmt.Print("Liczba prawdopodobnie pierwsza")
	}
}

func mayBePrime(num, reps int) bool {
	for _ = range reps {
		if fermat(num) == false {
			return false
		}
	}

	return true
}

func fermat(num int) bool {
	r := rand.Intn(num-2) + 2

	if gcd(r, num) != 1 {
		return false
	}

	if fastModExp(r, num-1, num) != 1 {
		return false
	}

	return true
}

func fastModExp(num, exp, mod int) int {
	acc := 1
	for exp != 0 {
		if exp%2 == 1 {
			acc = (acc * num) % mod
		}

		num = (num * num) % mod
		exp /= 2
	}

	return acc
}

func gcd(a, b int) int {
	for a != 0 {
		a, b = b%a, a
	}

	return b
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
