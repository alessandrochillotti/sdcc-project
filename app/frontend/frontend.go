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

func main() {
	var selected_container int
	var text string
	var choice int
	var empty utils.Empty

	// Print menù
	list_container := get_list_container()
	fmt.Println("Insert the number of one of following containers:")
	fmt.Printf("%s", list_container)
	fmt.Scanf("%d\n", &selected_container)

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
		fmt.Println("3. Test the algorithms")
		fmt.Println("4. Exit")

		fmt.Scanf("%d\n", &choice)

		switch choice {
		case 1:
			in := bufio.NewReader(os.Stdin)
			text, err = in.ReadString('\n')
			text = strings.TrimSpace(text)

			peer.Go("Peer.Get_message_from_frontend", &text, &empty, nil)

			clear_shell()

			break
		case 2:
			clear_shell()

			print_log(path_file, verbose)

			break
		case 3:
			clear_shell()

			choice = 0
			for choice != 1 && choice != 2 && choice != 3 {
				fmt.Println("Select the test to perform:")
				fmt.Println("1. Only sender message multicast")
				fmt.Println("2. More sender message multicast")
				fmt.Println("3. Back")

				fmt.Scanf("%d\n", &choice)
			}

			if choice == 1 {
				err := peer.Call("Peer.Test_one_sender", &empty, &empty)
				check_error(err)
			} else if choice == 2 {
				// err := peer.Call("Peer.Get_message_from_frontend", &text, &empty)
				// check_error(err)
			}

			break
		case 4:
			return
		default:
			fmt.Println("Codice operativo non supportato.")
			break
		}

	}
}
