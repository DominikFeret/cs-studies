package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	num, reps, _ := getInput()

	if num.Cmp(big.NewInt(2)) == 0 {
		fmt.Print("Liczba jest prawdopodobnie pierwsza")
		return
	}

	if num.Cmp(big.NewInt(2)) == -1 {
		fmt.Print("Liczba złożona")
	}

	r := mayBePrime(num, reps)

	if r == false {
		fmt.Print("Liczba złożona:")
		fs := getPrimeFactors(num)

		for i, n := range fs {

			fmt.Printf(" %d", n)
			if i != len(fs)-1 {
				fmt.Print(" *")
			}
		}
	} else {
		fmt.Print("Liczba jest prawdopodobnie pierwsza")
	}
}

func getPrimeFactors(num *big.Int) []*big.Int {
	cap := new(big.Int).Div(num, big.NewInt(2))
	numCpy := new(big.Int).Set(num)

	factors := []*big.Int{}
	for i := big.NewInt(2); i.Cmp(cap) == -1; i.Add(i, big.NewInt(1)) {
		if new(big.Int).Mod(numCpy, i).Cmp(big.NewInt(0)) == 0 { // num%i == 0
			iCpy := new(big.Int).Set(i)
			factors = append(factors, iCpy)
			numCpy.Div(numCpy, i)
			cap.Div(cap, i)
			i.Sub(i, big.NewInt(1))
			if mayBePrime(numCpy, 10) == true {
				factors = append(factors, numCpy)
				break
			}
		}
	}

	return factors
}

func mayBePrime(num *big.Int, reps uint64) bool {
	for range reps {
		if rabinMiller(num) == false {
			return false
		}
	}

	return true
}

func rabinMiller(num *big.Int) bool {
	if tmp := big.NewInt(0); tmp.Mod(num, big.NewInt(2)) == big.NewInt(0) {
		return false
	}

	// second step
	a, err := rand.Int(rand.Reader, new(big.Int).Sub(num, big.NewInt(2)))
	if err != nil {
		panic(err)
	}
	a.Add(a, big.NewInt(2))
	if big.NewInt(1).GCD(nil, nil, num, a).Cmp(big.NewInt(1)) != 0 { // gcd(num, a) == 1
		return false
	}

	// third step
	k, m := decomposeNumber(big.NewInt(1).Set(num).Sub(num, big.NewInt(1)))

	if new(big.Int).Exp(a, m, num).Cmp(big.NewInt(1)) == 0 { // fastModExp(a, m, num) == 1
		return true
	}

	for i := uint(0); k.Cmp(big.NewInt(int64(i))) == -1; i++ {
		pwr := big.NewInt(0).Mul(m, new(big.Int).Lsh(big.NewInt(2), i)) // pwr := m * (2 << (i))

		if new(big.Int).Exp(a, pwr, num).Cmp(big.NewInt(1)) == 1 { // fastModExp(a, pwr, num) == 1
			if i == 0 {
				pwr = m
			} else {
				pwr.Mul(m, new(big.Int).Lsh(big.NewInt(2), i-1)) // pwr = m * (2 << (i - 1))
			}

			if fastModExp(a, pwr, num) != new(big.Int).Sub(num, big.NewInt(1)) {
				return false
			}
			return true
		}
	}

	return false
}

func decomposeNumber(num *big.Int) (*big.Int, *big.Int) {
	exp := big.NewInt(0)
	for numCpy := new(big.Int).Set(num); big.NewInt(1).Mod(numCpy, big.NewInt(2)) == big.NewInt(0); {
		numCpy.Div(numCpy, big.NewInt(2))
		exp.Add(exp, big.NewInt(1))
	}

	return exp, num
}

func fastModExp(num, exp, mod *big.Int) *big.Int {
	acc := big.NewInt(1)
	numCpy := new(big.Int).Set(num)
	expCpy := new(big.Int).Set(exp)
	for expCpy.Cmp(big.NewInt(0)) == 1 {
		if tmp := big.NewInt(0); tmp.Mod(expCpy, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 { // exp%2 == 1
			acc.Mul(acc, numCpy).Mod(acc, mod) // acc = (acc * num) % mod
		}

		numCpy.Mul(numCpy, numCpy).Mod(numCpy, mod) // num = (num * num) % mod
		expCpy.Div(expCpy, big.NewInt(2))           // exp = exp / 2
	}

	return acc
}

func gcd(a, b uint64) uint64 {
	for a != 0 {
		a, b = b%a, a
	}

	return b
}

func getInput() (*big.Int, uint64, uint64) {
	inputData := readInput()
	s := strings.Split(inputData[0], " ")

	a, ok := new(big.Int).SetString(s[0], 10)
	if !ok {
		panic("Failed to parse input")
	}

	b, err := strconv.ParseUint(s[1], 10, 64)
	if err != nil {
		panic(err)
	}
	n, err := strconv.ParseUint(s[2], 10, 64)
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
