/*
This file contains general utils useful to implement the algorithms.
*/
package utils

import (
	"log"
	"math/rand"
	"time"
)

// This function wraps the pseudo-random generation of delay to send a packet.
func Delay(max int) {
	// Set the initial seed of PRNG
	rand.Seed(time.Now().UnixNano())
	// Extract a number that is between 0 and 2
	n := rand.Intn(max)
	// Simule the delay computed above
	time.Sleep(time.Duration(n) * time.Second)
}

// This function allow to verify if there is error and return it.
func Check_error(err error) error {
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	return err
}

// This function returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
