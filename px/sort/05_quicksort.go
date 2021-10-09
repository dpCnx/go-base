package main

import (
	"fmt"
)

func quicksort() {

	qSort(0, len(arr))

	fmt.Println(arr)
}

func qSort(begin, end int) {

	if end-begin < 2 {
		return
	}

	mid := pivotIndex(begin, end)
	qSort(begin, mid)
	qSort(mid+1, end)

}

func pivotIndex(begin, end int) int {

	poivt := arr[begin]
	end--

	for begin < end {

		for begin < end {
			if arr[end] > poivt {
				end--
			} else {
				arr[begin] = arr[end]
				begin++
				break
			}
		}

		for begin < end {
			if arr[begin] < poivt {
				begin++
			} else {
				arr[end] = arr[begin]
				end--
				break
			}
		}
	}
	arr[begin] = poivt
	return begin
}
