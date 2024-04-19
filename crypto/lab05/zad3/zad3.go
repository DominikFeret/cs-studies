package main

import (
	"bufio"
	"crypto"
	"fmt"
	"os"

	_ "crypto/md5"
)

func main() {
	h := crypto.MD5.New()
	s := getInput()

	h.Write(s)
	hs := h.Sum(nil)

	fmt.Printf("%x", hs)
}

func getInput() []byte {
	i := readInput()
	s := []byte(i[0])

	return s
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
