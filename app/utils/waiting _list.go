/*
This file contains the utils to manage the list of pending messages (i.e. ALGO 3)
*/

package utils

// This struct define tha waiting list
type Queue_2 struct {
	head *Node_2
	tail *Node_2
}

// This struct is the node of queue
type Node_2 struct {
	Update Update_2
	Next   *Node_2
	Ack    int
}

// This struct emulate the update message in algorithm 3 that the peer send in multicast
type Update_2 struct {
	Timestamp Vector_clock
	Packet    Packet
}

// This function insert update message into queue maintaining it sorted for timestamp
func (l *Queue_2) Append(update *Node_2) {
	if l.head == nil {
		l.head = update
		l.tail = update
	} else {
		l.tail.Next = update
		l.tail = update
	}
}

// Remove head
func (l *Queue) Remove_head() {
	if l.head != nil {
		// fmt.Println("Sto rimuovendo la testa che ha id", l.head.Update.Packet.Id)
		l.head = l.head.Next
		if l.head == nil {
			l.tail = nil
		}
		// fmt.Println("Ora la lista è:")
		l.Display()
	}
}

// Remove head
func (l *Queue_2) Remove_node(node_target *Node_2) bool {
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
func (l *Queue_2) Get_node(counter int) *Node_2 {
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
