/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"bufio"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strings"

	"alessandro.it/app/lib"
)

type Node int

// Constant values
const MAX_PACKET_BUFFERED = 100

// Global variables
var addresses [lib.NUMBER_NODES]string /* Contains ip addresses of each node in multicast group */
var current_id = 0
var buffer chan (lib.Packet_sequencer)

/* This function return the ip address of current node */
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

/* This function build a struct that contains the info to register the node */
func build_whoami_struct(whoami_to_register *lib.Whoami) {
	whoami_to_register.Ip_address = getIpAddress()
	whoami_to_register.Port = "1234"
}

/* This function allows to register the node to communicate in multicast group */
func register_into_group() {
	var whoami_to_register lib.Whoami

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 4321
	addr_register_node := "10.5.0.254:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	lib.Check_error(err)

	defer client.Close()

	build_whoami_struct(&whoami_to_register)

	// Call remote procedure and reply will store the RPC result
	client.Call("Register.Register_node", &whoami_to_register, &addresses)
}

/* This function log message into file  */
func log_message(pkt *lib.Packet) {
	// Open file into volume docker
	path_file := "/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	lib.Check_error(err)
	defer f.Close()

	// Write into file the ip address of registered node
	_, err = f.WriteString(pkt.Source_address + " -> " + pkt.Message + "\n")
	lib.Check_error(err)
}

func check_buffered_packet() {
	for {
		current_packet := <-buffer
		if current_id+1 == current_packet.Id { // If the packet is the expected packet
			current_id = current_id + 1
			log_message(&current_packet.Pkt)

			// Clear shell
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			// Print chat
			content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
			lib.Check_error(err)

			list := string(content)

			print(list)
		} else {
			buffer <- current_packet
		}
	}
}

/* This function is called by sequencer node for sending message */
func (node *Node) Get_Message(pkt *lib.Packet_sequencer, res *lib.Outcome) error {
	if current_id+1 == pkt.Id { // If the packet is the expected packet
		current_id = current_id + 1
		log_message(&pkt.Pkt)

		// Clear shell
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		// Print chat
		content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
		lib.Check_error(err)

		list := string(content)

		print(list)
	} else {
		// Buffered packet
		buffer <- *pkt
	}

	*res = true

	return nil
}

/* This function allow to wait the input of user to send the message to each node of group multicast */
func open_standard_input() {
	var res lib.Outcome

	for {
		// Take in input the content of message to send
		in := bufio.NewReader(os.Stdin)
		text, err := in.ReadString('\n')
		text = strings.TrimSpace(text)

		// Build packet to send
		pkt := lib.Packet{Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}

		// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 4321
		addr_sequencer_node := "10.5.0.253:4321"

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_sequencer_node)
		lib.Check_error(err)

		defer client.Close()

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Sequencer.Send_packet", &pkt, &res)
		lib.Check_error(err)
	}
}

func main() {
	// For first thing, the node communicates with the register node to register his info
	register_into_group()

	node := new(Node)
	buffer = make(chan lib.Packet_sequencer)

	// Create file for log of messages
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
	lib.Check_error(err)

	defer f.Close()

	// Register the Node methods
	receiver := rpc.NewServer()
	err = receiver.RegisterName("Node", node)
	lib.Check_error(err)

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	lib.Check_error(err)

	// Use goroutine to implement a lightweight thread to manage the coming of new messages
	go receiver.Accept(lis)

	// go check_buffered_packet()

	open_standard_input()
}
