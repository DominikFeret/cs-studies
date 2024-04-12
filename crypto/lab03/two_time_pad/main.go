package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"unicode"
)

type freqKey interface {
	rune | int
}

type freqDist[E freqKey] struct {
	key       E
	frequency int
}

var dict map[rune]int
var dictRev map[int]rune

const dictLen = 26

func main() {
	str := getInput()
	engDict := getEnglishDict()

	strDec, key := cryptoanalysis(engDict, str, decryptCharStd)
	str1 := strDec[:len(strDec)/2]
	str2 := strDec[len(strDec)/2:]

	fmt.Printf("%s\n%s\n%s\n", string(str1), string(str2), string(key))
}

func cryptoanalysis(engDict map[string]bool, encStr []rune, decrypt func(c rune, k int) rune) ([]rune, []rune) {
	initDict()

	mostFreqSubkeys := findMostFreqSubkeys(encStr, decrypt)

	possibleFullKeys := getPossibleFullKeys(mostFreqSubkeys)

	possibleWordKeys := dictionaryAttackKeys(engDict, possibleFullKeys)

	for _, key := range possibleWordKeys {
		decStr := dictionaryAttackCryptogram(engDict, encStr, key, decrypt)

		if decStr != nil {
			return decStr, key
		}
	}

	// if no match found, try all possible keys in case the key is not a word
	for _, key := range possibleFullKeys {
		decStr := dictionaryAttackCryptogram(engDict, encStr, key, decrypt)
		if decStr != nil {
			return decStr, key
		}
	}

	return nil, nil
}

func dictionaryAttackCryptogram(engDict map[string]bool, encStr []rune, key []rune, decrypt func(c rune, k int) rune) []rune {
	firstWord := make([]rune, len(key))
	for i := 0; i < len(key); i++ {
		firstWord[i] = decrypt(encStr[i], dict[key[i]])
	}

	secondWord := make([]rune, len(key))
	for i := 0; i < len(key); i++ {
		secondWord[i] = decrypt(encStr[i+len(key)], dict[key[i]])
	}

	if _, ok := engDict[string(firstWord)]; ok {
		if _, ok := engDict[string(secondWord)]; ok {
			return append(firstWord, secondWord...)
		}
	}

	return nil
}

func dictionaryAttackKeys(engDict map[string]bool, possibleFullKeys [][]rune) [][]rune {
	possibleKeys := make([][]rune, 0)

	for _, key := range possibleFullKeys {
		if _, ok := engDict[string(key)]; ok {
			possibleKeys = append(possibleKeys, []rune{})
			i := len(possibleKeys) - 1
			possibleKeys[i] = append(possibleKeys[i], key...)
		}

	}

	return possibleKeys
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
}

func findMostFreqSubkeys(str []rune, decrypt func(c rune, k int) rune) [][]rune {
	partitionedStr := partitionString(str, len(str)/2)

	// get candidate subkeys for each partition
	probableSubkeys := make([][]rune, len(partitionedStr))
	for pi, partition := range partitionedStr {
		keysFreqs := make([]freqDist[rune], dictLen)
		for k := range dictLen {
			partitiondecrypted := make([]rune, len(partition))

			for j, c := range partition {
				partitiondecrypted[j] = decrypt(c, k)
			}
			keysFreqs[k].key = dictRev[k]
			keysFreqs[k].frequency = analyzeFrequency(partitiondecrypted)
		}

		slices.SortFunc(keysFreqs, compareFreqDistDesc[rune])
		probableSubkeys[pi] = reduceKeys(keysFreqs)
	}

	return probableSubkeys
}

func reduceKeys(keysFreqs []freqDist[rune]) []rune {
	// get top 6 keys + the ones that have max frequency
	topKeys := make([]rune, 0)
	for i := 0; i < dictLen; i++ {
		if i < 6 || keysFreqs[0].frequency == keysFreqs[i].frequency {
			topKeys = append(topKeys, keysFreqs[i].key)
		}
	}

	return topKeys
}

func analyzeFrequency(str []rune) int {
	freqDist := getFreqDist(str)

	mostFreq := []rune{'e', 't', 'a', 'o', 'i', 'n'}
	leastFreq := []rune{'z', 'q', 'x', 'j', 'k', 'v'}

	frequency := 0
	for i := 0; i < len(mostFreq)-3; i++ {
		if slices.Contains(mostFreq, freqDist[i].key) {
			frequency++
		}
	}

	for i := dictLen - len(leastFreq) + 3; i < dictLen; i++ {
		if slices.Contains(mostFreq, freqDist[i].key) {
			frequency++
		}
	}

	return frequency
}

func getFreqDist(str []rune) []freqDist[rune] {
	freqDist := make([]freqDist[rune], dictLen)
	for i := 0; i < dictLen; i++ {
		freqDist[i].key = dictRev[i]
		freqDist[i].frequency = 0
	}

	for _, c := range str {
		if !isLetter(c) {
			continue
		}
		c = unicode.ToLower(c)

		freqDist[dict[c]].frequency++
	}

	slices.SortFunc(freqDist, compareFreqDistDesc[rune])

	return freqDist
}

func partitionString(str []rune, keyLen int) [][]rune {
	partitionedStr := make([][]rune, keyLen)
	for i, c := range str {
		partitionedStr[i%keyLen] = append(partitionedStr[i%keyLen], c)
	}

	return partitionedStr
}

func getEnglishDict() map[string]bool {
	engDict, err := http.Get("https://stepik.org/media/attachments/lesson/668860/dictionary.txt")
	if err != nil {
		panic(err)
	}
	defer engDict.Body.Close()
	engDictStr := make([]string, 0)
	scanner := bufio.NewScanner(engDict.Body)
	for scanner.Scan() {
		engDictStr = append(engDictStr, scanner.Text())
	}

	engdic := make(map[string]bool, len(engDictStr))
	for _, word := range engDictStr {
		word = strings.ToLower(word)
		engdic[word] = true
	}

	return engdic
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
	for _, c := range input {
		if c == "\n" {
			continue
		}
		str = append(str, []rune(c)...)
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

func isLetter(c rune) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func compareFreqDistDesc[E freqKey](i, j freqDist[E]) int {

	if i.frequency == j.frequency {
		if i.key > j.key {
			return -1
		} else if i.key < j.key {
			return 1
		}
		return 0
	} else if i.frequency > j.frequency {
		return -1
	}

	return 1
}

func decryptCharStd(char rune, key int) rune {
	return dictRev[(dict[char]-key+dictLen)%dictLen]
}

func decryptCharXor(char rune, key int) rune {
	return dictRev[dict[char]^key]
}
