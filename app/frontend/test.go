package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"alessandro.it/app/utils"
)

var peer []*rpc.Client

func check_log_equal(number int) bool {
	equal := true

	path_file := "./volumes/log_node/10.5.0.2_log.txt"
	file_1, err := ioutil.ReadFile(path_file)
	check_error(err)
	for i := 1; i < number && equal; i++ {
		path_file := "./volumes/log_node/10.5.0." + strconv.Itoa(i+2) + "_log.txt"
		file_2, err := ioutil.ReadFile(path_file)
		check_error(err)

		if string(file_1) == string(file_2) {
			file_1 = file_2
		} else {
			equal = false
		}
	}

	return equal
}

func check_log_empty(number int) bool {
	for i := 0; i < number; i++ {
		path_file := "./volumes/log_node/10.5.0." + strconv.Itoa(i+2) + "_log.txt"
		file, _ := os.Stat(path_file)

		if file.Size() != 0 {
			return false
		}
	}

	return true
}

func test_one_sender_1(number int) bool {
	var empty *utils.Empty

	message_1 := "messaggio"

	err := peer[0].Call("Peer.Get_message_from_frontend", &message_1, &empty)
	check_error(err)

	return check_log_equal(number)
}

func test_more_sender_1(number int) bool {
	var empty *utils.Empty

	for i := 0; i < number; i++ {
		message_1 := "messaggio " + strconv.Itoa(i)

		peer[i].Go("Peer.Get_message_from_frontend", &message_1, &empty, nil)
	}

	return check_log_equal(number)
}

func test_one_sender_2(number int) bool {
	var empty *utils.Empty

	message_1 := "messaggio"

	err := peer[0].Call("Peer.Get_message_from_frontend", &message_1, &empty)
	check_error(err)

	return check_log_empty(number)
}

func test_more_sender_2() bool {

	return true
}

func test_one_sender_3() bool {

	return true
}

func test_more_sender_3() bool {

	return true
}

func main() {
	var algo int
	var choice int

	// Print menù
	peer_number := len(strings.Split(get_list_container(), "\n")) - 1
	peer = make([]*rpc.Client, peer_number)
	fmt.Println(peer_number)
	// Dial of peer
	for i := 0; i < peer_number; i++ {
		var err error

		_, _, algo = handshake(i + 1)
		addr_node := "127.0.0.1:" + get_local_port(i+1, (uint16(1234)))
		peer[i], err = rpc.Dial("tcp", addr_node)
		utils.Check_error(err)
	}

	// Clear the shell
	clear_shell()

	for {
		fmt.Println("Select the test to perform:")
		fmt.Println("1. Only sender message multicast")
		fmt.Println("2. More sender message multicast")
		fmt.Println("3. Exit")

		fmt.Scanf("%d\n", &choice)

		switch choice {
		case 1:
			clear_shell()
			outcome := false

			switch algo {
			case 1:
				outcome = test_one_sender_1(peer_number)
				break
			case 2:
				outcome = test_one_sender_2(peer_number)
				break
			case 3:
				break
			}

			if outcome {
				fmt.Println("Test passed")
			} else {
				fmt.Println("Test NO passed")
			}
			break
		case 2:
			clear_shell()
			outcome := false

			switch algo {
			case 1:
				outcome = test_more_sender_1(peer_number)
				break
			case 2:

				break
			case 3:
				break
			}

			if outcome {
				fmt.Println("Test passed")
			} else {
				fmt.Println("Test NO passed")
			}
			break
		case 3:
			return
		default:
			fmt.Println("Codice operativo non supportato.")
			break
		}

	}
}
