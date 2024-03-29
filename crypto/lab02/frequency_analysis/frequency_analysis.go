package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type freqDist struct {
	char      rune
	frequency int
}

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 26

func main() {
	initDict()
	str := getInput()

	strFreq := analyzeFrequency(str)

	fmt.Printf("%d", strFreq)
}

func analyzeFrequency(str []rune) int {
	freqDist := getFreqDist(str)

	mostFreq := []rune{'e', 't', 'a', 'o', 'i', 'n'}
	leastFreq := []rune{'z', 'q', 'x', 'j', 'k', 'v'}

	freqIndex := 0
	for i := 0; i < 6; i++ {
		c := freqDist[i].char
		if slices.Contains(mostFreq, c) {
			freqIndex++
		}
	}

	for i := len(freqDist) - 6; i < len(freqDist); i++ {
		c := freqDist[i].char
		if slices.Contains(leastFreq, c) {
			freqIndex++
		}
	}

	return freqIndex
}

func getFreqDist(str []rune) []freqDist {
	freqDist := initFreqDist()

	for _, c := range str {
		if !isLiteral(c) {
			continue
		}
		c = toLower(c)

		freqDist[dict[c]].frequency++
	}

	quickSort(freqDist, 0, len(freqDist)-1)

	return freqDist
}

func isLiteral(c rune) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func isUppercase(c rune) bool {
	return (c >= 65 && c <= 90)
}

func toLower(c rune) rune {
	if !isUppercase(c) {
		return c
	}
	return c + 32
}

func toUpper(c rune) rune {
	if isUppercase(c) {
		return c
	}
	return c - 32
}

func mod(num int, z int) int {
	num = num % z
	if num < 0 {
		num += z
	}
	return num
}

func quickSort(arr []freqDist, low int, high int) {
	if low < high {
		pivot := partition(arr, low, high)

		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

func partition(arr []freqDist, low int, high int) int {
	pivot := arr[high].frequency
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j].frequency > pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]

	return i + 1
}

func initFreqDist() []freqDist {
	freqDist := make([]freqDist, dictLen)

	for i := 0; i < dictLen; i++ {
		freqDist[i].char = dictRev[i]
		freqDist[i].frequency = 0
	}

	return freqDist
}

func initDict() {
	dict = make(map[rune]int)
	dictRev = make(map[int]rune)
	for i := 97; i <= 122; i++ {
		dict[rune(i)] = i - 97
		dictRev[i-97] = rune(i)
	}
}

func getInput() []rune {
	input := readInput()

	var str []rune
	for i, c := range input {
		str = append(str, []rune(c)...)
		if i != len(input)-1 {
			str = append(str, '\n')
		}
	}

	return str
}

func readInput() []string {
	r := bufio.NewScanner(os.Stdin)
	s := []string{}

	for r.Scan() {
		s = append(s, r.Text())
	}

	return s
}
