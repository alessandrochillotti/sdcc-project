package utils

import (
	"fmt"

	"alessandro.it/app/lib"
)

type Update struct {
	Timestamp int
	Packet    lib.Packet
}

type Ack bool

/*
Definition of struct used for linked list to mantain the queue of packet sorted by timestamp
*/
type Node struct {
	Packet Update
	Next   *Node
}

type Queue struct {
	length int // NON SO SE SERVE
	head   *Node
	tail   *Node
}

// This function insert update message into queue maintaining it sorted for timestamp
func (l *Queue) Update_into_queue(update *Node) {
	if l.head == nil {
		l.head = update
		l.tail = update
		l.length++
	} else {
		previous_node := l.head
		current_node := previous_node
		inserted := false
		for current_node != nil && inserted == true {
			if update.Packet.Timestamp < current_node.Packet.Timestamp {
				node := &Node{Packet: update.Packet}
				previous_node.Next = node
				node.Next = current_node
				inserted = true
			}
			previous_node = current_node
			current_node = current_node.Next
		}
		if inserted == false {
			l.tail.Next = current_node
		}
	}
}

// Debug function
func (l Queue) Display() {
	for l.head != nil {
		fmt.Printf("%v -> ", l.head.Packet.Timestamp)
		l.head = l.head.Next
	}
	fmt.Println()
}
