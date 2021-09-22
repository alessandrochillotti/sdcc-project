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

	"alessandro.it/app/lib"
)

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
	log.Printf("Synchronous call to RPC server")
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
	log.Printf("Synchronous call to RPC server")
	err = client.Call("Register.List_of_nodes", "foo", &res)
	if err != nil {
		log.Fatal("Error in Register.Register_node: ", err)
	}
}

func main() {
	// For first thing, the node communicates with the register node to register his info
	register_node()
}
