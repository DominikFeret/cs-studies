package main

import (
	"bufio"
	"crypto"
	"fmt"
	"os"

	_ "crypto/md5"
)

func main() {
	s1, s2 := getInput()

	bits1 := hashToBits(s1, crypto.MD5)
	bits2 := hashToBits(s2, crypto.MD5)

	x := diffBits(bits1, bits2)
	result := (float64(x) / float64(len(bits1))) * 100

	p := fmt.Sprintf("%.2f", result)
	fmt.Print(p + "%")
}

func diffBits(s1, s2 string) int {
	var n int
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			n++
		}
	}

	return n
}

func hashToBits(s []byte, hf crypto.Hash) string {
	h := hf.New()
	h.Write(s)
	hs := h.Sum(nil)

	var strOut string
	for _, b := range hs {
		strOut += fmt.Sprintf("%08b", b)
	}

	return strOut
}

func getInput() ([]byte, []byte) {
	i := readInput()
	s1 := []byte(i[0])
	s2 := []byte(i[1])

	return s1, s2
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
