package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	action, mods, input := getInput()

	out := []int{}
	if action == "szyfruj" {
		out = encrypt(mods, input[0])
	} else if action == "deszyfruj" {
		out = decrypt(mods, input)
	} else {
		panic("Invalid action")
	}

	for i, v := range out {
		if i == len(out)-1 {
			print(v)
		} else {
			print(v, " ")
		}
	}
}

func encrypt(mods []int, data int) []int {
	out := []int{}

	for _, m := range mods {
		out = append(out, data%m)
	}

	return out
}

func decrypt(mods []int, data []int) []int {
	totalM := func() int {
		total := 1
		for _, v := range mods {
			total *= v
		}
		return total
	}()

	ms := make([]int, len(mods))
	ns := make([]int, len(mods))

	for i, m := range mods {
		ms[i] = totalM / m
		ns[i], _, _ = extendedEuclid(ms[i], m) // ns is the inverse of ms[i] modulo m
	}

	a := func() int {
		total := 0
		for i, d := range data {
			total += d * ms[i] * ns[i]
		}
		return mod(total, totalM)
	}()

	return []int{a}
}

func mod(a, b int) int {
	if a < 0 {
		return b + a%b
	}

	return a % b
}

func extendedEuclid(a, b int) (x, y, d int) {
	xa, xb := 1, 0
	ya, yb := 0, 1
	c := 0

	for a*b != 0 {
		if a >= b {
			c = a / b
			a %= b
			xa -= c * xb
			ya -= c * yb
		} else {
			c = b / a
			b %= a
			xb -= c * xa
			yb -= c * ya
		}
	}

	if a > 0 {
		return xa, ya, a
	}

	return xb, yb, b
}

func getInput() (string, []int, []int) {
	inputData := readInput()
	action := inputData[0]

	s := strings.Split(inputData[1], " ")
	mods := []int{}
	for _, v := range s {
		a, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		mods = append(mods, a)
	}

	s = strings.Split(inputData[2], " ")

	data := []int{}
	for _, v := range s {
		a, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		data = append(data, a)
	}

	return action, mods, data
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
