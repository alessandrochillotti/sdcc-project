/*
This file define the vector clock and the utils useful to implement the algorithm 3.
*/
package utils

import (
	"fmt"
)

type Vector_clock struct {
	Clocks []int
}

func (v *Vector_clock) Init(nodes int) {
	v.Clocks = make([]int, nodes)
	for i := 0; i < len(v.Clocks); i++ {
		v.Clocks[i] = 0
	}
}

func (v *Vector_clock) Increment(index int) {
	v.Clocks[index] = v.Clocks[index] + 1
}

func (v *Vector_clock) Print() {
	fmt.Printf("[ ")
	var i int
	for i = 0; i < len(v.Clocks)-1; i++ {
		fmt.Printf("%d, ", v.Clocks[i])
	}
	fmt.Printf("%d ]\n", v.Clocks[i])
}
