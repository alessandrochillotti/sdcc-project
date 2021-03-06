package utils

import (
	"time"
)

/* Structs use for develop the different algorithms */

// This struct identify a specific node
type Whoami struct {
	Ip_address string
	Username   string
	Port       int
}

// This is a packet that contain the message to send to each node og group multicast
type Packet struct {
	Source_address string
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
	Test     bool
}
type Hand_reply struct {
	Ip_address string
	Algorithm  int
}

// This struct is used as reply by register node to send the list of nodes
type List_of_nodes struct {
	List_str string
}

// This struct is used to send packet from frontend to peer
type Message struct {
	Text  string
	Delay []int
}

// This struct is use for RPC method when the reply is not important
type Empty struct{}

// This struct emulate the update message in algorithm 3 that the peer send in multicast
type Update_vector struct {
	Timestamp []int
	Packet    Packet
}

// This struct is used to configurate the variables useful for envinroment configuration
type Configuration struct {
	Algorithm int
	Nodes     int
	Verbose   bool
	Test      bool
	Timer     *time.Timer
}
