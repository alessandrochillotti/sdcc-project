/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
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
func register_node() {
	var whoami_to_register lib.Whoami
	var res lib.Outcome

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 4321
	addr_register_node := "10.5.0.254:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	build_whoami_struct(&whoami_to_register)

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.Register_node", &whoami_to_register, &res)
	if err != nil {
		log.Fatal("Error in Register.Register_node: ", err)
	}

	// Print the outcome of registration phase
	if res == true {
		fmt.Printf("The registration phase is ok\n")
	} else {
		fmt.Printf("Errore in registraion phase\n")
	}
}

/* This function allows to retrieve the list of nodes in group of multicast */
func get_nodes_in_group() {
	var res lib.Outcome

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 4321
	addr_register_node := "10.5.0.254:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.List_of_nodes", "foo", &res)
	if err != nil {
		log.Fatal("Error in Register.Register_node: ", err)
	}
}

/* This function is called by register node for sending list of ip addresses and it load list into array */
func (node *Node) Get_List_Nodes(pkt *lib.Packet, res *lib.Outcome) error {

	// Parse the list and load it into array of ip
	addr_tmp := strings.Split(pkt.Message, "\n")
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addresses[i] = addr_tmp[i]
	}

	*res = true

	return nil
}

func log_message(pkt *lib.Packet) {
	// Open file into volume docker
	path_file := "/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	// Write into file the ip address of registered node
	if _, err := f.WriteString(pkt.Source_address + " -> " + pkt.Message + "\n"); err != nil {
		log.Println(err)
	}
}

/* This function is called by sequencer node for sending message */
func (node *Node) Get_Message(pkt *lib.Packet_sequencer, res *lib.Outcome) error {

	if current_id+1 == pkt.Id {
		current_id = current_id + 1
		log_message(&pkt.Pkt)
		fmt.Print("\033[H\033[2J")

		// Print chat
		content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		list := string(content)

		print(list)

		// print_msg(&pkt.Pkt)
	} else {
		// TODO: buffer packet
		buffer <- *pkt
		// TODO: implement a goroutine that check if there are packet buffered that contains next id
	}

	*res = true

	return nil
}

func send_packet() {
	var text string
	var res lib.Outcome

	// Take in input the content of message to send
	fmt.Println("Insert a text to send to each node of group multicast.")

	in := bufio.NewReader(os.Stdin)
	text, err := in.ReadString('\n')
	text = strings.TrimSpace(text)

	// Build packet to send
	pkt := lib.Packet{Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}

	// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 4321
	addr_sequencer_node := "10.5.0.253:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_sequencer_node)
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Sequencer.Send_packet", &pkt, &res)
	if err != nil {
		log.Fatal("Error in Register.Register_node: ", err)
	}
}

func receive_message() {
	node := new(Node)

	receiver := rpc.NewServer()
	err := receiver.RegisterName("Node", node)
	if err != nil {
		fmt.Println("Format of service is not correct: ", err)
	}

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		fmt.Println("Listen error: ", err)
	}

	// Use goroutine to implement a lightweight thread to manage new connection
	go receiver.Accept(lis)

	for {
		var text string
		var res lib.Outcome

		// Take in input the content of message to send
		// fmt.Println("Insert a text to send to each node of group multicast.")

		in := bufio.NewReader(os.Stdin)
		text, err := in.ReadString('\n')
		text = strings.TrimSpace(text)

		// Build packet to send
		pkt := lib.Packet{Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}

		// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 4321
		addr_sequencer_node := "10.5.0.253:4321"

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_sequencer_node)
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		defer client.Close()

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Sequencer.Send_packet", &pkt, &res)
		if err != nil {
			log.Fatal("Error in Register.Register_node: ", err)
		}
	}
}

func main() {
	// For first thing, the node communicates with the register node to register his info
	register_node()

	node := new(Node)

	receiver := rpc.NewServer()
	err := receiver.RegisterName("Node", node)
	if err != nil {
		fmt.Println("Format of service is not correct: ", err)
	}

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		fmt.Println("Listen error: ", err)
	}

	// Use goroutine to implement a lightweight thread to manage new connection
	go receiver.Accept(lis)

	for {
		var text string
		var res lib.Outcome

		in := bufio.NewReader(os.Stdin)
		text, err := in.ReadString('\n')
		text = strings.TrimSpace(text)

		// Build packet to send
		pkt := lib.Packet{Source_address: getIpAddress(), Source_pid: os.Getpid(), Message: text}

		// The sequencer node has ip address set to 10.5.0.253 and it is listening in port 4321
		addr_sequencer_node := "10.5.0.253:4321"

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_sequencer_node)
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		defer client.Close()

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Sequencer.Send_packet", &pkt, &res)
		if err != nil {
			log.Fatal("Error in Register.Register_node: ", err)
		}

	}
}
