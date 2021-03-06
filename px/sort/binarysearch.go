package main

/**
 * 查找v在有序数组array中的位置
 */

func indexOf(arr []int, v int) int {
	begin := 0
	end := len(arr)

	for begin < end {

		mid := (begin + end) >> 1

		if arr[mid] < v {
			begin = mid + 1
		} else if arr[mid] > v {
			end = mid
		} else {
			return mid
		}
	}

	return -1
}
