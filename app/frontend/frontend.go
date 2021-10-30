/*
This file is an application that allow to join into multicast program.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"

	"alessandro.it/app/utils"
)

// This function check if the flag verbose is setted.
func catch_verbose_flag() bool {
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-V" {
			return true
		}
	}

	return false
}

// This function print the log of files.
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

func main() {
	var selected_container int
	var text string
	var choice int
	var empty utils.Empty

	// Print menù
	list_container := get_list_container()
	fmt.Println("Insert the number of one of following containers:")
	fmt.Printf("%s", list_container)

	number_container := len(strings.Split(list_container, "\n")) - 1

	// Replicate do while
	fmt.Scanf("%d\n", &selected_container)
	for selected_container < 1 || selected_container > number_container {
		fmt.Println("Insert a number between 1 to", number_container)
		fmt.Scanf("%d\n", &selected_container)
	}

	ip_addr, username, _ := handshake(selected_container, false)

	// Dial of peer
	addr_node := "127.0.0.1:" + get_local_port(selected_container, (uint16(1234)))
	peer, err := rpc.Dial("tcp", addr_node)
	utils.Check_error(err)

	// Prepare information to print log
	path_file := "./volumes/log_node/" + ip_addr + "_log.txt"
	verbose := catch_verbose_flag()

	// Clear the shell
	clear_shell()

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

			msg_to_send := &utils.Message{Text: text, Delay: make([]int, 0)}

			peer.Go("Peer.Get_message_from_frontend", msg_to_send, &empty, nil)

			clear_shell()

			fmt.Println("Message sent\n")

			break
		case 2:
			clear_shell()

			print_log(path_file, verbose)

			break
		case 3:
			return
		default:
			clear_shell()

			fmt.Println("Operation code not valid\n")
			break
		}

	}
}
