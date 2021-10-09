package main

import (
	"fmt"
)

var leftArray []int

func mergeSort() {

	leftArray = make([]int, len(arr)>>1)
	mSort(0, len(arr))

	fmt.Println(arr)

}

func mSort(begin, end int) {

	if end-begin < 2 {
		return
	}

	mid := (begin + end) >> 1

	mSort(begin, mid)
	mSort(mid, end)
	merge(begin, mid, end)

}

func merge(begin, mid, end int) {

	var li, le = 0, mid - begin
	var ri, re = mid, end
	var ai = begin

	for i := li; i < le; i++ {
		leftArray[i] = arr[begin+i]
	}

	for li < le {
		if ri < re && ((arr[ri] - leftArray[li]) < 0) {
			arr[ai] = arr[ri]
			ai++
			ri++
		} else {
			arr[ai] = leftArray[li]
			ai++
			li++
		}

	}

}
