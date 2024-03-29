package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
)

type freqKey interface {
	rune | int
}

type freqDist[E freqKey] struct {
	number    E
	frequency int
}

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 26

func main() {
	initDict()
	str := getInput()
	peeledStr := peelString(str)
	engDict := getEnglishDict()

	distances := findSequenceDistances(peeledStr)
	possibleLengths := getPossibleLenghts(distances, peeledStr)

	strDec, key := cryptoanalysis(engDict, possibleLengths, str)

	fmt.Printf("%s\n%s", string(strDec), string(key))
}

func cryptoanalysis(engDict [][]rune, possibleLengths []freqDist[int], str []rune) ([]rune, []rune) {
	peeledStr := peelString(str)
	for _, val := range possibleLengths {
		fmt.Println("checking length: ", val.number)
		partitionedStr := partitionString(peeledStr, val.number)

		mostFreqKeys := getMostFreqKeys(partitionedStr)

		possibleFullKeys := getPossibleFullKeys(mostFreqKeys)

		peeledDict := peelDictByLength(engDict, val.number)

		possibleKeys := dictionaryAttackKeys(peeledDict, possibleFullKeys)

		for _, key := range possibleKeys {
			strDec := dictionaryAttackCryptogram(engDict, str, key)
			if strDec != nil {
				return strDec, key
			}
		}
	}

	return nil, nil
}

func dictionaryAttackCryptogram(engDict [][]rune, str []rune, key []rune) []rune {
	keyAsInt := make([]int, 0)
	for _, c := range key {
		keyAsInt = append(keyAsInt, dict[c])
	}

	strDec := vigenereDecrypt(str, keyAsInt)

	strPeeled := []rune(regexp.MustCompile(`[^a-zA-Z \n]+`).ReplaceAllString(string(strDec), ""))

	strSeparated := strings.Fields(string(strPeeled))

	numOfTries := 20
	expectedAccuracy := 0.2
	numOfHits := 0
	for i := 0; i < numOfTries; i++ {
		r := rand.Intn(len(strSeparated))
		randWord := strSeparated[r]

		for _, word := range engDict {
			if len(word) != len(randWord) {
				continue
			}
			if string(word) == randWord {
				numOfHits++
			}
		}
	}

	if float64(numOfHits)/float64(numOfTries) >= expectedAccuracy {
		return strDec
	}

	return nil
}

func dictionaryAttackKeys(engDict [][]rune, possibleFullKeys [][]rune) [][]rune {
	possibleKeys := make([][]rune, 0)

	for _, key := range possibleFullKeys {
		matched := false
		for _, word := range engDict {
			if string(key) == string(word) {
				matched = true
				break
			}
		}

		if matched {
			possibleKeys = append(possibleKeys, []rune{})
			i := len(possibleKeys) - 1
			possibleKeys[i] = append(possibleKeys[i], key...)
		}

	}

	return possibleKeys
}

func peelDictByLength(engDict [][]rune, l int) [][]rune {
	peeledDict := make([][]rune, 0)
	for _, word := range engDict {
		if len(word) == l {
			peeledDict = append(peeledDict, word)
		}
	}
	return peeledDict
}

func getPossibleFullKeys(possibleSubkeys [][]rune) [][]rune {
	size := 1
	for _, val := range possibleSubkeys {
		size *= len(val)
	}

	possibleFullKeys := make([][]rune, size)
	var i int
	generateCombinations(possibleSubkeys, 0, []rune{}, possibleFullKeys, &i)

	return possibleFullKeys
}

func generateCombinations(possibleSubkeys [][]rune, index int, currentComb []rune, possibleFullKeys [][]rune, count *int) {
	if index == len(possibleSubkeys) {
		possibleFullKeys[*count] = append(possibleFullKeys[*count], currentComb...)
		*count++
		return
	}

	for _, val := range possibleSubkeys[index] {
		generateCombinations(possibleSubkeys, index+1, append(currentComb, val), possibleFullKeys, count)
	}

	return
}

func getMostFreqKeys(partitionedStr [][]rune) [][]rune {
	mostFreqKeys := make([][]rune, len(partitionedStr))

	for i, partition := range partitionedStr {
		frequencies := make([]int, dictLen)
		maxFreq := 0
		for j := range frequencies {
			strDec := make([]rune, len(partition))
			for k, c := range partition {
				dec := mod(dict[c]-j, dictLen)
				strDec[k] = dictRev[dec]
			}
			freq := analyzeFrequency(strDec)
			frequencies[j] = freq
			if freq > maxFreq {
				maxFreq = freq
			}
		}

		for j, val := range frequencies {
			// IMPORTANT NOTE - in case of shorter texts, the strictness of the condition can bo lowered,
			// as maxFreq has a higher chance of being a false positive in such cases
			// otherwise the condition should be kept as is for performance reasons
			if maxFreq == val {
				mostFreqKeys[i] = append(mostFreqKeys[i], dictRev[j])
			}
		}
	}

	return mostFreqKeys
}

func partitionString(str []rune, distance int) [][]rune {
	partitionedStr := make([][]rune, 0)

	for i := 0; i < distance; i++ {
		partitionedStr = append(partitionedStr, make([]rune, 0))
		for j := i; j < len(str); j += distance {
			partitionedStr[i] = append(partitionedStr[i], str[j])
		}
	}

	return partitionedStr
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

func getPossibleLenghts(distances []int, str []rune) []freqDist[int] {
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
	trimmedPossibleLengths := make([]freqDist[int], 0)
	for _, val := range possibleLengths {
		if val.frequency > 0 && val.number < 17 {
			trimmedPossibleLengths = append(trimmedPossibleLengths, val)
		}
	}

	quickSort(possibleLengths, 0, len(possibleLengths)-1)

	return trimmedPossibleLengths
}

func analyzeFrequency(str []rune) int {
	freqDist := getFreqDist(str)

	mostFreq := []rune{'e', 't', 'a', 'o', 'i', 'n'}
	leastFreq := []rune{'z', 'q', 'x', 'j', 'k', 'v'}

	freqIndex := 0
	for i := 0; i < 6; i++ {
		c := freqDist[i].number
		if slices.Contains(mostFreq, c) {
			freqIndex++
		}
	}

	for i := len(freqDist) - 6; i < len(freqDist); i++ {
		c := freqDist[i].number
		if slices.Contains(leastFreq, c) {
			freqIndex++
		}
	}

	return freqIndex
}

func vigenereDecrypt(str []rune, keys []int) []rune {
	dec := ""
	j := 0
	for _, c := range str {
		if !isLetter(c) {
			dec += string(c)
			continue
		}

		isUpper := isUppercase(c)
		c = toLower(c)

		decChar := dictRev[mod(int(dict[c])-keys[j], dictLen)]
		if isUpper {
			decChar = toUpper(decChar)
		}
		dec += string(decChar)
		j++
		j %= len(keys)
	}

	return []rune(dec)
}

func getEnglishDict() [][]rune {
	engDict, err := http.Get("https://stepik.org/media/attachments/lesson/668860/dictionary.txt")
	if err != nil {
		panic(err)
	}
	defer engDict.Body.Close()
	engDictStr := make([][]rune, 0)
	scanner := bufio.NewScanner(engDict.Body)
	for scanner.Scan() {
		engDictStr = append(engDictStr, []rune(scanner.Text()))
	}

	for i, word := range engDictStr {
		for j, c := range word {
			engDictStr[i][j] = toLower(c)
		}
	}

	return engDictStr
}

func getFreqDist(str []rune) []freqDist[rune] {
	freqDist := initFreqDist()

	for _, c := range str {
		if !isLetter(c) {
			continue
		}
		c = toLower(c)

		freqDist[dict[c]].frequency++
	}

	quickSort(freqDist, 0, len(freqDist)-1)

	return freqDist
}

func initFreqDist() []freqDist[rune] {
	freqDist := make([]freqDist[rune], dictLen)

	for i := 0; i < dictLen; i++ {
		freqDist[i].number = dictRev[i]
		freqDist[i].frequency = 0
	}

	return freqDist
}

func analyseBuckets(buckets []int) []freqDist[int] {
	possibleLengths := make([]freqDist[int], 0)

	for i, val := range buckets {
		if val > 0 {
			possibleLengths = append(possibleLengths, freqDist[int]{i, val})
		}
	}

	return possibleLengths
}

func peelString(str []rune) []rune {
	peeledStr := make([]rune, 0)
	for _, c := range str {
		if isLetter(c) {
			peeledStr = append(peeledStr, toLower(c))
		}
	}
	return peeledStr
}

func isLetter(c rune) bool {
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

func quickSort[E freqKey](arr []freqDist[E], low int, high int) {
	if low < high {
		pivot := partition(arr, low, high)

		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

func partition[E freqKey](arr []freqDist[E], low int, high int) int {
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
