package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func main() {
	num := getInput()

	e, d, k := generateKeys(num)

	fmt.Printf("Klucz publiczny: %d,%d,%d\n", num, k, e)
	fmt.Printf("Klucz prywatny: %d,%d,%d", num, k, d)
}

func generateKeys(num *big.Int) (*big.Int, *big.Int, *big.Int) {
	bottomCap := new(big.Int).Lsh(big.NewInt(1), uint(num.Uint64())-1)
	topCap := new(big.Int).Lsh(big.NewInt(1), uint(num.Uint64()))

	// generate prime p and q
	p := new(big.Int)
	for {
		p = randomInt(bottomCap, topCap)
		if mayBePrime(p, 5) {
			break
		}
	}
	q := new(big.Int)
	for {
		q = randomInt(bottomCap, topCap)
		if mayBePrime(q, 5) {
			break
		}
	}

	key := new(big.Int).Mul(p, q)
	phi := phiFunc(p, q)

	e := new(big.Int)
	for {
		e = randomInt(bottomCap, topCap)
		if big.NewInt(1).GCD(nil, nil, e, phi).Cmp(big.NewInt(1)) == 0 {
			break
		}
	}
	d := new(big.Int).ModInverse(e, phi)

	return e, d, key
}

func phiFunc(p, q *big.Int) *big.Int {
	return new(big.Int).Mul(new(big.Int).Sub(p, big.NewInt(1)), new(big.Int).Sub(q, big.NewInt(1)))
}

func randomInt(bottom, top *big.Int) *big.Int {
	r, err := rand.Int(rand.Reader, new(big.Int).Sub(top, bottom))
	if err != nil {
		panic(err)
	}
	r.Add(r, bottom)

	return r
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

func getInput() *big.Int {
	inputData := readInput()
	s := strings.Split(inputData[0], " ")

	a, ok := new(big.Int).SetString(s[0], 10)
	if !ok {
		panic("Failed to parse input")
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
