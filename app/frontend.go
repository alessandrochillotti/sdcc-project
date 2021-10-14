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
			fmt.Printf("[%s] %s -> %s\n", line[0], line[1], line[2])
		} else {
			fmt.Printf("%s -> %s\n", line[1], line[2])
		}
	}
	fmt.Println()

	log_file.Close()
}

func main() {
	var text string
	var choice int
	var empty utils.Empty
	var hand_reply utils.Hand_reply

	// Print menù
	list_container := get_list_container()
	fmt.Println("Insert the number of one of following containers:")
	fmt.Printf("%s", list_container)
	fmt.Scanf("%d\n", &selected_container)

	// Dial of peer
	addr_node := "127.0.0.1:" + get_free_port(selected_container)
	client, err := rpc.Dial("tcp", addr_node)
	utils.Check_error(err)

	// Handshake with peer
	err = client.Call("General.Handshake", &empty, &hand_reply)
	check_error(err)

	// Prepare information to print log
	path_file := "./volumes/log_node/" + hand_reply.Ip_address + "_log.txt"
	verbose := catch_verbose_flag()

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

			client.Go("Peer.Get_message_from_frontend", &text, &empty, nil)

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
