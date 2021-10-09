package main

import (
	"fmt"
)

func shellSort() {

	steps := shellStepSequence()
	for _, step := range steps {
		sortShell(step)
	}

	fmt.Println(arr)

}

func sortShell(step int) {

	// i : 第几列
	for i := 0; i < step; i++ {
		// j、j+step、j+2*step、j+3*step
		for j := 0; j < len(arr); j += step {
			cur := j
			for cur > i && (arr[cur]-arr[cur-step] < 0) {
				arr[cur], arr[cur-step] = arr[cur-step], arr[cur]
				cur -= step
			}
		}
	}
}

func shellStepSequence() []int {

	var stepSequence []int
	step := len(arr)

	for step > 0 {
		step = step >> 1
		stepSequence = append(stepSequence, step)
	}

	return stepSequence
}
