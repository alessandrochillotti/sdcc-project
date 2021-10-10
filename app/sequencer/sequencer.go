/*
	This is the sequencer node that allow fully ordered multicast implemented centrally.
	It has ip address equal to 10.5.0.253 and it is listening in port 1234.
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

// Global variables
var f *os.FileMode
var registered_nodes = 0
var current_id = 0

/*
	This function send a specific message to each node of group multicast.
*/
func send_multicast_message(ip_address string, arg *lib.Packet, empty *lib.Empty) error {
	// Prepare packet to send
	pkt_seq := lib.Packet_sequencer{Id: current_id, Pkt: *arg}

	//Compute address destination
	addr_node := ip_address + ":1234"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_node)
	if err != nil {
		log.Println("Error in dialing: ", err)
		return err
	}
	defer client.Close()

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Node.Get_Message", &pkt_seq, &empty)
	if err != nil {
		log.Fatal("Error in Node.Get_Message: ", err)
		return err
	}

	return nil
}

/* This function is called by each generic node to send packet to each node of group multicast */
func (reg *Sequencer) Send_packet(arg *lib.Packet, empty *lib.Empty) error {
	// Open file
	file, err := os.Open("/docker/register_volume/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file line by line, so scan every ip address
	scanner := bufio.NewScanner(file)
	current_id = current_id + 1

	// Send to each node of group multicast the message
	for scanner.Scan() {
		go send_multicast_message(scanner.Text(), arg, empty)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func main() {
	reg := new(Sequencer)

	// Register a sequencer methods
	sequencer := rpc.NewServer()
	err := sequencer.RegisterName("Sequencer", reg)
	if err != nil {
		fmt.Println("Format of service is not correct: ", err)
	}

	// Listen for incoming messages on port 1234
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Listen error: ", err)
	}

	// Use goroutine to implement a lightweight thread to manage new connection
	go sequencer.Accept(lis)

	for {
		/*
			Since that the partnership is static, the sequencer stay up and running for all time
			without looking if a peer is up or down.
		*/
	}
}
