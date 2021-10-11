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
				fmt.Println(port)
			}
		}
	}

	return strconv.Itoa(int(port))
}

func get_list_containers() {
	cnt := 0

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	// TODO: control if container is already selected
	for _, container := range containers {
		fmt.Println(container.Names[0])
		if container.Names[0] == "app_node_1" {
			cnt++
			fmt.Printf("%d. %s\n", cnt, (container.Names[0])[1:])
		}
	}
}

func main() {
	var text string
	var choice int
	var empty lib.Empty
	var ip_container string

	// Print menù
	fmt.Println("Welcome to App")
	fmt.Println("Inserisci il numero di container che vuoi utilizzare (1, 2, ..., 3)")
	fmt.Scanf("%d", &selected_container)

	addr_node := "127.0.0.1:" + get_free_port(selected_container)
	client, err := rpc.Dial("tcp", addr_node)
	lib.Check_error(err)

	client.Call("Node.Handshake", &empty, &ip_container)
	check_error(err)

	path_file := "./volumes/log_node/" + ip_container + "_log.txt"

	for {
		fmt.Println("Quale operazione vuoi effettuare:")
		fmt.Println("1. Invio messaggio")
		fmt.Println("2. Stampa chat")
		fmt.Println("3. Uscire")

		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			in := bufio.NewReader(os.Stdin)
			text, err = in.ReadString('\n')
			text = strings.TrimSpace(text)

			err = client.Call("Node.Get_message_from_frontend", &text, &empty)
			check_error(err)

			// Clear the shell
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()

			fmt.Println("Messaggio inviato correttamente.")
			fmt.Println()
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
			fmt.Println("Arrivederci")
			return
		default:
			fmt.Println("Codice operativo non supportato.")
			break
		}

	}
}
