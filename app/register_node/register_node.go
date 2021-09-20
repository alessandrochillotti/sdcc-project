/*
	This is the special node that allow to register each node in the network.
	It has ip address equal to 10.5.0.254 and it is listening in port 4321.
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

type Register int

var f *os.File

func check_error(e error, res *lib.Outcome) {
	if e != nil {
		*res = false
		os.Exit(1)
	}
}

func (reg *Register) Register_node(arg *lib.Whoami, res *lib.Outcome) error {

	// d1 := []byte(arg.Ip_address + "\t" + arg.Port + "\n")

	// err := os.WriteFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt", d1, 0644)
	// if err != nil {
	// 	panic(err)
	// }

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
	fmt.Printf("The registration is for the ip address : %s\n", arg.Ip_address)

	return nil
}

func (reg *Register) List_of_nodes(msg string, res *lib.Outcome) error {
	fmt.Printf("Dummy list \n")

	*res = true

	return nil
}

func main() {

	reg := new(Register)

	// Register a new RPC server and the struct we created above
	server := rpc.NewServer()
	err := server.RegisterName("Register", reg)
	if err != nil {
		fmt.Printf("Format of service is not correct: ", err)
	}

	// Create file to maintain ip address and number port of all registered nodes
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/nodes.txt")
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// Listen for incoming messages on port 4321
	lis, err := net.Listen("tcp", ":4321")
	if err != nil {
		fmt.Printf("Listen error: ", err)
	}
	fmt.Printf("RPC server on port %d\n", 4321)

	server.Accept(lis)
}
