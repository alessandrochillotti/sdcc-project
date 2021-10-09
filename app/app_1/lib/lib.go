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

// This struct is used to send 'first version' of packet to sequencer
type Packet struct {
	Source_address string
	Message        string
}

// This struct is used by sequencer to reply to each node of group multicast
type Packet_sequencer struct {
	Id  int
	Pkt Packet
}

// This struct is use for RPC method when the reply is not important
type Empty struct{}

type Addresses struct {
	Addresses [NUMBER_NODES]string
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
