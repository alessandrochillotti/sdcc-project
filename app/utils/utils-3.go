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
func (l *Queue) Append(update *Node) {
	if l.head == nil {
		l.head = update
		l.tail = update
	} else {
		l.tail.Next = update
		l.tail = update
	}
}

// Remove head
func (l *Queue) Remove_node(id int) bool {
	current_node := l.head
	previous_node := current_node

	for current_node != nil {
		if current_node.Update.Packet.Id == id {
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

// This function return the head of queue
func (l *Queue) Get_node(id int) *Node {
	current_node := l.head

	for current_node != nil {
		if current_node.Update.Packet.Id == id {
			return current_node
		}
		current_node = current_node.Next
	}

	return nil
}
