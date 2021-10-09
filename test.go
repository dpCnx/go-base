package main

import (
	"fmt"
)

func main() {
	s := make([]int, 5, 5)
	fmt.Println(s)
	fmt.Println(len(s))
	fmt.Println(cap(s))
	s = append(s, 1)
	fmt.Println(s)
	fmt.Println(len(s))
	fmt.Println(cap(s))

	fmt.Println("--------------")
	var a []int
	fmt.Println(a)
	fmt.Println(len(a))
	fmt.Println(cap(a))
	a = append(a, 1,2,3)
	fmt.Println(a)
	fmt.Println(len(a))
	fmt.Println(cap(a))
}
