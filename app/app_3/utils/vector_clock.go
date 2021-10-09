package utils

import (
	"alessandro.it/app/lib"
)

const NUMBER_PEER = lib.NUMBER_NODES

type Vector_clock struct {
	clocks [NUMBER_PEER]int
}

func (v *Vector_clock) Init() {
	for i := 0; i < NUMBER_PEER; i++ {
		v.clocks[i] = 0
	}
}

func (v *Vector_clock) Increment(index int) {
	v.clocks[index] = v.clocks[index] + 1
}

func (v *Vector_clock) Update_with_max(v2 Vector_clock) {
	for i := 0; i < NUMBER_PEER; i++ {
		v.clocks[i] = lib.Max(v.clocks[i], v2.clocks[i])
	}
}

/*
This function return:
	- True, if the vector v1 has, in each position, value less or equal than v2.
	- False, otherwise.
*/
func Less_equal(v1 Vector_clock, v2 Vector_clock) bool {
	for i := 0; i < NUMBER_PEER; i++ {
		if v1.clocks[i] > v2.clocks[i] {
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
		if v1.clocks[i] != v2.clocks[i] {
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
		if v1.clocks[i] > v2.clocks[i] {
			return false
		} else if v1.clocks[i] < v2.clocks[i] {
			less_value_exist = true
		}
	}

	if less_value_exist {
		return true
	} else {
		return false
	}
}
