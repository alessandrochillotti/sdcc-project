/*
This file build a peer that run following the rules of algorithm 3.
*/
package main

import (
	"os"
	"strconv"
	"sync"
	"time"

	"alessandro.it/app/utils"
)

// Definition of third type of Peer
type Peer_3 struct {
	Peer         Peer
	vector_clock *utils.Vector_clock
	waiting_list chan (utils.Update_vector)
	mutex_clock  sync.Mutex
	wg           sync.WaitGroup
}

// Initialization of peer
func (p3 *Peer_3) init_peer_3() {
	p3.vector_clock = &utils.Vector_clock{}
	p3.waiting_list = make(chan utils.Update_vector, MAX_QUEUE)
}

// This function log message into file: this has the value of delivery to application layer.
func (p3 *Peer_3) log_message(update_to_deliver *utils.Update_vector) {
	// Open file into volume docker
	path_file := "/docker/node_volume/" + p3.Peer.Ip_address + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Check_error(err)

	timestamp_str := ""
	for i := 0; i < conf.Nodes-1; i++ {
		timestamp_str = timestamp_str + strconv.Itoa(update_to_deliver.Timestamp[i]) + ","
	}
	timestamp_str = timestamp_str + strconv.Itoa(update_to_deliver.Timestamp[conf.Nodes-1])

	_, err = f.WriteString(timestamp_str + ";" + update_to_deliver.Packet.Timestamp.Format(time.RFC1123)[17:25] + ";" + conn.GetUsername(update_to_deliver.Packet.Source_address) + ";" + update_to_deliver.Packet.Message + "\n")
	utils.Check_error(err)

	f.Close()
}

// This RPC method of Node allow to get update from the other node of group multicast
func (p3 *Peer_3) Get_update(update utils.Update_vector, empty *utils.Empty) error {
	// Reset timer
	conf.Timer.Reset(time.Duration(utils.TIMER_NODE*conf.Nodes) * time.Second)

	// Insert update node into queue
	p3.waiting_list <- update

	return nil
}

/*
This function check if there are packet to deliver, so the following conditions must be checked:
	1. The message sended from p_j to current process is the expected message from p_j.
	2. The current process has seen every messahe that p_j has seen.
*/
func (p3 *Peer_3) deliver_packet() {
	for {
		deliver := false
		node_to_deliver := <-p3.waiting_list
		index_pid_to_deliver := conn.GetIndex(node_to_deliver.Packet.Source_address)

		t_i := node_to_deliver.Timestamp[index_pid_to_deliver]
		v_j_i := p3.vector_clock.Clocks[index_pid_to_deliver]

		if p3.Peer.Index == index_pid_to_deliver {
			deliver = true
		} else if t_i == v_j_i+1 {
			deliver = true
			for k := 0; k < conf.Nodes && deliver; k++ {
				if k != index_pid_to_deliver {
					t_k := node_to_deliver.Timestamp[k]
					v_j_k := p3.vector_clock.Clocks[k]
					if t_k > v_j_k {
						deliver = false
					}
				}
			}
		}

		if deliver {
			// Update the vector clock
			if p3.Peer.Index != index_pid_to_deliver {
				p3.mutex_clock.Lock()
				p3.vector_clock.Increment(index_pid_to_deliver)
				p3.mutex_clock.Unlock()
			}

			// Deliver the packet to application layer
			p3.log_message(&node_to_deliver)
		} else {
			p3.waiting_list <- node_to_deliver
		}
	}
}

// This function send a single message to a single node
func (p3 *Peer_3) send_single_message(index_pid int, delay int, update utils.Update_vector, empty_reply *utils.Empty) {
	if conf.Test {
		time.Sleep(time.Duration(delay) * time.Second)
	} else if index_pid != conn.GetIndex(update.Packet.Source_address) {
		utils.Delay(MAX_DELAY)
	}

	err := conn.Peer[index_pid].Call("Peer.Get_update", update, empty_reply)
	utils.Check_error(err)

	p3.wg.Done()
}

// This function get the message from frontend and send it in multicast
func (p3 *Peer_3) Get_message_from_frontend(msg *utils.Message, empty_reply *utils.Empty) error {
	// Reset timer
	conf.Timer.Reset(time.Duration(utils.TIMER_NODE*conf.Nodes) * time.Second)

	// Build packet
	pkt := utils.Packet{Source_address: p3.Peer.Ip_address, Message: msg.Text, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	// Update the scalar clock
	p3.mutex_clock.Lock()
	p3.vector_clock.Increment(p3.Peer.Index)

	// Make timestamp
	timestamp := make([]int, conf.Nodes)
	for i := 0; i < conf.Nodes; i++ {
		timestamp[i] = p3.vector_clock.Clocks[i]
	}
	// Build update packet to send
	update := utils.Update_vector{Timestamp: timestamp, Packet: pkt}
	p3.mutex_clock.Unlock()

	// Send to each node of group multicast the message
	p3.wg.Add(conf.Nodes)
	for i := 0; i < conf.Nodes; i++ {
		if msg.Delay == nil {
			go p3.send_single_message(i, 0, update, empty_reply)
		} else {
			go p3.send_single_message(i, msg.Delay[i], update, empty_reply)
		}
	}
	p3.wg.Wait()

	return nil
}
