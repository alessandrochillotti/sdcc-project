/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"alessandro.it/app/lib"
)

type Node int

var addresses [lib.NUMBER_NODES]string

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

func (node *Node) Get_List_Nodes(pkt *lib.Packet, res *lib.Outcome) error {

	// Parse the list and load it into array of ip
	addr_tmp := strings.Split(pkt.Message, "\n")
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addresses[i] = addr_tmp[i]
	}

	*res = true

	return nil
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
		// Arrivo del pacchetto di lista di nodi

	}
}
