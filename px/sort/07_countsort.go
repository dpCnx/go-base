package main

import (
	"fmt"
)

func countSort(arr []int) {

	var max = arr[0]
	var min = arr[0]

	for _, v := range arr {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	var counts = make([]int, max-min+1)

	for i := range arr {
		counts[arr[i]-min] += 1
	}

	for i := 1; i < len(counts); i++ {
		counts[i] += counts[i-1]
	}

	var newArr = make([]int, len(arr))

	for i := len(newArr) - 1; i >= 0; i-- {
		newArr[counts[arr[i]-min]-1] = arr[i]
		counts[arr[i]-min] = counts[arr[i]-min] - 1
	}

	fmt.Println(newArr)
}
