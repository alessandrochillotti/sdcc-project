package utils

import (
	"fmt"

	"alessandro.it/app/lib"
)

type Update struct {
	Timestamp int
	Packet    lib.Packet
}

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

// This function insert update message into queue maintaining it sorted for timestamp
func (l *Queue) Update_into_queue(update *Node) {
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

	l.Display()
}

// Get number ack of head
func (l *Queue) Get_ack_head() Ack {
	if l.head == nil {
		return 0
	} else {
		return l.head.Ack
	}

}

// Retrive head
func (l *Queue) Get_head() *Node {
	head := l.head
	l.head = l.head.Next

	return head
}

// Debug function
func (l Queue) Display() {
	for l.head != nil {
		fmt.Printf("%v -> %s \n", l.head.Update.Timestamp, l.head.Update.Packet.Message)
		l.head = l.head.Next
	}
	fmt.Println()
}
