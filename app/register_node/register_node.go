/*
	This is the special node that allow to register each node in the network.
	It has ip address equal to 10.5.0.254 and it is listening in port 4321.
*/

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"

	"alessandro.it/app/lib"
)

type Register int

var f *os.FileMode
var mutex_write sync.Mutex
var chan_reg chan (bool)

/*
This function is called by each generic node to:
	1. Register its ip address into group multicast
	2. When the number of node is equal to NUMBER_NODES, the register_node send the list of nodes
*/
func (reg *Register) Register_node(arg *lib.Whoami, empty *lib.Empty) error {
	// Open file into volume docker
	f, err := os.OpenFile("/docker/register_volume/nodes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// Write into file the ip address of registered node
	mutex_write.Lock()
	_, err = f.WriteString(arg.Ip_address + "\n")
	if err != nil {
		log.Println(err)
		return err
	}
	mutex_write.Unlock()

	chan_reg <- true

	return nil
}

/*
This function allow to send to each node of group multicast the list of nodes registered.
*/
func send_list() {
	var list_nodes lib.List_of_nodes
	var empty lib.Empty

	// Read whole file
	nodes, err := ioutil.ReadFile("/docker/register_volume/nodes.txt")
	if err != nil {
		log.Printf("Unable to read file: %v", err)
	}
	list_nodes.List_str = string(nodes)

	// Open file
	file, err := os.Open("/docker/register_volume/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file line by line, so scan every ip address
	scanner := bufio.NewScanner(file)

	// Send to each node of group multicast the message
	for scanner.Scan() {
		addr_node := scanner.Text() + ":1234"
		client, err := rpc.Dial("tcp", addr_node)
		if err != nil {
			log.Println("Error in dialing: ", err)
		}

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Node.Get_list", &list_nodes, &empty)
		if err != nil {
			log.Fatal("Error in Node.Get_list: ", err)
		}

		client.Close()
	}

}

func main() {
	reg := new(Register)

	// Register a new RPC server and the struct we created above
	server := rpc.NewServer()
	err := server.RegisterName("Register", reg)
	lib.Check_error(err)

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/docker/register_volume/nodes.txt")
	lib.Check_error(err)
	f.Close()

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	lib.Check_error(err)

	chan_reg = make(chan bool)
	go server.Accept(lis)

	defer lis.Close()

	// Wait that every nodes is registered
	for i := 0; i < lib.NUMBER_NODES; i++ {
		<-chan_reg
	}

	// Once that every node is registered, then the register node send the list of nodes to each node of group multicast
	send_list()

	os.Exit(0)
}
