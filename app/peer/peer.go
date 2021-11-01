/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"alessandro.it/app/utils"
)

// Interface Peer base, so the common features to each type of peer
type Peer struct {
	Index      int
	Ip_address string
	Port       int
	Username   string
}

/* Constant values */
const MAX_QUEUE = 100
const MAX_DELAY = 3
const PORT = 1234

/* Global variables */
var conf *utils.Configuration
var conn *utils.Connection
var channel_connection chan bool
var channel_handshake chan bool

// This function allows to register the node to communicate in multicast group.
func (p *Peer) register_into_group() {
	var whoami_to_register utils.Whoami
	var empty utils.Empty

	whoami_to_register.Ip_address = getIpAddress()
	whoami_to_register.Port = "1234"
	whoami_to_register.Username = p.Username

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

// This function, after reception of list from register node, allow to setup connection with each node of group multicast
func setup_connection(p *Peer) error {
	var err error

	for i := 0; i < conf.Nodes; i++ {
		addr_node := conn.Addresses[i] + ":1234"
		conn.Peer[i], err = rpc.Dial("tcp", addr_node)
		utils.Check_error(err)

		conn.Index[conn.Addresses[i]] = i

		if conn.Addresses[i] == p.Ip_address {
			p.Index = i
		}
	}

	return nil
}

// This RPC method of Node allow to get list from the registered node
func (p *Peer) Get_list(list *utils.List_of_nodes, reply *utils.Empty) error {
	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < conf.Nodes; i++ {
		mapping := strings.Split(addr_tmp[i], ";")
		conn.Addresses[i] = mapping[0]
		conn.Username[mapping[0]] = mapping[1]
	}

	channel_connection <- true

	return nil
}

// This function allow to init the information valid to manage the system
func init_configuration() {
	algo, _ := strconv.Atoi(os.Getenv("ALGORITHM"))
	nodes, _ := strconv.Atoi(os.Getenv("NODES"))

	conf = &utils.Configuration{Algorithm: algo, Nodes: nodes}

	conn = new(utils.Connection)
	conn.Init_connection(nodes)

	// Build channel
	channel_connection = make(chan bool)
	channel_handshake = make(chan bool)
}

/*
This function, after 30 seconds without the arrival of new messages,
closes all active connections and the application process is killed.
*/
func (p *Peer) manage_connection() {
	// Wait timer
	<-conf.Timer.C

	for i := 0; i < conf.Nodes; i++ {
		conn.Peer[i].Close()
	}

	os.Exit(0)
}

func main() {
	hand_peer := new(Handshake)
	init_configuration()

	// Handshake phase
	listener_handshake := hand_peer.frontend_handshake()
	defer (*listener_handshake).Close()

	// Wait the end of handshake phase
	<-channel_handshake

	// Build a new peer
	peer_base := &Peer{Index: hand_peer.New_peer.Index, Ip_address: hand_peer.New_peer.Ip_address, Port: hand_peer.New_peer.Port, Username: hand_peer.New_peer.Username}

	// Register the base services of general Peer
	receiver := rpc.NewServer()
	err := receiver.RegisterName("General", peer_base)
	utils.Check_error(err)
	listener, err := net.Listen("tcp", ":1234")
	utils.Check_error(err)
	defer listener.Close()

	// The node communicates with the recorder node for registration in the multicast group
	peer_base.register_into_group()

	// Create file for log of messages
	f, err := os.Create("/docker/node_volume/" + peer_base.Ip_address + "_log.txt")
	utils.Check_error(err)
	defer f.Close()

	conf.Timer = time.NewTimer(time.Duration(utils.TIMER_NODE*conf.Nodes) * time.Second)
	go peer_base.manage_connection()

	// Allocate object to use it into program execution
	if conf.Algorithm == 1 {
		peer_1 := &Peer_1{Peer: *peer_base}
		peer_1.init_peer_1(peer_base.Username)

		err = receiver.RegisterName("Peer", peer_1)
		utils.Check_error(err)

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(listener)

		// Setup the connection with the peer of group multicast after the reception of list
		<-channel_connection

		setup_connection(&peer_1.Peer)

		peer_1.deliver_packet()

	} else if conf.Algorithm == 2 {
		peer_2 := &Peer_2{Peer: *peer_base}
		peer_2.init_peer_2(peer_base.Username)

		err := receiver.RegisterName("Peer", peer_2)
		utils.Check_error(err)

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(listener)

		// Setup the connection with the peer of group multicast after the reception of list
		<-channel_connection

		setup_connection(&peer_2.Peer)

		peer_2.deliver_packet()

	} else if conf.Algorithm == 3 {
		peer_3 := &Peer_3{Peer: *peer_base}
		peer_3.init_peer_3()
		peer_3.vector_clock.Init(conf.Nodes)

		err := receiver.RegisterName("Peer", peer_3)
		utils.Check_error(err)

		// Use goroutine to implement a lightweight thread to manage the coming of new messages
		go receiver.Accept(listener)

		// Setup the connection with the peer of group multicast after the reception of list
		<-channel_connection

		setup_connection(&peer_3.Peer)

		peer_3.deliver_packet()
	}
}
