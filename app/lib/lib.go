package lib

/* Constant values */
const NUMBER_NODES = 3

type Outcome bool /* If is true, then ok. */
type List string

type Whoami struct {
	Ip_address string
	Port       string
}

type Packet struct {
	Source_address string
	Source_pid     int
	Message        string
}

type Packet_sequencer struct {
	Id  int
	Pkt Packet
}
