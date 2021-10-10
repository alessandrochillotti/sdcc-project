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
	Id             int // This is useful to ack the packet
	Source_address string
	Index_pid      int
	Message        string
}

// This struct is used by sequencer to reply to each node of group multicast
type Packet_sequencer struct {
	Id  int
	Pkt Packet
}

// This struct is use for RPC method when the reply is not important
type Empty struct{}

type List_of_nodes struct {
	List_str string
}

type Deliver struct {
	Ok bool
}

/* Utility */
func Delay(max int) {
	// Set the initial seed of PRNG
	rand.Seed(time.Now().UnixNano())
	// Extract a number that is between 0 and 2
	n := rand.Intn(max)
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

// This function returns the larger of x or y
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
