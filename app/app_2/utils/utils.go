package utils

import (
	"fmt"

	"alessandro.it/app/lib"
)

type Update struct {
	Timestamp int
	Packet    lib.Packet
}

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
	max_id int
	head   *Node
	tail   *Node
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
			if update.Update.Timestamp < current_node.Update.Timestamp {
				previous_node.Next = update
				update.Next = current_node
				inserted = true
			} else {
				previous_node = current_node
				current_node = current_node.Next
			}
		}
		if inserted == false {
			l.tail.Next = update
			l.tail = update
		}
	}

	if update.Update.Packet.Id > l.Get_max_id() {
		l.Set_max_id(update.Update.Packet.Id)
	}

	// fmt.Println("Ho inserito il nodo con id", update.Update.Packet.Id)

	// l.Display()
}

// Put ack for a specific timestamp
func (l *Queue) Ack_node(id int) bool {
	acked := false
	current_node := l.head

	// fmt.Println("Sto cercando di dare l'ack al pacchetto con id", id)

	for current_node != nil && current_node.Update.Packet.Id != id {
		current_node = current_node.Next
	}

	if current_node != nil {
		current_node.Ack = current_node.Ack + 1
		acked = true
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
	}
}

// Debug function
func (l *Queue) Display() {
	for l.head != nil {
		// fmt.Printf("%v -> %s \n", l.head.Update.Timestamp, l.head.Update.Packet.Message)
		l.head = l.head.Next
	}
	fmt.Println()
}

// Return the min timestamp that is inserted into queue, so the timestamp of head node
func (l *Queue) Get_min_timestamp() Timestamp {
	if l.head != nil {
		return Timestamp(l.head.Update.Timestamp)
	}

	return 0
}

// Set the max id
func (l *Queue) Set_max_id(id int) {
	l.max_id = id
}

// Return the max id that is inserted into queue
func (l *Queue) Get_max_id() int {
	return l.max_id
}
