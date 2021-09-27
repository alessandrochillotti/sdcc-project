/*
	This is the sequencer node that allow fully ordered multicast implemented centrally.
	It has ip address equal to 10.5.0.253 and it is listening in port 4321.
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"alessandro.it/app/lib"
)

type Sequencer int

var f *os.FileMode

// Global variables
var registered_nodes = 0
var current_id = 0

/* This function is called by each generic node to send packet to each node of group multicast */
func (reg *Sequencer) Send_packet(arg *lib.Packet, res *lib.Outcome) error {
	// Open file
	file, err := os.Open("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file line by line, so scan every ip address
	scanner := bufio.NewScanner(file)
	current_id = current_id + 1

	for scanner.Scan() {
		// Prepare packet to send
		pkt_seq := lib.Packet_sequencer{Id: current_id, Pkt: *arg}

		//Compute address destination
		addr_node := scanner.Text() + ":4321"

		// Try to connect to addr_register_node
		client, err := rpc.Dial("tcp", addr_node)
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		defer client.Close()

		// Call remote procedure and reply will store the RPC result
		err = client.Call("Node.Get_Message", &pkt_seq, &res)
		if err != nil {
			log.Fatal("Error in Node.Get_Message: ", err)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ricevuto %s da %s\n", arg.Message, arg.Source_address)

	return nil
}

func main() {

	reg := new(Sequencer)

	// Register a new RPC server and the struct we created above
	sequencer := rpc.NewServer()
	err := sequencer.RegisterName("Sequencer", reg)
	if err != nil {
		fmt.Println("Format of service is not correct: ", err)
	}

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		fmt.Println("Listen error: ", err)
	}

	// Use goroutine to implement a lightweight thread to manage new connection
	go sequencer.Accept(lis)

	for {
		// TODO: implement a control that if nobody is up, then sequencer exit
	}
}