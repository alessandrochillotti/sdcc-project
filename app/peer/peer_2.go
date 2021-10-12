/*
This file build a peer that run following the rules of algorithm 2.
*/
package main

import (
	"sync"
	"time"

	"alessandro.it/app/utils"
)

type Peer_2 struct {
	peer          Peer
	scalar_clock  int
	ordered_queue *utils.Queue
	mutex_queue   sync.Mutex
	mutex_clock   sync.Mutex
}

func (p2 *Peer_2) init_peer_2() {
	p2.peer.init_peer()
	p2.scalar_clock = 0
	p2.ordered_queue = &utils.Queue{}
}

/*
Algorithm: 2

This RPC method of Node allow to get update from the other node of group multicast
*/
func (p2 *Peer_2) Get_update(update *utils.Update, ack *utils.Ack) error {
	p2.mutex_clock.Lock()
	p2.scalar_clock = utils.Max(p2.scalar_clock, update.Timestamp)
	p2.scalar_clock = p2.scalar_clock + 1
	p2.mutex_clock.Unlock()

	// Build update node to insert the packet into queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	p2.mutex_queue.Lock()
	p2.ordered_queue.Update_into_queue(update_node)
	p2.ordered_queue.Display()
	p2.mutex_queue.Unlock()

	// Send ack message in multicast
	ack_to_send := &utils.Ack{Ip_addr: update.Packet.Source_address, Timestamp: update.Timestamp}

	for i := 0; i < conf.Nodes; i++ {
		var empty utils.Empty
		peer[i].Go("Peer.Get_ack", &ack_to_send, &empty, nil)
	}

	return nil
}

/*
Algorithm: 2

This RPC method of Node allow to receive ack from other nodes of group multicast.
*/
func (p2 *Peer_2) Get_ack(ack *utils.Ack, empty *utils.Empty) error {
	acked := false
	for acked == false {
		p2.mutex_queue.Lock()
		acked = p2.ordered_queue.Ack_node(*ack)
		p2.mutex_queue.Unlock()
	}

	return nil
}

/*
Algorithm: 2

This function check if there are packet to deliver, so the following conditions must be checked:
	1. The firse message in the local queue must have acked by every node
	2. The other node of group must not have a packet with timestamp less than the considering packet to deliver
*/
func (p2 *Peer_2) deliver_packet() {
	for {
		p2.mutex_queue.Lock()
		head := p2.ordered_queue.Get_head()
		p2.mutex_queue.Unlock()
		if head != nil && head.Ack == conf.Nodes {
			deliver := true
			head_node := head.Update

			for i := 0; i < conf.Nodes; i++ {
				if i != p2.peer.index {
					p2.mutex_queue.Lock()
					update_max_timestamp := p2.ordered_queue.Get_update_max_timestamp(addresses[i])
					p2.mutex_queue.Unlock()
					deliver = deliver && (update_max_timestamp.Timestamp > head_node.Timestamp || (update_max_timestamp.Timestamp == head_node.Timestamp && head_node.Packet.Index_pid < update_max_timestamp.Packet.Index_pid))
				}
			}

			if deliver {
				// Deliver the packet to application layer
				log_message(&head_node.Packet)

				// Remove the node that is just delivered
				p2.mutex_queue.Lock()
				p2.ordered_queue.Remove_head()
				p2.mutex_queue.Unlock()
			}
		}
	}
}

// Frontend communication
func (p2 *Peer_2) Get_message_from_frontend(text *string, empty_reply *utils.Empty) error {
	var ack utils.Ack

	// Build packet
	pkt := utils.Packet{Source_address: getIpAddress(), Message: *text, Index_pid: p2.peer.index, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// Update the scalar clock and build update packet to send
	p2.mutex_clock.Lock()
	p2.scalar_clock = p2.scalar_clock + 1
	update := utils.Update{Timestamp: p2.scalar_clock, Packet: pkt}
	p2.mutex_clock.Unlock()

	// Send to each node of group multicast the message
	for i := 0; i < conf.Nodes; i++ {
		utils.Delay(3)
		err := peer[i].Call("Peer.Get_update", &update, &ack)
		utils.Check_error(err)
	}

	return nil
}
