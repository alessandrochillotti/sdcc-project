/*
This file build a peer that run following the rules of algorithm 2.
*/
package main

import (
	"os"
	"strconv"
	"sync"
	"time"

	"alessandro.it/app/utils"
)

// Definition of second type of Peer
type Peer_2 struct {
	Peer          Peer
	scalar_clock  int
	ordered_queue *utils.Queue
	mutex_queue   sync.Mutex
	mutex_clock   sync.Mutex
	wg            sync.WaitGroup
}

// Initialization of peer
func (p2 *Peer_2) init_peer_2(username string) {
	p2.scalar_clock = 0
	p2.ordered_queue = &utils.Queue{}
}

// This RPC method of Node allow to get update from the other node of group multicast
func (p2 *Peer_2) Get_update(update *utils.Update, empty *utils.Empty) error {
	// Reset timer
	conf.Timer.Reset(time.Duration(utils.TIMER_NODE*conf.Nodes) * time.Second)

	// Manage scalar clock
	p2.mutex_clock.Lock()
	p2.scalar_clock = utils.Max(p2.scalar_clock, update.Timestamp)
	p2.scalar_clock = p2.scalar_clock + 1
	p2.mutex_clock.Unlock()

	// Build update node to insert the packet into queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	p2.mutex_queue.Lock()
	p2.ordered_queue.Update_into_queue(update_node, conn.Index)
	p2.ordered_queue.Display(conn.Index)
	p2.mutex_queue.Unlock()

	// Send ack message in multicast
	ack_to_send := &utils.Ack{Ip_addr: update.Packet.Source_address, Timestamp: update.Timestamp}

	for i := 0; i < conf.Nodes; i++ {
		var empty utils.Empty
		conn.Peer[i].Go("Peer.Get_ack", &ack_to_send, &empty, nil)
	}

	return nil
}

// This RPC method of Node allow to receive ack from other nodes of group multicast.
func (p2 *Peer_2) Get_ack(ack *utils.Ack, empty *utils.Empty) error {
	acked := false
	for !acked {
		p2.mutex_queue.Lock()
		acked = p2.ordered_queue.Ack_node(*ack)
		p2.mutex_queue.Unlock()
	}

	return nil
}

// This function log message into file: this has the value of delivery to application layer.
func (p2 *Peer_2) log_message(update_to_deliver *utils.Update) {
	// Open file into volume docker
	path_file := "/docker/node_volume/" + p2.Peer.Ip_address + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Check_error(err)

	_, err = f.WriteString(strconv.Itoa(update_to_deliver.Timestamp) + ";" + update_to_deliver.Packet.Timestamp.Format(time.RFC1123)[17:25] + ";" + conn.GetUsername(update_to_deliver.Packet.Source_address) + ";" + update_to_deliver.Packet.Message + "\n")
	utils.Check_error(err)

	f.Close()
}

/*
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
				if i != p2.Peer.Index {
					p2.mutex_queue.Lock()
					update_max_timestamp := p2.ordered_queue.Get_update_max_timestamp(conn.Addresses[i])
					p2.mutex_queue.Unlock()
					deliver = deliver && (update_max_timestamp.Timestamp > head_node.Timestamp || (update_max_timestamp.Timestamp == head_node.Timestamp && conn.GetIndex(head_node.Packet.Source_address) < conn.GetIndex(update_max_timestamp.Packet.Source_address)))
				}
			}

			if deliver {
				// Deliver the packet to application layer
				p2.log_message(&head_node)

				// Remove the node that is just delivered
				p2.mutex_queue.Lock()
				p2.ordered_queue.Remove_head()
				p2.mutex_queue.Unlock()
			}
		}
	}
}

// This function send a single message to a single node
func (p2 *Peer_2) send_single_message(index_pid int, update *utils.Update, empty_reply *utils.Empty) {
	if index_pid != conn.GetIndex(update.Packet.Source_address) {
		utils.Delay(MAX_DELAY)
	}
	err := conn.Peer[index_pid].Call("Peer.Get_update", update, empty_reply)
	utils.Check_error(err)

	p2.wg.Done()
}

// Frontend communication
func (p2 *Peer_2) Get_message_from_frontend(msg *utils.Message, empty_reply *utils.Empty) error {
	// Reset timer
	conf.Timer.Reset(time.Duration(utils.TIMER_NODE*conf.Nodes) * time.Second)

	// Build packet
	pkt := utils.Packet{Source_address: p2.Peer.Ip_address, Message: msg.Text, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// Update the scalar clock and build update packet to send
	p2.mutex_clock.Lock()
	p2.scalar_clock = p2.scalar_clock + 1
	update := utils.Update{Timestamp: p2.scalar_clock, Packet: pkt}
	p2.mutex_clock.Unlock()

	// Send to each node of group multicast the message
	p2.wg.Add(conf.Nodes)
	for i := 0; i < conf.Nodes; i++ {
		go p2.send_single_message(i, &update, empty_reply)
	}
	p2.wg.Wait()

	return nil
}
