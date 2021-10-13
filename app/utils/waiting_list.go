/*
This file contains the utils to manage the list of pending messages (i.e. ALGO 3)
*/

package utils

// This struct define tha waiting list
type Waiting_list struct {
	head *Waiting_node
	tail *Waiting_node
}

// This struct is the node of queue
type Waiting_node struct {
	Update Update_vector
	Next   *Waiting_node
	Ack    int
}

// This struct emulate the update message in algorithm 3 that the peer send in multicast
type Update_vector struct {
	Timestamp Vector_clock
	Packet    Packet
}

// This function insert update message into queue maintaining it sorted for timestamp
func (l *Waiting_list) Append(update *Waiting_node) {
	if l.head == nil {
		l.head = update
		l.tail = update
	} else {
		l.tail.Next = update
		l.tail = update
	}
}

// Remove head
func (l *Waiting_list) Remove_node(node_target *Waiting_node) bool {
	current_node := l.head
	previous_node := current_node

	for current_node != nil {
		if current_node == node_target {
			if current_node == l.head {
				l.head = l.head.Next
			} else {
				previous_node.Next = current_node.Next
			}

			return true
		}
		previous_node = current_node
		current_node = current_node.Next
	}

	return false
}

// This function return the head of queue
func (l *Waiting_list) Get_node(counter int) *Waiting_node {
	current_node := l.head
	if current_node == nil {
		return current_node
	}

	for i := 0; i < counter; i++ {
		if current_node.Next == nil {
			current_node = l.head
		} else {
			current_node = current_node.Next
		}
	}

	return current_node
}
