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

	result := compareBits(bits1, bits2)

	fmt.Print(result)
}

func compareBits(s1, s2 string) string {
	var str string
	for i := 0; i < len(s1); i++ {
		if s1[i] == s2[i] {
			str += "_"
		} else {
			str += "X"
		}
	}

	return str
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
