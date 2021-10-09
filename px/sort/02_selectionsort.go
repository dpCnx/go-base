package main

import "fmt"

func selectSort() {

	for i := len(arr) - 1; i > 0; i-- {
		index := 0
		for j := 1; j <= i; j++ {
			if arr[j] > arr[index] {
				index = j
			}
		}
		arr[i], arr[index] = arr[index], arr[i]
	}

	fmt.Println(arr)
}
