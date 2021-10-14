/*
This file build a peer that run following the rules of algorithm 3.
*/
package main

import (
	"sync"
	"time"

	"alessandro.it/app/utils"
)

// Definition of third type of Peer
type Peer_3 struct {
	peer         Peer
	vector_clock *utils.Vector_clock
	waiting_list *utils.Waiting_list
	mutex_queue  sync.Mutex
	mutex_clock  sync.Mutex
}

// Initialization of peer
func (p3 *Peer_3) init_peer_3() {
	p3.peer.init_peer()
	p3.vector_clock = &utils.Vector_clock{}
	p3.waiting_list = &utils.Waiting_list{}
}

// This RPC method of Node allow to get update from the other node of group multicast
func (p3 *Peer_3) Get_update(update *utils.Update_vector, empty *utils.Empty) error {
	if update.Packet.Source_address != getIpAddress() {
		p3.mutex_clock.Lock()
		p3.vector_clock.Increment(p3.peer.index)
		p3.mutex_clock.Unlock()
	}

	// Build update node to insert the packet into queue
	update_node := &utils.Waiting_node{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	p3.mutex_queue.Lock()
	p3.waiting_list.Append(update_node)
	p3.mutex_queue.Unlock()

	p3.vector_clock.Print()

	return nil
}

/*
This function check if there are packet to deliver, so the following conditions must be checked:
	1. The message inviato from p_j to current process is the expected message from p_j
	2. The current process has seen every messahe that p_j has seen
*/
func (p3 *Peer_3) deliver_packet() {
	current_index := 1
	for {
		p3.mutex_queue.Lock()
		node_to_deliver := p3.waiting_list.Get_node(current_index)
		p3.mutex_queue.Unlock()
		deliver := true
		index_pid_to_deliver := 0
		if node_to_deliver == nil {
			deliver = false
		} else {
			index_pid_to_deliver = node_to_deliver.Update.Packet.Index_pid
			current_index = current_index + 1
		}

		if deliver && node_to_deliver.Update.Timestamp.Clocks[index_pid_to_deliver] == p3.vector_clock.Clocks[index_pid_to_deliver]+1 {
			for k := 0; k < conf.Nodes && deliver; k++ {
				if k != index_pid_to_deliver && node_to_deliver.Update.Timestamp.Clocks[k] > p3.vector_clock.Clocks[k] {
					deliver = false
				}
			}
		}

		if deliver {
			// Update the vector clock
			p3.vector_clock.Update_with_max(node_to_deliver.Update.Timestamp)

			// Deliver the packet to application layer
			log_message(&node_to_deliver.Update.Packet)

			// Remove the node that is just delivered
			p3.waiting_list.Remove_node(node_to_deliver)
		}
	}
}

// This function send a single message to a single node
func (p3 *Peer_3) send_single_message(index_pid int, update *utils.Update_vector, empty_reply *utils.Empty) {
	// first := true
	/*
		The following 3 lines allow to test the algorithm 3 in case of scenario that we saw in class.
	*/
	// if first && i == 2 {
	// 	time.Sleep(time.Duration(10) * time.Second)
	// 	first = false
	// }
	utils.Delay(MAX_DELAY)

	err := conn.Peer[index_pid].Call("Peer.Get_update", update, empty_reply)
	utils.Check_error(err)
}

// This function get the message from frontend and send it in multicast
func (p3 *Peer_3) Get_message_from_frontend(text *string, empty_reply *utils.Empty) error {
	// Build packet
	pkt := utils.Packet{Source_address: getIpAddress(), Message: *text, Index_pid: p3.peer.index, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// Update the scalar clock and build update packet to send
	p3.mutex_clock.Lock()
	p3.vector_clock.Increment(p3.peer.index)
	update := utils.Update_vector{Timestamp: *p3.vector_clock, Packet: pkt}
	p3.mutex_clock.Unlock()

	// Send to each node of group multicast the message
	for i := 0; i < conf.Nodes; i++ {
		go p3.send_single_message(i, &update, empty_reply)
	}

	return nil
}
