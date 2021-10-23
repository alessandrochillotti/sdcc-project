package main

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"strconv"
	"strings"
	"time"

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
		file, err := ioutil.ReadFile(path_file)
		check_error(err)

		if len(string(file)) != 0 {
			return false
		}
	}

	return true
}

func check_log_2(number int) bool {
	equal := true
	minor_len := 10000

	for i := 0; i < number; i++ {
		path_file := "./volumes/log_node/10.5.0." + strconv.Itoa(i+2) + "_log.txt"
		file, err := ioutil.ReadFile(path_file)
		check_error(err)
		lines := len(strings.Split(string(file), "\n")) - 1

		if lines < int(minor_len) {
			minor_len = lines
		}
	}

	path_file := "./volumes/log_node/10.5.0.2_log.txt"
	file_1, err := ioutil.ReadFile(path_file)
	lines_1 := strings.Split(string(file_1), "\n")
	check_error(err)
	for i := 1; i < number && equal; i++ {
		path_file := "./volumes/log_node/10.5.0." + strconv.Itoa(i+2) + "_log.txt"
		file_2, err := ioutil.ReadFile(path_file)
		check_error(err)
		lines_2 := strings.Split(string(file_2), "\n")

		for j := 0; j < minor_len; j++ {
			if lines_1[j] != lines_2[j] {
				return false
			}
		}

		lines_1 = lines_2
	}

	return true
}

// OK
func test_one_sender_1(number int) bool {
	var empty *utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}

	for i := 0; i < 6; i++ {
		err := peer[0].Call("Peer.Get_message_from_frontend", &messages[i], &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	return check_log_equal(number)
}

// OK
func test_more_sender_1(number int) bool {
	var empty *utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}

	for i := 0; i < 6; i++ {
		err := peer[i%number].Call("Peer.Get_message_from_frontend", &messages[i], &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	return check_log_equal(number)
}

// OK
func test_one_sender_2(number int) bool {
	var empty *utils.Empty

	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}

	for i := 0; i < 6; i++ {
		err := peer[0].Call("Peer.Get_message_from_frontend", &messages[i], &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	return check_log_empty(number)
}

func test_more_sender_2(number int) bool {
	var empty string
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}

	for i := 0; i < 6; i++ {
		err := peer[i%number].Call("Peer.Get_message_from_frontend", &messages[i], &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(10) * time.Second)

	return check_log_2(number)
}

// NO OK
func test_one_sender_3(number int) bool {
	return test_one_sender_1(number)
}

func test_more_sender_3(number int) bool {
	var empty utils.Empty
	message_1 := "causa"
	message_2 := "effetto"

	err := peer[0].Call("Peer.Get_message_from_frontend", &message_1, &empty)
	check_error(err)
	err = peer[1].Call("Peer.Get_message_from_frontend", &message_2, &empty)
	check_error(err)

	time.Sleep(time.Duration(10) * time.Second)

	return check_log_equal(number)
}

func main() {
	var algo int
	var choice int

	// Print menù
	peer_number := len(strings.Split(get_list_container(), "\n")) - 1
	peer = make([]*rpc.Client, peer_number)

	// Dial of peer
	for i := 0; i < peer_number; i++ {
		var err error

		_, _, algo = handshake(i+1, true)
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
				outcome = test_one_sender_3(peer_number)
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
				outcome = test_more_sender_2(peer_number)
				break
			case 3:
				outcome = test_more_sender_3(peer_number)
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
