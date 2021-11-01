/*
	This is the special node that allow to register each node in the network.
	It has ip address equal to 10.5.0.254 and it is listening in port 1234.
*/

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"alessandro.it/app/utils"
)

type Register struct{}

var mutex_write sync.Mutex
var chan_reg chan (bool)

/*
This function is called by each generic node to:
	1. Register its ip address into group multicast
	2. When the number of node is equal to NUMBER_NODES, the register_node send the list of nodes
*/
func (reg *Register) Register_node(arg *utils.Whoami, empty *utils.Empty) error {
	// Open file into volume docker
	f, err := os.OpenFile("/docker/register_volume/nodes.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// Write into file the ip address of registered node
	mutex_write.Lock()
	_, err = f.WriteString(arg.Ip_address + ";" + arg.Username + "\n")
	if err != nil {
		log.Println(err)
		return err
	}
	mutex_write.Unlock()

	chan_reg <- true

	return nil
}

// This function allow to send to each node of group multicast the list of nodes registered.
func send_list() {
	var list_nodes utils.List_of_nodes
	var empty utils.Empty

	// Read whole file
	nodes, err := ioutil.ReadFile("/docker/register_volume/nodes.txt")
	if err != nil {
		log.Printf("Unable to read file: %v", err)
	}
	list_nodes.List_str = string(nodes)

	algo, _ := strconv.Atoi(os.Getenv("ALGORITHM"))
	if algo == 1 {
		// Send list of node to sequencer
		addr_node := "10.5.0.253:1234"
		client_sequencer, err := rpc.Dial("tcp", addr_node)
		if err != nil {
			log.Println("Error in dialing: ", err)
		}

		// Call remote procedure and reply will store the RPC result
		err = client_sequencer.Call("Sequencer.Get_list", &list_nodes, &empty)
		if err != nil {
			log.Fatal("Error in General.Get_list: ", err)
		}

		client_sequencer.Close()
	} else {
		// Open file
		file, err := os.Open("/docker/register_volume/nodes.txt")
		if err != nil {
			log.Fatal(err)
		}

		// Read file line by line, so scan every ip address
		scanner := bufio.NewScanner(file)

		// Send to each node of group multicast the message
		for scanner.Scan() {
			data := strings.Split(scanner.Text(), ";")
			addr_node := data[0] + ":1234"
			client, err := rpc.Dial("tcp", addr_node)
			if err != nil {
				log.Println("Error in dialing: ", err)
			}

			// Call remote procedure and reply will store the RPC result
			err = client.Call("General.Get_list", &list_nodes, &empty)
			if err != nil {
				log.Fatal("Error in General.Get_list: ", err)
			}

			client.Close()
		}

		file.Close()
	}
}

/*
This function, after (10 * nodes) seconds without the arrival of new messages, close the application
*/
func manage_connection(nodes int) {
	quit_timer := time.NewTimer(time.Duration(10*nodes) * time.Second)

	// Wait timer
	<-quit_timer.C

	os.Exit(0)
}

func main() {
	// Build useful structures
	chan_reg = make(chan bool)

	nodes, _ := strconv.Atoi(os.Getenv("NODES"))

	reg := new(Register)

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/docker/register_volume/nodes.txt")
	utils.Check_error(err)
	f.Close()

	// Wait timer
	go manage_connection(nodes)

	// Register a new RPC server and the struct we created above
	server := rpc.NewServer()
	err = server.RegisterName("Register", reg)
	utils.Check_error(err)

	// Listen for incoming messages on port 1234
	lis, err := net.Listen("tcp", ":1234")
	utils.Check_error(err)
	go server.Accept(lis)

	// Wait that every nodes is registered
	for i := 0; i < nodes; i++ {
		<-chan_reg
	}

	// Once that every node is registered, then the register node send the list of nodes to each node of group multicast
	send_list()
}
