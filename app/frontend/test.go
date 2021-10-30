/*
This file is an application that allow to test the algorithms.
*/
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

const WAIT int = 3

var peer []*rpc.Client

/*
Test check for verification
*/

// This function check if the files log have the same content.
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

// This function check if the files log are empty.
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

// This function check if the first N lines of file log are equal, where N is the minimum number line of files log.
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

/*
Test scenarios
*/

/*
Algorithm: 1

Scenario: Only one peer send message in multicast.

Expected outcome: The log files have the same content.
*/
func test_one_sender_1(number int) bool {
	var empty utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}
	delay := make([]int, number)

	for i := 0; i < 6; i++ {
		msg := &utils.Message{Text: messages[i], Delay: delay}
		err := peer[0].Call("Peer.Get_message_from_frontend", msg, &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(WAIT) * time.Second)

	return check_log_equal(number)
}

/*
Algorithm: 1

Scenario: More than one peer send message in multicast.

Expected outcome: The log files have the same content.
*/
func test_more_sender_1(number int) bool {
	var empty utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}
	delay := make([]int, number)

	divCall := make([]*rpc.Call, number)

	for i := 0; i < number; i++ {
		msg := &utils.Message{Text: messages[i], Delay: delay}
		divCall[i] = peer[i].Go("Peer.Get_message_from_frontend", msg, &empty, nil)

		go wait_to_send(i, divCall[i], messages[i+3], delay, &empty)
	}

	time.Sleep(time.Duration(WAIT) * time.Second)

	return check_log_equal(number)
}

/*
Algorithm: 2

Scenario: Only one peer send message in multicast.

Expected outcome: The log files are empty.
*/
func test_one_sender_2(number int) bool {
	var empty utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}
	delay := make([]int, number)

	for i := 0; i < 6; i++ {
		msg := &utils.Message{Text: messages[i], Delay: delay}
		err := peer[0].Call("Peer.Get_message_from_frontend", msg, &empty)
		check_error(err)
	}

	time.Sleep(time.Duration(WAIT) * time.Second)

	return check_log_empty(number)
}

/*
Algorithm: 2

Scenario: More than one peer send message in multicast.

Expected outcome: The first N lines of file log are equal, where N is the minimum number line of files log.
*/
func test_more_sender_2(number int) bool {
	var empty utils.Empty
	var messages = [6]string{"uno", "due", "tre", "quattro", "cinque", "sei"}
	delay := make([]int, number)

	divCall := make([]*rpc.Call, number)

	for i := 0; i < number; i++ {
		msg := &utils.Message{Text: messages[i], Delay: delay}
		divCall[i] = peer[i].Go("Peer.Get_message_from_frontend", msg, &empty, nil)

		go wait_to_send(i, divCall[i], messages[i+3], delay, &empty)
	}

	time.Sleep(time.Duration(WAIT) * time.Second)

	return check_log_2(number)
}

/*
Algorithm: 3

Scenario: Only one peer send message in multicast.

Expected outcome: The log files have the same content.
*/
func test_one_sender_3(number int) bool {
	return test_one_sender_1(number)
}

/*
Algorithm: 3

Scenario: Example in class.

Expected outcome: Respect causality.
*/
func test_more_sender_3(number int) bool {
	var empty utils.Empty
	delay := make([]int, number)

	delay[2] = 3
	msg_1 := &utils.Message{Text: "causa", Delay: delay}
	peer[0].Go("Peer.Get_message_from_frontend", msg_1, &empty, nil)

	time.Sleep(time.Duration(1) * time.Second)

	delay[2] = 0
	msg_2 := &utils.Message{Text: "effetto", Delay: delay}
	err := peer[1].Call("Peer.Get_message_from_frontend", msg_2, &empty)
	check_error(err)

	time.Sleep(time.Duration(WAIT) * time.Second)

	return check_log_equal(number)
}

/*
Utility
*/
func wait_to_send(index int, divCall *rpc.Call, message string, delay []int, empty *utils.Empty) {
	msg := &utils.Message{Text: message, Delay: delay}
	<-divCall.Done
	peer[index].Go("Peer.Get_message_from_frontend", msg, empty, nil)
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
				fmt.Println("Test passed\n")
			} else {
				fmt.Println("Test NO passed\n")
			}
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
