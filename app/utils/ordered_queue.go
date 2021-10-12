/*
This file contains the utils to manage the list of update message sorted by timestamp (i.e. ALGO 2)
*/

package utils

import "fmt"

// This struct define the queue sorted by timestamp.
type Queue struct {
	head *Node
	tail *Node
}

// This struct is the node of queue
type Node struct {
	Update Update
	Next   *Node
	Ack    int
}

// This struct emulate the update message in algorithm 2 that the peer send in multicast
type Update struct {
	Timestamp int
	Packet    Packet
}

// This function insert update message into queue maintaining it sorted for timestamp
func (l *Queue) Update_into_queue(update *Node) {
	// fmt.Println("Sto inserendo il nodo con id", update.Update.Packet.Id)
	if l.head == nil {
		l.head = update
		l.tail = update
	} else {
		previous_node := l.head
		current_node := previous_node
		inserted := false
		for current_node != nil && inserted == false {
			if (update.Update.Timestamp < current_node.Update.Timestamp) || ((update.Update.Timestamp == current_node.Update.Timestamp) && (update.Update.Packet.Index_pid < current_node.Update.Packet.Index_pid)) {
				if previous_node != current_node {
					previous_node.Next = update
					update.Next = current_node
				} else {
					l.head = update
					l.head.Next = current_node
				}
				inserted = true
			} else {
				previous_node = current_node
				current_node = current_node.Next
			}
		}
		if inserted == false {
			l.tail.Next = update
			l.tail = update
			l.tail.Next = nil
		}
	}
}

// Put ack for a specific timestamp
func (l *Queue) Ack_node(ack_received Ack) bool {
	acked := false
	current_node := l.head

	for current_node != nil && (current_node.Update.Packet.Source_address != ack_received.Ip_addr || current_node.Update.Timestamp != ack_received.Timestamp) {
		current_node = current_node.Next
	}

	if current_node != nil {
		current_node.Ack = current_node.Ack + 1
		acked = true
	}

	return acked
}

// Get number ack of head
func (l *Queue) Get_ack_head() int {
	if l.head != nil {
		return l.head.Ack
	}

	return 0
}

// Get head
func (l *Queue) Get_head() *Node {
	return l.head
}

// Debug function
func (l *Queue) Display() {
	current_node := l.head
	fmt.Println("[timestamp, pid, id]")
	for current_node != nil {
		fmt.Printf("[%v, %d] -> ", current_node.Update.Timestamp, current_node.Update.Packet.Index_pid)
		current_node = current_node.Next
	}
	fmt.Printf("\n")
}

// Return the min timestamp with a specific ip_addr that is inserted into queue
func (l *Queue) Get_update_max_timestamp(ip_addr string) Update {
	current_node := l.head
	var update_max_timestamp Update
	for current_node != nil {
		if current_node.Update.Packet.Source_address == ip_addr {
			update_max_timestamp = current_node.Update
		}
		current_node = current_node.Next
	}

	return update_max_timestamp
}
