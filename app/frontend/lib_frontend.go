package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"alessandro.it/app/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clear_shell() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func get_local_port(index int, port_number uint16) string {
	var port uint16
	port = 0

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	cnt := 0
	for _, container := range containers {
		cnt++
		for i := 0; cnt == index && i < len(container.Ports) && port == 0; i++ {
			if container.Ports[i].PrivatePort == port_number {
				port = container.Ports[i].PublicPort
			}
		}
	}

	return strconv.Itoa(int(port))
}

func get_list_container() string {
	var list string

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	cnt := 0
	for _, container := range containers {
		if container.Names[0][1:10] == "app_peer_" {
			list = list + fmt.Sprintf("%d. %s\n", cnt+1, container.Names[0][1:])
			cnt++
		}
	}

	if cnt == 0 {
		fmt.Println("There are no containers.")
		os.Exit(0)
	}

	return list
}

func handshake(container int, test bool) (string, string, int) {
	var reply *utils.Hand_reply

	addr_node := "127.0.0.1:" + get_local_port(container, (uint16(4444)))

	peer_handshake, err := rpc.Dial("tcp", addr_node)
	if err != nil {
		log.Println("Error in dialing: ", err)
	}

	fmt.Println("Choose a username")
	in := bufio.NewReader(os.Stdin)
	username, err := in.ReadString('\n')
	username = strings.TrimSpace(username)

	// Call remote procedure and reply will store the RPC result
	request := &utils.Hand_request{Username: username, Test: test}
	err = peer_handshake.Call("Handshake.Handshake", &request, &reply)
	clear_shell()
	check_error(err)

	peer_handshake.Close()

	return reply.Ip_address, username, reply.Algorithm
}

func check_error(err error) {
	if err != nil {
		fmt.Println("Something went wrong. Retry.")
		os.Exit(1)
	}
}
