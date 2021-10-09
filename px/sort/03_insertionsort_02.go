package main

import "fmt"

func insertionSort2() {

	for i := 1; i < len(arr); i++ {
		insert(i, search(i, arr), arr)
	}

	fmt.Println(arr)
}

func insert(souseIndex int, dest int, arr []int) {

	v := arr[souseIndex]

	for i := souseIndex; i > dest; i-- {
		arr[i] = arr[i-1]
	}

	arr[dest] = v
}

func search(index int, arr []int) int {

	begin := 0
	end := index

	for begin < end {

		mid := (begin + end) >> 1

		if arr[mid] > arr[index] {
			end = mid
		} else {
			begin = mid + 1
		}
	}

	return begin
}
