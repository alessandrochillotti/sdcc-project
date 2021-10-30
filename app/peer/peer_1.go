/*
This file build a peer that run following the rules of algorithm 1.
*/
package main

import (
	"net/rpc"
	"os"
	"strconv"
	"time"

	"alessandro.it/app/utils"
)

// Definition of first type of Peer
type Peer_1 struct {
	Peer              Peer
	current_packet_id int
	buffer            chan (utils.Packet_sequencer)
}

// Initialization of peer
func (p1 *Peer_1) init_peer_1(username string) {
	p1.buffer = make(chan utils.Packet_sequencer, MAX_QUEUE)
}

// This function is called by sequencer node for sending message: the message is received, not delivered.
func (p1 *Peer_1) Get_Message(pkt *utils.Packet_sequencer, empty *utils.Empty) error {
	// The packet is received, so it is buffered
	p1.buffer <- *pkt

	return nil
}

// This function log message into file: this has the value of delivery to application layer.
func (p1 *Peer_1) log_message(pkt_to_deliver *utils.Packet_sequencer) {
	// Open file into volume docker
	path_file := "/docker/node_volume/" + p1.Peer.Ip_address + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Check_error(err)

	_, err = f.WriteString(strconv.Itoa(pkt_to_deliver.Id) + ";" + pkt_to_deliver.Pkt.Timestamp.Format(time.RFC1123)[17:25] + ";" + conn.GetUsername(pkt_to_deliver.Pkt.Source_address) + ";" + pkt_to_deliver.Pkt.Message + "\n")
	utils.Check_error(err)

	f.Close()
}

// This function check if there are packets to deliver, according to current_id + 1 == current_packet.Id.
func (p1 *Peer_1) deliver_packet() {
	for {
		current_packet := <-p1.buffer

		if p1.current_packet_id+1 == current_packet.Id {
			// Update expected id of packet
			p1.current_packet_id = p1.current_packet_id + 1

			// Deliver the packet to application layer
			p1.log_message(&current_packet)
		} else {
			p1.buffer <- current_packet
		}
	}
}

// Frontend communication
func (p1 *Peer_1) Get_message_from_frontend(msg *utils.Message, empty_reply *utils.Empty) error {
	var empty utils.Empty

	// Build packet
	pkt := utils.Packet{Source_address: p1.Peer.Ip_address, Message: msg.Text, Index_pid: p1.Peer.Index, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 1234
	addr_sequencer_node := "10.5.0.253:1234"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_sequencer_node)
	utils.Check_error(err)

	utils.Delay(MAX_DELAY)

	// Send packet to sequencer node
	err = client.Call("Sequencer.Send_packet", pkt, &empty)
	utils.Check_error(err)

	// Close connection with him
	client.Close()

	return nil
}
