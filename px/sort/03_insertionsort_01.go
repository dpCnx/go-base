package main

import (
	"fmt"
)

func insertionSort1() {

	for i := 1; i < len(arr); i++ {
		cur := i
		for cur > 0 && (arr[cur]-arr[cur-1]) < 0 {
			arr[cur], arr[cur-1] = arr[cur-1], arr[cur]
			cur--
		}
	}

	fmt.Println(arr)
}
