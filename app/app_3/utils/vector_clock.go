package utils

import (
	"fmt"

	"alessandro.it/app/lib"
)

const NUMBER_PEER = lib.NUMBER_NODES

type Vector_clock struct {
	Clocks [NUMBER_PEER]int
}

func (v *Vector_clock) Init() {
	for i := 0; i < NUMBER_PEER; i++ {
		v.Clocks[i] = 0
	}
}

func (v *Vector_clock) Increment(index int) {
	v.Clocks[index] = v.Clocks[index] + 1
}

func (v *Vector_clock) Update_with_max(v2 Vector_clock) {
	for i := 0; i < NUMBER_PEER; i++ {
		v.Clocks[i] = lib.Max(v.Clocks[i], v2.Clocks[i])
	}
}

func (v *Vector_clock) Print() {
	fmt.Printf("[ ")
	var i int
	for i = 0; i < lib.NUMBER_NODES-1; i++ {
		fmt.Printf("%d, ", v.Clocks[i])
	}
	fmt.Printf("%d ]\n", v.Clocks[i])
}

/*
This function return:
	- True, if the vector v1 has, in each position, value less or equal than v2.
	- False, otherwise.
*/
func Less_equal(v1 Vector_clock, v2 Vector_clock) bool {
	for i := 0; i < NUMBER_PEER; i++ {
		if v1.Clocks[i] > v2.Clocks[i] {
			return false
		}
	}

	return true
}

/*
This function return:
	- True, if the vectors have the same value for each position.
	- False, otherwise.
*/
func Equal(v1 Vector_clock, v2 Vector_clock) bool {
	for i := 0; i < NUMBER_PEER; i++ {
		if v1.Clocks[i] != v2.Clocks[i] {
			return false
		}
	}

	return true
}

/*
This function return:
	- True, if the vector v1 is less than vector v2
	- False, otherwise.
*/
func Less(v1 Vector_clock, v2 Vector_clock) bool {
	less_value_exist := false

	for i := 0; i < NUMBER_PEER; i++ {
		if v1.Clocks[i] > v2.Clocks[i] {
			return false
		} else if v1.Clocks[i] < v2.Clocks[i] {
			less_value_exist = true
		}
	}

	if less_value_exist {
		return true
	} else {
		return false
	}
}
