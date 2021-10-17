package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"alessandro.it/app/utils"
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

func catch_verbose_flag() bool {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-V" {
			return true
		}
	}

	return false
}

func print_log(path_file string, verbose bool) {
	log_file, err := os.Open(path_file)
	if err != nil {
		log.Fatal(err)
	}

	scanner_log := bufio.NewScanner(log_file)
	if err := scanner_log.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delivered messages")
	for scanner_log.Scan() {
		line := strings.Split(scanner_log.Text(), ";")
		if verbose {
			fmt.Printf("[%s] %s -> %s [%s]\n", line[1], line[2], line[3], line[0])
		} else {
			fmt.Printf("%s -> %s\n", line[2], line[3])
		}
	}
	fmt.Println()

	log_file.Close()
}

func handshake(container int) (string, string) {
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
	request := &utils.Hand_request{Username: username}
	err = peer_handshake.Call("Handshake.Handshake", &request, &reply)
	if err != nil {
		log.Fatal("Error in General.Get_list: ", err)
	}

	peer_handshake.Close()

	return reply.Ip_address, username
}

func main() {
	var text string
	var choice int
	var empty utils.Empty

	// Print menù
	list_container := get_list_container()
	fmt.Println("Insert the number of one of following containers:")
	fmt.Printf("%s", list_container)
	fmt.Scanf("%d\n", &selected_container)

	ip_addr, username := handshake(selected_container)

	// Dial of peer
	addr_node := "127.0.0.1:" + get_local_port(selected_container, (uint16(1234)))
	peer, err := rpc.Dial("tcp", addr_node)
	utils.Check_error(err)

	// Prepare information to print log
	path_file := "./volumes/log_node/" + ip_addr + "_log.txt"
	verbose := catch_verbose_flag()

	// Clear the shell
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for {
		fmt.Println("Welcome", username)
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

			peer.Go("Peer.Get_message_from_frontend", &text, &empty, nil)

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

			print_log(path_file, verbose)

			break
		case 3:
			return
		default:
			fmt.Println("Codice operativo non supportato.")
			break
		}

	}
}
