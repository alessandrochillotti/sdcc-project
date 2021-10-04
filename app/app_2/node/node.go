/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"alessandro.it/app/lib"
	"alessandro.it/app/utils"
)

type Node int

// Constant values
const MAX_QUEUE = 100

// Global variables
var scalar_clock int = 0
var addresses [lib.NUMBER_NODES]string /* Contains ip addresses of each node in multicast group */
var queue *utils.Queue

/*
This function return the ip address of current node
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
This function build a struct that contains the info to register the node
*/
func build_whoami_struct(whoami_to_register *lib.Whoami) {
	whoami_to_register.Ip_address = getIpAddress()
	whoami_to_register.Port = "1234"
}

/*
This function allows to register the node to communicate in multicast group
*/
func register_into_group() {
	var whoami_to_register lib.Whoami
	var empty lib.Empty

	build_whoami_struct(&whoami_to_register)

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 4321
	addr_register_node := "10.5.0.254:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	lib.Check_error(err)
	defer client.Close()

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.Register_node", &whoami_to_register, &empty)
	lib.Check_error(err)

}

func (node *Node) Get_list(list *lib.List_of_nodes, reply *lib.Empty) error {
	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addresses[i] = addr_tmp[i]
	}

	return nil
}

/*
	Functions used to develop the algorith number 1:
*/

/*
This function log message into file
*/
func log_message(pkt *lib.Packet, id int) {
	// Open file into volume docker
	path_file := "/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	lib.Check_error(err)
	defer f.Close()

	// Write into file the ip address of registered node
	_, err = f.WriteString(pkt.Source_address + " -> " + pkt.Message + "[" + strconv.Itoa(id) + "]\n")
	lib.Check_error(err)
}

/*
This function check if there are packet to deliver.
*/
func deliver_packet() {
	for {
		if queue.Get_ack_head() == lib.NUMBER_NODES {
			fmt.Println("Consegnato", queue.Get_head().Update.Packet.Message)
			// TODO: consegnare il messaggio a livello applicativo e avvisare gli altri nodi
		}
	}
}

func (node *Node) Get_update(update *utils.Update, ack *utils.Ack) error {
	scalar_clock = lib.Max(scalar_clock, update.Timestamp)
	scalar_clock = scalar_clock + 1

	// Put update in queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 0}
	queue.Update_into_queue(update_node)

	// Send ack to sender of update message
	*ack = *ack + 1

	return nil
}

func send_update(addr_node string, update_node *utils.Node) error {
	// Try to connect to node
	client, err := rpc.Dial("tcp", addr_node)
	if err != nil {
		log.Println("Error in dialing: ", err)
		return err
	}
	defer client.Close()

	// Build update to send
	var ack utils.Ack = 0

	// Delay the send of update message
	// lib.Delay()

	err = client.Call("Node.Get_update", update_node.Update, &ack)
	lib.Check_error(err)

	fmt.Println("Ack ricevuto = ", ack)
	update_node.Ack = update_node.Ack + ack

	return nil
}

/*
This function allow to wait the input of user to send the message to each node of group multicast
*/
func open_standard_input() {
	for {
		// Take in input the content of message to send
		in := bufio.NewReader(os.Stdin)
		text, _ := in.ReadString('\n')
		text = strings.TrimSpace(text)

		// Build packet
		pkt := lib.Packet{Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}

		// Update the scalar clock
		scalar_clock = scalar_clock + 1

		// Build update to send
		update := utils.Update{Timestamp: scalar_clock, Packet: pkt}
		update_node := utils.Node{Update: update, Next: nil, Ack: 1}
		queue.Update_into_queue(&update_node)

		my_ip := getIpAddress()
		// Send to each node of group multicast the message
		for i := 0; i < lib.NUMBER_NODES; i++ {
			if addresses[i] != my_ip {
				addr_node := addresses[i] + ":1234"
				go send_update(addr_node, &update_node)
			}
		}
	}
}

func main() {
	// For first thing, the node communicates with the register node to register his info
	register_into_group()

	node := new(Node)

	queue = &utils.Queue{}

	// Create file for log of messages
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
	lib.Check_error(err)
	defer f.Close()

	// Register the Node methods
	receiver := rpc.NewServer()
	err = receiver.RegisterName("Node", node)
	lib.Check_error(err)

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":1234")
	lib.Check_error(err)

	// Use goroutine to implement a lightweight thread to manage the coming of new messages
	go receiver.Accept(lis)

	// This goroutine check always if there are packet to deliver
	go deliver_packet()

	// The user can insert text to send to each node of group multicast
	open_standard_input()
}
