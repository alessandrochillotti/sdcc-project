/*
	This is the sequencer node that allow fully ordered multicast implemented centrally.
	It has ip address equal to 10.5.0.253 and it is listening in port 1234.
*/

package main

import (
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"

	"alessandro.it/app/utils"
)

type Sequencer struct {
	current_id int
	Mutex_id   sync.Mutex
	peer       []*rpc.Client
}

// This function send a specific message to each node of group multicast
func (seq *Sequencer) send_single_message(peer_id int, arg *utils.Packet_sequencer, empty *utils.Empty) error {
	// Call remote procedure and reply will store the RPC result
	err := seq.peer[peer_id].Call("Peer.Get_Message", &arg, &empty)
	if err != nil {
		return err
	}

	return nil
}

// This function is called by each generic node to send packet to each node of group multicast
func (seq *Sequencer) Send_packet(arg *utils.Packet, empty *utils.Empty) error {
	// Read file line by line, so scan every ip address
	seq.Mutex_id.Lock()
	seq.current_id = seq.current_id + 1
	// Prepare packet to send
	pkt_seq := utils.Packet_sequencer{Id: seq.current_id, Pkt: *arg}
	seq.Mutex_id.Unlock()

	// Send to each node of group multicast the message
	for i := 0; i < len(seq.peer); i++ {
		go seq.send_single_message(i, &pkt_seq, empty)
	}

	return nil
}

// This RPC method of Node allow to get list from the registered node
func (seq *Sequencer) Get_list(list *utils.List_of_nodes, reply *utils.Empty) error {
	var err error
	nodes, _ := strconv.Atoi(os.Getenv("NODES"))

	seq.peer = make([]*rpc.Client, nodes)

	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < nodes; i++ {
		addr_node := addr_tmp[i] + ":1234"
		seq.peer[i], err = rpc.Dial("tcp", addr_node)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	seq := &Sequencer{current_id: 0}

	// Register a sequencer methods
	sequencer := rpc.NewServer()
	err := sequencer.RegisterName("Sequencer", seq)
	utils.Check_error(err)

	// Listen for incoming messages on port 1234
	lis, err := net.Listen("tcp", ":1234")
	utils.Check_error(err)

	// Use goroutine to implement a lightweight thread to manage new connection
	go sequencer.Accept(lis)

	for {
	}
}
