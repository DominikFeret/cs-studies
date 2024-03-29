package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type freqDist struct {
	number    int
	frequency int
}

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 26

func main() {
	initDict()
	str := getInput()
	peeledStr := peelString(str)
	distances := findSequenceDistances(peeledStr)
	fmt.Println(distances)
	possibleLengths := getPossibleLenghts(distances, peeledStr)

	for i, val := range possibleLengths {
		if i == len(possibleLengths)-1 {
			fmt.Printf("%d", val.number)
			continue
		}
		fmt.Printf("%d ", val.number)
	}
}

func findSequenceDistances(str []rune) []int {
	seqLen := 3
	possibleLengths := make([]int, 0)

	for i := 0; i < len(str)-seqLen; i++ {
		seq := str[i : i+seqLen]
		for j := i + 1; j < len(str)-seqLen; j++ {
			consideredSeq := str[j : j+seqLen]
			if string(seq) == string(consideredSeq) && !slices.Contains(possibleLengths, j-i) {
				possibleLengths = append(possibleLengths, j-i)
			}
		}
	}

	return possibleLengths
}

func getPossibleLenghts(distances []int, str []rune) []freqDist {
	factorBuckets := make([]int, len(str)-3)

	for _, val := range distances {
		for i := 2; i <= val/2; i++ {
			if val%i == 0 {
				factorBuckets[i]++
			}
		}
		factorBuckets[val]++
	}

	possibleLengths := analyseBuckets(factorBuckets)
	trimmedPossibleLengths := make([]freqDist, 0)
	for _, val := range possibleLengths {
		if val.frequency > 0 && val.number < 17 {
			trimmedPossibleLengths = append(trimmedPossibleLengths, val)
		}
	}

	quickSort(possibleLengths, 0, len(possibleLengths)-1)

	return trimmedPossibleLengths
}

func analyseBuckets(buckets []int) []freqDist {
	possibleLengths := make([]freqDist, 0)

	for i, val := range buckets {
		if val > 0 {
			possibleLengths = append(possibleLengths, freqDist{i, val})
		}
	}

	return possibleLengths
}

func peelString(str []rune) []rune {
	peeledStr := make([]rune, 0)
	for _, c := range str {
		if isLiteral(c) {
			peeledStr = append(peeledStr, toLower(c))
		}
	}
	return peeledStr
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
		if arr[j].frequency == pivot {
			if arr[j].number < arr[high].number {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
		if arr[j].frequency > pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]

	return i + 1
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
