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
	"time"

	"alessandro.it/app/utils"
)

// Interface Peer base, so the common features to each type of peer.
type Peer struct {
	index int
}

/* Constant values */
const MAX_QUEUE = 100
const MAX_DELAY = 3

/* Global variables */
var conf *utils.Configuration
var conn *utils.Connection

// Initialization of features of Peer
func (p Peer) init_peer() {
	p.index = 0
}

// This function allows to register the node to communicate in multicast group.
func register_into_group() {
	var whoami_to_register utils.Whoami
	var empty utils.Empty

	build_whoami_struct(&whoami_to_register)

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 1234
	addr_register_node := "10.5.0.254:1234"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	utils.Check_error(err)

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.Register_node", &whoami_to_register, &empty)
	utils.Check_error(err)

	client.Close()
}

// This function log message into file: this has the value of delivery to application layer.
func log_message(pkt *utils.Packet) {
	// Open file into volume docker
	path_file := "/docker/node_volume/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.Check_error(err)

	_, err = f.WriteString(pkt.Timestamp.Format(time.RFC1123)[17:25] + ";" + pkt.Source_address + ";" + pkt.Message + "\n")
	utils.Check_error(err)

	f.Close()
}

// This function has the goal to clear the shell and print all messages received and sended by the current peer
func print_chat() {
	// Clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Print all messages, received and sended
	content, err := ioutil.ReadFile("/docker/node_volume/" + getIpAddress() + "_log.txt")
	utils.Check_error(err)
	list := string(content)
	print(list)
}

// This function, after reception of list from register node, allow to setup connection with each node of group multicast
func setup_connection(p *Peer) error {
	var err error

	for i := 0; i < conf.Nodes; i++ {
		addr_node := conn.Addresses[i] + ":1234"
		conn.Peer[i], err = rpc.Dial("tcp", addr_node)
		utils.Check_error(err)
		if conn.Addresses[i] == getIpAddress() {
			p.index = i
		}
	}

	return nil
}

// This RPC method of Node allow to get list from the registered node
func (p *Peer) Get_list(list *utils.List_of_nodes, reply *utils.Empty) error {
	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < conf.Nodes; i++ {
		conn.Addresses[i] = addr_tmp[i]
	}

	conn.Channel_connection <- true

	return nil
}

// This function allow to manage the log file to frontend
func (p *Peer) Handshake(emoty *utils.Empty, reply *utils.Hand_reply) error {
	reply.Ip_address = getIpAddress()

	return nil
}

// This function allow to init the information valid to manage the system
func init_configuration() {
	algo, _ := strconv.Atoi(os.Getenv("ALGORITHM"))
	nodes, _ := strconv.Atoi(os.Getenv("NODES"))

	conf = &utils.Configuration{Algorithm: algo, Nodes: nodes}

	conn = new(utils.Connection)
	conn.Init_connection(nodes)
}

func main() {
	// Init phase
	init_configuration()

	// Create file for log of messages
	f, err := os.Create("/docker/node_volume/" + getIpAddress() + "_log.txt")
	utils.Check_error(err)
	defer f.Close()

	// The node communicates with the register node to register his info
	register_into_group()

	// Register the service that allow the communication with frontend
	receiver := rpc.NewServer()

	peer_base := new(Peer)
	err = receiver.RegisterName("General", peer_base)
	utils.Check_error(err)

	// Allocate object to use it into program execution
	if conf.Algorithm == 1 {
		peer_1 := new(Peer_1)
		peer_1.init_peer_1()

		err = receiver.RegisterName("Peer", peer_1)
		utils.Check_error(err)

		// Listen for incoming messages on port 1234
		lis, err := net.Listen("tcp", ":1234")
		utils.Check_error(err)
		defer lis.Close()

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(lis)

		// Setup the connection with the peer of group multicast after the reception of list
		<-conn.Channel_connection

		setup_connection(&peer_1.peer)
		go peer_1.deliver_packet()

	} else if conf.Algorithm == 2 {
		peer_2 := new(Peer_2)
		peer_2.init_peer_2()

		err := receiver.RegisterName("Peer", peer_2)
		utils.Check_error(err)

		// Listen for incoming messages on port 1234
		lis, err := net.Listen("tcp", ":1234")
		utils.Check_error(err)
		defer lis.Close()

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(lis)

		// Setup the connection with the peer of group multicast after the reception of list
		<-conn.Channel_connection

		setup_connection(&peer_2.peer)
		go peer_2.deliver_packet()

	} else if conf.Algorithm == 3 {
		peer_3 := new(Peer_3)
		peer_3.init_peer_3()
		peer_3.vector_clock.Init(conf.Nodes)

		err := receiver.RegisterName("Peer", peer_3)
		utils.Check_error(err)

		// Listen for incoming messages on port 1234
		lis, err := net.Listen("tcp", ":1234")
		utils.Check_error(err)
		defer lis.Close()

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(lis)

		// Setup the connection with the peer of group multicast after the reception of list
		<-conn.Channel_connection

		setup_connection(&peer_3.peer)
		go peer_3.deliver_packet()
	}

	for {

	}
}
