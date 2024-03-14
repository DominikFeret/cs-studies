package main

import "fmt"

func main() {
	var num int
	var exp int
	var mod int
	fmt.Scanf("%d %d %d", &num, &exp, &mod)

	fmt.Println(modulo(num, exp, mod))
}

func modulo(num int, exp int, mod int) int {
	if exp == 0 {
		return 1
	}
	if exp%2 == 1 {
		fmt.Println("num: ", num, " exp: ", exp, " mod: ", mod, "I'M ODD")
		return (num * modulo(num, exp-1, mod)) % mod
	}
	fmt.Println("num: ", num, " exp: ", exp, " mod: ", mod, "I'M EVEN")
	return modulo((num*num)%mod, exp/2, mod)
}
