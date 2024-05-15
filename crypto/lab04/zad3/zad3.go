package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	str := getInput()
	s1 := [][]byte{
		[]byte("101"),
		[]byte("010"),
		[]byte("001"),
		[]byte("110"),
		[]byte("011"),
		[]byte("100"),
		[]byte("111"),
		[]byte("000"),
		[]byte("001"),
		[]byte("100"),
		[]byte("110"),
		[]byte("010"),
		[]byte("000"),
		[]byte("111"),
		[]byte("101"),
		[]byte("011"),
	}
	s2 := [][]byte{
		[]byte("100"),
		[]byte("000"),
		[]byte("110"),
		[]byte("101"),
		[]byte("111"),
		[]byte("001"),
		[]byte("011"),
		[]byte("010"),
		[]byte("101"),
		[]byte("011"),
		[]byte("000"),
		[]byte("111"),
		[]byte("110"),
		[]byte("010"),
		[]byte("001"),
		[]byte("100"),
	}
	perm := []int{
		0, 1, 3, 2, 3, 2, 4, 5,
	}
	key := []byte("10101010")

	enc := miniDesEnc(str, s1, s2, perm, key)

	fmt.Print(enc)
}

func miniDesEnc(str []byte, s1, s2 [][]byte, perm []int, key []byte) string {
	bLen := len(str) / 2
	pl, pr := binToInt(str[:bLen]), binToInt(str[bLen:])
	cl, cr := 0, 0

	for i := 0; i < 7; i++ {
		// shift key slice left by 1 with each iteration
		e := permute(intToBin6(pr), perm)
		shiftLeft(key)

		xored := intToBin8(binToInt(e) ^ binToInt(key))

		i1, i2 := binToInt(xored[:4]), binToInt(xored[4:])

		f := binToInt([]byte(string(s1[i1]) + string(s2[i2])))

		cl = pr
		cr = pl ^ f
		pl, pr = cl, cr
	}

	return string(append(intToBin6(cr), intToBin6(cl)...))
}

func binToInt(bin []byte) int {
	s := string(bin)

	n := 0

	for i := 0; i < len(s); i++ {
		if s[len(s)-i-1] == '1' {
			n += 1 << i
		}
	}

	return int(n)
}

func intToBin6(n int) []byte {
	s := fmt.Sprintf("%06b", n)

	return []byte(s)
}

func intToBin8(n int) []byte {
	s := fmt.Sprintf("%08b", n)

	return []byte(s)
}

func permute(bin []byte, perm []int) []byte {
	permuted := make([]byte, len(perm))
	for i, v := range perm {
		permuted[i] = bin[v]
	}

	return permuted
}

func shiftLeft(key []byte) {
	first := key[0]
	for i := 0; i < len(key)-1; i++ {
		key[i] = key[i+1]
	}
	key[len(key)-1] = first
}

func getInput() []byte {
	input := readInput()

	return input[0]
}

func readInput() [][]byte {
	input := make([][]byte, 0)
	for r := bufio.NewReader(os.Stdin); ; {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		input = append(input, line)
	}

	return input
}
