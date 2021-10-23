package utils

import (
	"time"
)

/* Structs use for develop the different algorithms */

// This struct identify a specific node
type Whoami struct {
	Ip_address string
	Port       string
}

// This is a packet that contain the message to send to each node og group multicast
type Packet struct {
	Username       string
	Source_address string
	Index_pid      int
	Message        string
	Timestamp      time.Time
}

// This struct is used by sequencer to reply to each node of group multicast
type Packet_sequencer struct {
	Id  int
	Pkt Packet
}

// This struct emulate ack packet
type Ack struct {
	Ip_addr   string
	Timestamp int
}

// This struct is used by frontend to estabilsh a specific type of communication
type Hand_request struct {
	Username string
}
type Hand_reply struct {
	Ip_address string
	Algorithm  int
}

// This struct is used as reply by register node to send the list of nodes
type List_of_nodes struct {
	List_str string
}

// This struct is use for RPC method when the reply is not important
type Empty struct{}

// This struct is used to configurate the variables useful for envinroment configuration
type Configuration struct {
	Algorithm int
	Nodes     int
	Verbose   bool
}
