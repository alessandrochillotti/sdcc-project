/*
	This is the special node that allow to register each node in the network.
	It has ip address equal to 10.5.0.254 and it is listening in port 4321.
*/

package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"alessandro.it/app/lib"
)

type Register int

var f *os.FileMode
var chan_reg chan bool
var chan_send chan bool
var chan_exit chan bool

/*
This function is called by each generic node to:
	1. Register its ip address into group multicast
	2. When the number of node is equal to NUMBER_NODES, the register_node send the list of nodes
*/
func (reg *Register) Register_node(arg *lib.Whoami, addresses *lib.Addresses) error {

	// Open file into volume docker
	f, err := os.OpenFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// Write into file the ip address of registered node
	_, err = f.WriteString(arg.Ip_address + "\n")
	if err != nil {
		log.Println(err)
		return err
	}

	// Communicate to the main process that the registration is completed
	chan_reg <- true

	// Wait other goroutine
	<-chan_send

	// Read whole file
	content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Printf("Unable to read file: %v", err)
		return err
	}
	list_of_nodes := string(content)

	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list_of_nodes, "\n")
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addresses.Addresses_array[i] = addr_tmp[i]
	}

	// Communicate to the main process that goroutine has terminated its job
	chan_exit <- true

	return nil
}

func main() {
	reg := new(Register)

	// Register a new RPC server and the struct we created above
	server := rpc.NewServer()
	err := server.RegisterName("Register", reg)
	lib.Check_error(err)

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	lib.Check_error(err)
	f.Close()

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	lib.Check_error(err)

	defer lis.Close()

	// Those channels allow the syncronization between the goroutines
	chan_exit = make(chan bool)
	chan_send = make(chan bool)
	chan_reg = make(chan bool)

	// Use goroutine to implement a lightweight thread to manage new connection
	go server.Accept(lis)

	// Wait the registration
	for i := 0; i < lib.NUMBER_NODES; i++ {
		<-chan_reg
	}

	// Tell that the goroutine can send
	for i := 0; i < lib.NUMBER_NODES; i++ {
		chan_send <- true
	}

	// Wait the goroutines
	for i := 0; i < lib.NUMBER_NODES; i++ {
		<-chan_exit
	}

	os.Exit(0)
}
