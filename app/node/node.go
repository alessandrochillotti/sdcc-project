/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"alessandro.it/app/lib"
	"alessandro.it/app/utils"
)

type Node int

/* Constant values */
const MAX_QUEUE = 100
const MAX_DELAY = 3

/* Global variables */

// Process information
var conf *lib.Configuration
var my_index int
var verbose_flag bool

// Algorithm 1 global variables
var current_id = 0
var buffer chan (lib.Packet_sequencer)

// Algorithm 2 global variables
var scalar_clock int = 0
var addresses []string /* Contains ip addresses of each node in multicast group */
var queue *utils.Queue

// Algorithm 3 global variables
var vector_clock *utils.Vector_clock
var queue_2 *utils.Queue_2

// Connection variables
var peer []*rpc.Client

// Mutex variables
var mutex_queue sync.Mutex
var mutex_clock sync.Mutex

// Channel
var channel_connection chan bool

/*
This function return the ip address of current node.
*/
func getIpAddress() string {
	ip_address := "0.0.0.0"
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip_address = ipv4.String()
		}
	}

	return ip_address
}

/*
This function build a struct that contains the info to register the node.
*/
func build_whoami_struct(whoami_to_register *lib.Whoami) {
	whoami_to_register.Ip_address = getIpAddress()
	whoami_to_register.Port = "1234"
}

/*
This function allows to register the node to communicate in multicast group.
*/
func register_into_group() {
	var whoami_to_register lib.Whoami
	var empty lib.Empty

	build_whoami_struct(&whoami_to_register)

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 1234
	addr_register_node := "10.5.0.254:1234"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	lib.Check_error(err)

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.Register_node", &whoami_to_register, &empty)
	lib.Check_error(err)

	client.Close()
}

/*
This function log message into file: this has the value of delivery to application layer.
*/
func log_message(pkt *lib.Packet) {
	// Open file into volume docker
	path_file := "/docker/node_volume/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	lib.Check_error(err)

	if verbose_flag {
		_, err = f.WriteString("[" + pkt.Timestamp.Format(time.RFC1123)[17:25] + "] " + pkt.Source_address + " -> " + pkt.Message + "\n")
		lib.Check_error(err)
	} else {
		_, err = f.WriteString(pkt.Source_address + " -> " + pkt.Message + "\n")
		lib.Check_error(err)
	}

	f.Close()
}

/*
Algorithm: 1, 2, 3

This function has the goal to clear the shell and print all messages received and sended by the current peer.
*/
func print_chat() {
	// Clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Print all messages, received and sended
	content, err := ioutil.ReadFile("/docker/node_volume/" + getIpAddress() + "_log.txt")
	lib.Check_error(err)
	list := string(content)
	print(list)
}

/*
Algorithm: 1

This function check if there are packets to deliver, according to current_id + 1 == current_packet.Id.
*/
func deliver_packet_1() {
	for {
		current_packet := <-buffer
		if current_id+1 == current_packet.Id {
			// Update expected id of packet
			current_id = current_id + 1

			// Deliver the packet to application layer
			log_message(&current_packet.Pkt)

			// print_chat()

		} else {
			buffer <- current_packet
		}
	}
}

/*
Algorithm: 2

This function check if there are packet to deliver, so the following conditions must be checked:
	1. The firse message in the local queue must have acked by every node
	2. The other node of group must not have a packet with timestamp less than the considering packet to deliver
*/
func deliver_packet_2() {
	for {
		mutex_queue.Lock()
		head := queue.Get_head()
		mutex_queue.Unlock()
		if head != nil && head.Ack == conf.Nodes {
			deliver := true
			head_node := head.Update

			for i := 0; i < conf.Nodes; i++ {
				if i != my_index {
					mutex_queue.Lock()
					update_max_timestamp := queue.Get_update_max_timestamp(addresses[i])
					mutex_queue.Unlock()
					deliver = deliver && (update_max_timestamp.Timestamp > head_node.Timestamp || (update_max_timestamp.Timestamp == head_node.Timestamp && head_node.Packet.Index_pid < update_max_timestamp.Packet.Index_pid))
				}
			}

			if deliver {
				// Deliver the packet to application layer
				log_message(&head_node.Packet)

				// Remove the node that is just delivered
				mutex_queue.Lock()
				queue.Remove_head()
				mutex_queue.Unlock()
			}
		}
	}
}

/*
Algorithm: 3

This function check if there are packet to deliver, so the following conditions must be checked:
	1. The message inviato from p_j to current process is the expected message from p_j.
	2. The current process has seen every messahe that p_j has seen.
*/
func deliver_packet_3() {
	current_index := 1
	for {
		mutex_queue.Lock()
		node_to_deliver := queue_2.Get_node(current_index)
		mutex_queue.Unlock()
		deliver := true
		index_pid_to_deliver := 0
		if node_to_deliver == nil {
			deliver = false
		} else {
			index_pid_to_deliver = node_to_deliver.Update.Packet.Index_pid
			current_index = current_index + 1
		}

		if deliver && node_to_deliver.Update.Timestamp.Clocks[index_pid_to_deliver] == vector_clock.Clocks[index_pid_to_deliver]+1 {
			for k := 0; k < conf.Nodes && deliver; k++ {
				if k != index_pid_to_deliver && node_to_deliver.Update.Timestamp.Clocks[k] > vector_clock.Clocks[k] {
					deliver = false
				}
			}
		}

		if deliver {
			// Update the vector clock
			vector_clock.Update_with_max(node_to_deliver.Update.Timestamp)

			// Deliver the packet to application layer
			log_message(&node_to_deliver.Update.Packet)

			// Remove the node that is just delivered
			queue_2.Remove_node(node_to_deliver)
		}
	}
}

/*
Algorithm: 2, 3

This function, after reception of list from register node, allow to setup connection with each node of group multicast.
*/
func setup_connection() error {
	var err error

	for i := 0; i < conf.Nodes; i++ {
		addr_node := addresses[i] + ":1234"
		peer[i], err = rpc.Dial("tcp", addr_node)
		lib.Check_error(err)
		if addresses[i] == getIpAddress() {
			my_index = i
		}
	}

	return nil
}

/* RPC methods registered by Node */

/*
Algorithm: 2, 3

This RPC method of Node allow to get list from the registered node.
*/
func (node *Node) Get_list(list *lib.List_of_nodes, reply *lib.Empty) error {
	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < conf.Nodes; i++ {
		addresses[i] = addr_tmp[i]
	}

	channel_connection <- true

	return nil
}

/*
Algorithm: 2

This RPC method of Node allow to get update from the other node of group multicast
*/
func (node *Node) Get_update(update *utils.Update, ack *utils.Ack) error {
	mutex_clock.Lock()
	scalar_clock = lib.Max(scalar_clock, update.Timestamp)
	scalar_clock = scalar_clock + 1
	mutex_clock.Unlock()

	// Build update node to insert the packet into queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	mutex_queue.Lock()
	queue.Update_into_queue(update_node)
	queue.Display()
	mutex_queue.Unlock()

	// Send ack message in multicast
	ack_to_send := &utils.Ack{Ip_addr: update.Packet.Source_address, Timestamp: update.Timestamp}

	for i := 0; i < conf.Nodes; i++ {
		var empty lib.Empty
		peer[i].Go("Node.Get_ack", &ack_to_send, &empty, nil)
	}

	return nil
}

/*
Algorithm: 3

This RPC method of Node allow to get update from the other node of group multicast
*/
func (node *Node) Get_update_2(update *utils.Update_2, ack *utils.Ack) error {
	if update.Packet.Source_address != getIpAddress() {
		mutex_clock.Lock()
		vector_clock.Increment(my_index)
		mutex_clock.Unlock()
	}

	// Build update node to insert the packet into queue
	update_node := &utils.Node_2{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	mutex_queue.Lock()
	queue_2.Append(update_node)
	mutex_queue.Unlock()

	vector_clock.Print()

	return nil
}

/*
Algorithm: 2

This RPC method of Node allow to receive ack from other nodes of group multicast.
*/
func (node *Node) Get_ack(ack *utils.Ack, empty *lib.Empty) error {
	acked := false
	for acked == false {
		mutex_queue.Lock()
		acked = queue.Ack_node(*ack)
		mutex_queue.Unlock()
	}

	return nil
}

/*
Algorithm: 1

This function is called by sequencer node for sending message: the message is received, not delivered.
*/
func (node *Node) Get_Message(pkt *lib.Packet_sequencer, empty *lib.Empty) error {
	// The packet is received, so it is buffered
	buffer <- *pkt

	return nil
}

func (node *Node) Get_message_from_frontend(text *string, empty_reply *lib.Empty) error {
	var empty lib.Empty
	var ack utils.Ack

	// Build packet
	pkt := lib.Packet{Source_address: getIpAddress(), Message: *text, Index_pid: my_index, Timestamp: time.Now().Add(time.Duration(2) * time.Hour)}

	switch conf.Algorithm {
	case 1:
		// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 1234
		addr_sequencer_node := "10.5.0.253:1234"

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_sequencer_node)
		lib.Check_error(err)

		// Send packet to sequencer node
		err = client.Call("Sequencer.Send_packet", &pkt, &empty)
		lib.Check_error(err)

		// Close connection with him
		client.Close()

		break
	case 2:
		// Update the scalar clock and build update packet to send
		mutex_clock.Lock()
		scalar_clock = scalar_clock + 1
		update := utils.Update{Timestamp: scalar_clock, Packet: pkt}
		mutex_clock.Unlock()

		// Send to each node of group multicast the message
		for i := 0; i < conf.Nodes; i++ {
			lib.Delay(3)
			err := peer[i].Call("Node.Get_update", &update, &ack)
			lib.Check_error(err)
		}

		break
	case 3:
		// Update the scalar clock and build update packet to send
		mutex_clock.Lock()
		vector_clock.Increment(my_index)
		update := utils.Update_2{Timestamp: *vector_clock, Packet: pkt}
		mutex_clock.Unlock()

		first := true

		// Send to each node of group multicast the message
		for i := 0; i < conf.Nodes; i++ {
			// lib.Delay(3)
			/*
				The following 3 lines allow to test the algorithm 3 in case of scenario that we saw in class.
			*/
			if first && i == 2 {
				time.Sleep(time.Duration(10) * time.Second)
				first = false
			}
			err := peer[i].Call("Node.Get_update_2", &update, &ack)
			lib.Check_error(err)
		}

		break
	}

	return nil
}

func (node *Node) Handshake(request *lib.Hand_request, reply *lib.Hand_reply) error {
	verbose_flag = request.Verbose

	reply.Ip_address = getIpAddress()

	return nil
}

func init_configuration() {
	algo, _ := strconv.Atoi(os.Getenv("ALGORITHM"))
	nodes, _ := strconv.Atoi(os.Getenv("NODES"))

	conf = &lib.Configuration{Algorithm: algo, Nodes: nodes}

	addresses = make([]string, conf.Nodes)
	peer = make([]*rpc.Client, nodes)
}

func main() {
	// Init phase
	init_configuration()

	// The node communicates with the register node to register his info
	register_into_group()

	// Allocate object to use it into program execution
	node := new(Node)
	channel_connection = make(chan bool)

	// Create file for log of messages
	f, err := os.Create("/docker/node_volume/" + getIpAddress() + "_log.txt")
	lib.Check_error(err)
	defer f.Close()

	// Register the Node methods
	receiver := rpc.NewServer()
	err = receiver.RegisterName("Node", node)
	lib.Check_error(err)

	// Listen for incoming messages on port 1234
	lis, err := net.Listen("tcp", ":1234")
	lib.Check_error(err)
	defer lis.Close()

	// Setup counter
	switch conf.Algorithm {
	case 1:
		buffer = make(chan lib.Packet_sequencer, MAX_QUEUE)
		current_id = 0
		break
	case 2:
		queue = &utils.Queue{}
		scalar_clock = 0
		break
	case 3:
		queue_2 = &utils.Queue_2{}
		vector_clock = new(utils.Vector_clock)
		vector_clock.Init(conf.Nodes)
		break
	}

	// Use goroutine to implement a lightweight thread to manage the coming of new messages
	go receiver.Accept(lis)

	// Setup the connection with the peer of group multicast after the reception of list
	<-channel_connection
	if setup_connection() != nil {
		os.Exit(1)
	}

	switch conf.Algorithm {
	case 1:
		go deliver_packet_1()
		break
	case 2:
		go deliver_packet_2()
		break
	case 3:
		go deliver_packet_3()
		break
	}

	for {

	}
}
