package utils

import (
	"fmt"

	"alessandro.it/app/lib"
)

type Timestamp int
type Ack int

/*
Definition of struct used for linked list to mantain the queue of packet sorted by timestamp
*/
type Node struct {
	Update Update
	Next   *Node
	Ack    Ack
}

type Queue struct {
	head *Node
	tail *Node
}

/* Definition of update message used in algorithm 2 and 3. */
type Update struct {
	Timestamp Vector_clock
	Packet    lib.Packet
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
			if Less(update.Update.Timestamp, current_node.Update.Timestamp) {
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
func (l *Queue) Ack_node(id int) bool {
	acked := false
	current_node := l.head

	for current_node != nil && current_node.Update.Packet.Id != id {
		current_node = current_node.Next
	}

	if current_node != nil {
		current_node.Ack = current_node.Ack + 1
		acked = true
		// fmt.Println("Ora il numero di ack per il l'id", id, "è", current_node.Ack)
	}

	return acked
}

// Get number ack of head
func (l *Queue) Get_ack_head() Ack {
	if l.head != nil {
		return l.head.Ack
	}

	return 0
}

// Get head
func (l *Queue) Get_head() *Node {
	return l.head
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

// Debug function
func (l *Queue) Display() {
	current_node := l.head
	fmt.Println("[timestamp, pid, id]")
	for current_node != nil {
		// fmt.Printf("%v -> ", current_node.Update.Packet.Id)
		fmt.Printf("[%v, %d, %d] -> ", current_node.Update.Timestamp, current_node.Update.Packet.Index_pid, current_node.Update.Packet.Id)
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
