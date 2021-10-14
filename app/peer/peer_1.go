/*
This file build a peer that run following the rules of algorithm 1.
*/
package main

import (
	"net/rpc"
	"time"

	"alessandro.it/app/utils"
)

// Definition of first type of Peer
type Peer_1 struct {
	peer              Peer
	current_packet_id int
	buffer            chan (utils.Packet_sequencer)
}

// Initialization of peer
func (p1 *Peer_1) init_peer_1() {
	p1.peer.init_peer()
	p1.buffer = make(chan utils.Packet_sequencer, MAX_QUEUE)
}

/*
Algorithm: 1

This function is called by sequencer node for sending message: the message is received, not delivered.
*/
func (p1 *Peer_1) Get_Message(pkt *utils.Packet_sequencer, empty *utils.Empty) error {
	// The packet is received, so it is buffered
	p1.buffer <- *pkt

	return nil
}

// This function check if there are packets to deliver, according to current_id + 1 == current_packet.Id.
func (p1 *Peer_1) deliver_packet() {
	for {
		current_packet := <-p1.buffer

		if p1.current_packet_id+1 == current_packet.Id {
			// Update expected id of packet
			p1.current_packet_id = p1.current_packet_id + 1

			// Deliver the packet to application layer
			log_message(&current_packet.Pkt)
		} else {
			p1.buffer <- current_packet
		}
	}
}

// Frontend communication
func (p1 *Peer_1) Get_message_from_frontend(text *string, empty_reply *utils.Empty) error {
	var empty utils.Empty

	// Build packet
	pkt := utils.Packet{Source_address: getIpAddress(), Message: *text, Index_pid: p1.peer.index, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 1234
	addr_sequencer_node := "10.5.0.253:1234"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_sequencer_node)
	utils.Check_error(err)

	utils.Delay(MAX_DELAY)

	// Send packet to sequencer node
	err = client.Call("Sequencer.Send_packet", &pkt, &empty)
	utils.Check_error(err)

	// Close connection with him
	client.Close()

	return nil
}
