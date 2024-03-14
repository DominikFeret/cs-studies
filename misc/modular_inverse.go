package main

import "fmt"

func main() {
	var num int
	var mod int

	fmt.Scanf("%d %d", &num, &mod)

	fmt.Println(modularInverse(num, mod))
}

func modularInverse(num int, originalMod int) int {
	inverse := 0
	newInverse := 1
	mod := originalMod
	newMod := num

	for newMod != 0 {
		quotient := mod / newMod

		fmt.Printf("%d = %d * %d + %d", mod, newMod, quotient, mod%newMod)
		fmt.Printf("\t%d = %d - %d * %d\n", inverse, newInverse, quotient, inverse-quotient*newInverse)
		inverse, newInverse = newInverse, inverse-quotient*newInverse
		mod, newMod = newMod, mod-quotient*newMod
	}

	if mod != 1 {
		return -1
	}

	if inverse < 0 {
		inverse += originalMod
	}

	return inverse
}
