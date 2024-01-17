package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type adjacencyList[T int | tuple] struct {
	list       map[T][]T
	sortedKeys []T
}

type tuple struct {
	x, y int
}

func main() {
	inputList := getInput()

	adjList := adjacencyList[int]{list: inputList}
	adjList.sortedKeys = getSortedKeys(adjList.list)

	set1, set2 := depthFirstSearchBipartition(adjList, 1)

	if checkHall(adjList, set1, set2) {
		fmt.Print("Istnieje skojarzenie doskonałe")
	} else {
		fmt.Print("Nie istnieje skojarzenie doskonałe")
	}

	// fmt.Println(getMaximumMatching(adjList, set1, set2))
}

/* EXAMPLE PERFECT MATCHING
func getMaximumMatching(adjList adjacencyList[int], set1, set2 []int) []tuple {
	matching := make([]tuple, 0)
	smallerSet := make([]int, 0)
	largerSet := make([]int, 0)
	if len(set1) <= len(set2) {
		smallerSet = set1
		largerSet = set2
	} else {
		smallerSet = set2
		largerSet = set1
	}
	matchingAssignment := make([]int, len(adjList.list)+1) // contains redundant empty space for easier indexing

	for _, node := range smallerSet {
		visited := make([]bool, len(adjList.sortedKeys)+1)

		assign(adjList, node, visited, matchingAssignment, largerSet)
	}

	for i := 1; i < len(matchingAssignment); i++ {
		if matchingAssignment[i] != 0 {
			matching = append(matching, tuple{matchingAssignment[i], i})
		}
	}

	return matching
}


func assign(adjList adjacencyList[int], node int, visited []bool, matchingAssignment []int, largerSet []int) bool {
	for _, candidate := range largerSet {
		if slices.Contains(adjList.list[node], candidate) && !visited[candidate] {
			visited[candidate] = true

			// assign if node is free or its current assignment can be reassigned recursively
			if matchingAssignment[candidate] == 0 || assign(adjList, matchingAssignment[candidate], visited, matchingAssignment, largerSet) {
				matchingAssignment[candidate] = node
				return true
			}
		}
	}

	return false
}
*/

func checkHall(adjList adjacencyList[int], set1, set2 []int) bool {
	// only need to check the smaller set
	smallerSet := make([]int, 0)

	if len(set1) <= len(set2) {
		smallerSet = set1
	} else {
		smallerSet = set2
	}

	// check all subsets of the smaller set
	for i := 0; i < 1<<len(smallerSet); i++ {
		subset := make([]int, 0)
		for j := 0; j < len(smallerSet); j++ {
			if i&(1<<j) != 0 {
				subset = append(subset, smallerSet[j])
			}
		}

		neighborhood := getNeighborhood(adjList, subset)
		// if neighborhood is smaller than the subset, there is no perfect matching
		if len(subset) > len(neighborhood) {
			return false
		}
	}

	return true
}

func getNeighborhood(adjList adjacencyList[int], nodes []int) []int {
	neighborhood := make([]int, 0)
	addedNeighbors := make(map[int]bool, len(adjList.sortedKeys)) // additional map for a faster lookup

	for _, node := range nodes {
		for _, neighbor := range adjList.list[node] {
			if !addedNeighbors[neighbor] {
				neighborhood = append(neighborhood, neighbor)
				addedNeighbors[neighbor] = true
			}
		}
	}

	return neighborhood
}

func depthFirstSearchBipartition(graph adjacencyList[int], startingNode int) ([]int, []int) {
	stack := make([]int, 1, len(graph.sortedKeys))
	stack[0] = startingNode

	set1 := make([]int, 0)
	set2 := make([]int, 0)
	coloring := make([]int, len(graph.sortedKeys)+1)
	coloring[startingNode] = 1
	set1 = append(set1, startingNode)

	visited := make(map[int]bool)

	visited[startingNode] = true

	for len(stack) != 0 {
		hasUnvisitedNeighbors := false
		for _, neighbor := range graph.list[peek(stack)] {
			if !visited[neighbor] {
				// coloring part -- color the neighbor with the opposite color of the current node
				if coloring[peek(stack)] == 1 {
					coloring[neighbor] = 2
					set2 = append(set2, neighbor)
				} else {
					coloring[neighbor] = 1
					set1 = append(set1, neighbor)
				}

				// DFS part
				stack = append(stack, neighbor)
				visited[neighbor] = true
				hasUnvisitedNeighbors = true
				break
			}
		}
		if !hasUnvisitedNeighbors {
			pop(&stack)
		}
	}

	return set1, set2
}

func getInput() map[int][]int {
	inputArray := readInput()
	list := make(map[int][]int)

	for i := 0; i < len(inputArray); i++ {
		list[inputArray[i][0]] = inputArray[i][1:]
	}

	return list
}

func readInput() [][]int {
	reader := bufio.NewReader(os.Stdin)

	// read the input as strings separated by new line
	stringArray := []string{}
	for i := 0; true; i++ {
		line, err := reader.ReadString('\n')
		stringArray = append(stringArray, line)

		if err != nil {
			break
		}
	}

	// convert the input string array to a 2d array of substrings
	dividedStrings := make([][]string, len(stringArray), 10000000)
	for i := range stringArray {
		dividedLine := strings.Fields(stringArray[i])

		dividedStrings[i] = make([]string, len(dividedLine))
		dividedStrings[i] = dividedLine
	}

	// convert the string 2d array to a int 2d array
	inputAsInts := make([][]int, len(dividedStrings))
	for i, line := range dividedStrings {
		inputAsInts[i] = make([]int, len(dividedStrings[i]))

		for j, stringNumber := range line {
			number, err := strconv.Atoi(stringNumber)
			if err != nil {
				panic(err)
			}

			inputAsInts[i][j] = number
		}
	}

	return inputAsInts
}

func getSortedKeys(list map[int][]int) []int {
	keys := make([]int, 0)
	for key := range list {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

/* STACK FUNCTIONS */

func pop(slice *[]int) int {
	popped := (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
	return popped
}

func peek(slice []int) int {
	peeked := slice[len(slice)-1]
	return peeked
}
