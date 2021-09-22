/*
	This is the special node that allow to register each node in the network.
	It has ip address equal to 10.5.0.254 and it is listening in port 4321.
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

	"alessandro.it/app/lib"
)

type Register int

var f *os.FileMode
var registered_nodes = 0

func check_error(e error, res *lib.Outcome) {
	if e != nil {
		*res = false
		os.Exit(1)
	}
}

func (reg *Register) Register_node(arg *lib.Whoami, res *lib.Outcome) error {

	f, err := os.OpenFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(arg.Ip_address + "\n"); err != nil {
		log.Println(err)
	}

	*res = true
	registered_nodes = registered_nodes + 1
	fmt.Printf("The registration is for the ip address : %s\n", arg.Ip_address)

	return nil
}

func send_list_registered_nodes() {
	var res lib.Outcome

	// Read whole file
	content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	list := string(content)

	// Open file
	file, err := os.Open("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file line by line, so scan every ip address
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Prepare packet to send
		pkt := lib.Packet{Source_address: "10.5.0.254", Message: list}

		//Compute address destination
		addr_node := scanner.Text() + ":4321"

		fmt.Printf("Send list to %s\n", addr_node)

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_node)
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		defer client.Close()

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Node.Get_List_Nodes", &pkt, &res)
		if err != nil {
			log.Fatal("Error in Node.Get_List_Nodes: ", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func main() {

	reg := new(Register)

	// Register a new RPC server and the struct we created above
	server := rpc.NewServer()
	err := server.RegisterName("Register", reg)
	if err != nil {
		fmt.Println("Format of service is not correct: ", err)
	}

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		fmt.Println("Listen error: ", err)
	}

	// Use goroutine to implement a lightweight thread to manage new connection
	go server.Accept(lis)

	for {
		if registered_nodes == lib.NUMBER_NODES {
			send_list_registered_nodes()
			lis.Close()
			break
		}
	}
}
