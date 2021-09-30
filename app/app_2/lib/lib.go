package lib

import (
	"log"
	"math/rand"
	"time"
)

/* Constant values */
const NUMBER_NODES = 3

type Outcome bool /* If is true, then ok. */
type List string

/* Structs use for develop the different algorithms */

// This struct identify a specific node
type Whoami struct {
	Ip_address string
	Port       string
}

// This is a packet that contain the message to send to each node og group multicast
type Packet struct {
	Source_address string
	Source_pid     int
	Message        string
}

// This struct is used for RPC methods when the reply is not important
type Empty struct{}

type Addresses struct {
	Addresses_array [NUMBER_NODES]string
}

/* Utility */
func Delay() {
	// Set the initial seed of PRNG
	rand.Seed(time.Now().UnixNano())
	// Extract a number that is between 0 and 2
	n := rand.Intn(3)
	// Simule the delay computed above
	time.Sleep(time.Duration(n) * time.Second)
}

// This function allow to verify if there is error and return it
func Check_error(err error) error {
	if err != nil {
		log.Printf("unable to read file: %v", err)
	}
	return err
}
