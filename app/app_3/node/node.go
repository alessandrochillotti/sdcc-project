/*
	This is a generic node that must register in group multicast and, then, it can communicate
	with other nodes of newtork.
*/

package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"alessandro.it/app/lib"
	"alessandro.it/app/utils"
	"github.com/go-redis/redis"
)

type Node int

var channel_connection chan bool

// Constant values
const MAX_QUEUE = 100
const MAX_DELAY = 3

// Global variables
var vector_clock *utils.Vector_clock
var addresses [lib.NUMBER_NODES]string /* Contains ip addresses of each node in multicast group */
var peer [lib.NUMBER_NODES]*rpc.Client
var queue *utils.Queue
var my_index int
var client *redis.Client

var mutex_queue sync.Mutex
var mutex_clock sync.Mutex

/*
This function return the ip address of current node
*/
func getIpAddress() string {
	ip_address := "0.0.0.0"
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip_address = ipv4.String()
		}
	}

	return ip_address
}

/*
This function build a struct that contains the info to register the node
*/
func build_whoami_struct(whoami_to_register *lib.Whoami) {
	whoami_to_register.Ip_address = getIpAddress()
	whoami_to_register.Port = "1234"
}

/*
This function allows to register the node to communicate in multicast group
*/
func register_into_group() {
	var whoami_to_register lib.Whoami
	var empty lib.Empty

	build_whoami_struct(&whoami_to_register)

	// The RPC server has ip address set to 10.5.0.254 and it is listening in port 4321
	addr_register_node := "10.5.0.254:4321"

	// Try to connect to addr_register_node
	client, err := rpc.Dial("tcp", addr_register_node)
	lib.Check_error(err)

	// Call remote procedure and reply will store the RPC result
	err = client.Call("Register.Register_node", &whoami_to_register, &empty)
	lib.Check_error(err)

	client.Close()
}

/*
This function log message into file.
*/
func log_message(pkt *lib.Packet, id int) {
	// Open file into volume docker
	path_file := "/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt"
	f, err := os.OpenFile(path_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	lib.Check_error(err)
	defer f.Close()

	// Write into file the ip address of registered node
	_, err = f.WriteString(pkt.Source_address + " -> " + pkt.Message + "[" + strconv.Itoa(id) + "]\n")
	lib.Check_error(err)
}

/*
This function check if there are packet to deliver, so the following conditions must be checked:
	1. The message inviato from p_j to current process is the expected message from p_j.
	2. The current process has seen every messahe that p_j has seen.
*/
func deliver_packet() {
	current_index := 1
	for {
		mutex_queue.Lock()
		node_to_deliver := queue.Get_node(current_index)
		mutex_queue.Unlock()
		deliver := true
		index_pid_to_deliver := 0
		if node_to_deliver == nil {
			deliver = false
		} else {
			index_pid_to_deliver = node_to_deliver.Update.Packet.Index_pid
			current_index = current_index + 1
		}

		if deliver && node_to_deliver.Update.Timestamp.Clocks[index_pid_to_deliver] == vector_clock.Clocks[index_pid_to_deliver]+1 {
			for k := 0; k < lib.NUMBER_NODES && deliver; k++ {
				if k != index_pid_to_deliver && node_to_deliver.Update.Timestamp.Clocks[k] > vector_clock.Clocks[k] {
					deliver = false
				}
			}
		}

		if deliver {
			vector_clock.Update_with_max(node_to_deliver.Update.Timestamp)
			log_message(&node_to_deliver.Update.Packet, node_to_deliver.Update.Packet.Id)
			queue.Remove_node(node_to_deliver.Update.Packet.Id)

			// Clear shell
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			// Print chat
			content, err := ioutil.ReadFile("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
			lib.Check_error(err)

			// fmt.Println("Consegnato messaggio con timestamp", head_node.Timestamp)
			list := string(content)

			print(list)
		}
	}
}

/*
This function allow to increment ID packet into Redis container in transactional mode.
*/
func increment_id(id *int) {
	const routineCount = 100

	// Transactionally increments key using GET and SET commands.
	increment := func(key string) error {
		txf := func(tx *redis.Tx) error {
			// get current value or zero
			n, err := tx.Get("ID").Int()
			if err != nil && err != redis.Nil {
				return err
			}

			// actual opperation (local in optimistic lock)
			n++

			// runs only if the watched keys remain unchanged
			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
				// pipe handles the error case
				pipe.Set(key, n, 0)
				*id = n
				return nil
			})
			return err
		}

		for retries := routineCount; retries > 0; retries-- {
			err := client.Watch(txf, key)
			if err != redis.TxFailedErr {
				return err
			}
			// optimistic lock lost
		}
		return errors.New("increment reached maximum number of retries")
	}

	increment("ID")
}

/*
This function allow to wait the input of user to send the message to each node of group multicast
*/
func open_standard_input() {
	var pkt_id int
	var ack utils.Ack = 0
	for {
		// Take in input the content of message to send
		in := bufio.NewReader(os.Stdin)
		text, _ := in.ReadString('\n')
		text = strings.TrimSpace(text)

		// Increment ID in transactional mode to Redis container
		increment_id(&pkt_id)
		// Build packet
		pkt := lib.Packet{Id: pkt_id, Source_address: getIpAddress(), Index_pid: my_index, Message: text}

		// Update the scalar clock and build update packet to send
		mutex_clock.Lock()
		vector_clock.Increment(my_index)
		update := utils.Update{Timestamp: *vector_clock, Packet: pkt}
		mutex_clock.Unlock()

		// Send to each node of group multicast the message
		for i := 0; i < lib.NUMBER_NODES; i++ {
			lib.Delay(3)
			/*
				The following 3 lines allow to test the algorithm 3 in case of scenario that we saw in class.
			*/
			// if pkt_id == 1 && i == 2 {
			// 	time.Sleep(time.Duration(10) * time.Second)
			// }
			err := peer[i].Call("Node.Get_update", &update, &ack)
			lib.Check_error(err)
		}
	}
}

/*
This function, after reception of list from register node, allow to setup connection with each node of group multicast.
*/
func setup_connection() error {
	var err error

	for i := 0; i < lib.NUMBER_NODES; i++ {
		addr_node := addresses[i] + ":1234"
		peer[i], err = rpc.Dial("tcp", addr_node)
		lib.Check_error(err)
		if addresses[i] == getIpAddress() {
			my_index = i
		}
	}

	return nil
}

/* RPC methods registered by Node */

/*
This RPC method of Node allow to get list from the registered node.
*/
func (node *Node) Get_list(list *lib.List_of_nodes, reply *lib.Empty) error {
	// Parse the list and put the addresses into destination array
	addr_tmp := strings.Split(list.List_str, "\n")
	for i := 0; i < lib.NUMBER_NODES; i++ {
		addresses[i] = addr_tmp[i]
	}

	channel_connection <- true

	return nil
}

/*
This RPC method of Node allow to get update from the other node of group multicast
*/
func (node *Node) Get_update(update *utils.Update, ack *utils.Ack) error {
	if update.Packet.Source_address != getIpAddress() {
		mutex_clock.Lock()
		// vector_clock.Update_with_max(update.Timestamp)
		// vector_clock.Increment(my_index)
		mutex_clock.Unlock()
	}

	// Build update node to insert the packet into queue
	update_node := &utils.Node{Update: *update, Next: nil, Ack: 0}

	// Insert update node into queue
	mutex_queue.Lock()
	queue.Append(update_node)
	// queue.Display()
	mutex_queue.Unlock()

	vector_clock.Print()

	return nil
}

func main() {
	// For first thing, the node communicates with the register node to register his info
	register_into_group()

	// Allocate object to use it into program execution
	node := new(Node)
	queue = &utils.Queue{}
	channel_connection = make(chan bool)

	// Create file for log of messages
	f, err := os.Create("/home/alessandro/Dropbox/Università/SDCC/sdcc-project/mnt/" + getIpAddress() + "_log.txt")
	lib.Check_error(err)
	defer f.Close()

	// Register the Node methods
	receiver := rpc.NewServer()
	err = receiver.RegisterName("Node", node)
	lib.Check_error(err)

	// Listen for incoming messages on port 1234
	lis, err := net.Listen("tcp", ":1234")
	lib.Check_error(err)
	defer lis.Close()

	client = redis.NewClient(&redis.Options{
		Addr:     "10.5.0.250:6379",
		Password: "password",
		DB:       0,
	})

	err = client.Set("ID", 0, 0).Err()
	lib.Check_error(err)

	// Initialize vector clock
	vector_clock = new(utils.Vector_clock)
	vector_clock.Init()

	// Use goroutine to implement a lightweight thread to manage the coming of new messages
	go receiver.Accept(lis)

	// Setup the connection with the peer of group multicast after the reception of list
	<-channel_connection
	if setup_connection() != nil {
		os.Exit(1)
	}

	// This goroutine check always if there are packet to deliver
	go deliver_packet()

	// The user can insert text to send to each node of group multicast
	open_standard_input()
}
