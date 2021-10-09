package main

import (
	"fmt"
)

/*
	并查集
*/

var parents []int

func initUnionFind(length int) {

	parents = make([]int, length)

	for i := range parents {
		parents[i] = i
	}

	fmt.Println(parents)

}
