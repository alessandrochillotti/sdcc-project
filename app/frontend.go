package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"alessandro.it/app/lib"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var selected_container int

func check_error(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func get_free_port(index int) string {
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
			if container.Ports[i].PrivatePort == 1234 {
				port = container.Ports[i].PublicPort
			}
		}
	}

	return strconv.Itoa(int(port))
}

func get_list_container() string {
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
		if container.Names[0][1:10] == "app_node_" {
			fmt.Printf("%d. %s\n", cnt+1, container.Names[0][1:])
			cnt++
		}
	}

	return strconv.Itoa(int(port))
}

func handshake(client *rpc.Client, hand_reply *lib.Hand_reply) {
	verbose := false

	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-V" {
			verbose = true
		}
	}

	handshake_packet := &lib.Hand_request{Verbose: verbose}

	err := client.Call("Node.Handshake", &handshake_packet, hand_reply)
	check_error(err)
}

func main() {
	var text string
	var choice int
	var empty lib.Empty
	var hand_reply lib.Hand_reply

	// Print menù
	fmt.Println("Insert the number of one of following containers:")
	get_list_container()
	fmt.Scanf("%d\n", &selected_container)

	// Dial of peer
	addr_node := "127.0.0.1:" + get_free_port(selected_container)
	client, err := rpc.Dial("tcp", addr_node)
	lib.Check_error(err)

	// Handshake with peer
	handshake(client, &hand_reply)

	path_file := "./volumes/log_node/" + hand_reply.Ip_address + "_log.txt"

	// Clear the shell
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for {
		fmt.Println("Welcome", hand_reply.Ip_address)
		fmt.Println("Insert the operation code:")
		fmt.Println("1. Send message")
		fmt.Println("2. Print messaged delivered")
		fmt.Println("3. Exit")

		fmt.Scanf("%d\n", &choice)

		switch choice {
		case 1:
			in := bufio.NewReader(os.Stdin)
			text, err = in.ReadString('\n')
			text = strings.TrimSpace(text)

			client.Go("Node.Get_message_from_frontend", &text, &empty, nil)

			// Clear the shell
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			break
		case 2:
			// Clear the shell
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			content, _ := ioutil.ReadFile(path_file)
			fmt.Println(string(content))

			break
		case 3:
			return
		default:
			fmt.Println("Codice operativo non supportato.")
			break
		}

	}
}
