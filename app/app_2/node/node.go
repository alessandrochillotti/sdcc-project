/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
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
			head_node := queue.Get_head().Update
			// Deliver the packet to application layer
			log_message(&head_node.Packet, head_node.Timestamp)

			// Clear shell
			// cmd := exec.Command("clear")
			// cmd.Stdout = os.Stdout
			// cmd.Run()

			// Print chat
			content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
			lib.Check_error(err)

			list := string(content)

			print(list)
		}
	}
}

/*
This function allow to get update from the other node of group multicast
*/
func (node *Node) Get_update(update *utils.Update, ack *utils.Ack) error {
	scalar_clock = lib.Max(scalar_clock, update.Timestamp)
	scalar_clock = scalar_clock + 1

	// Put update in queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 1}
	queue.Update_into_queue(update_node)

	// Send ack message in multicast
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addr_node := addresses[i] + ":1234"
		go send_ack(addr_node, update.Packet.Id)
	}

	return nil
}

func (node *Node) Get_ack(id *int, empty *lib.Empty) error {
	if queue.Ack_node(*id) == true {
		return nil
	} else {
		err := errors.New("Element to acked not found")
		return err
	}
}

func send_ack(addr_node string, id int) error {
	var empty lib.Empty

	// Try to connect to node
	client, err := rpc.Dial("tcp", addr_node)
	if err != nil {
		log.Println("Error in dialing: ", err)
		return err
	}
	defer client.Close()

	err = client.Call("Node.Get_ack", &id, &empty)
	lib.Check_error(err)

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
		pkt := lib.Packet{Id: queue.Max_id + 1, Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}
		fmt.Println("Ho assegnato id = ", pkt.Id)

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

	queue = &utils.Queue{Max_id: 0}

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
