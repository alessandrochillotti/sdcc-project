/*
	This is a frontend that allow a generic peer to attach to a specific container and send messages.
*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func checkError(err error) {
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
			if container.Ports[i].PrivatePort == 4321 {
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
		if container.Names[0][1:9] == "app_node" {
			cnt++
			fmt.Printf("%d. %s\n", cnt, (container.Names[0])[1:])
		}
	}
}

func main() {

	var container int

	fmt.Println("Welcome to multicast group!")
	fmt.Printf("Insert the number of container that you want manage:\n")
	get_list_containers()

	fmt.Scanf("%d", &container)

	service := "127.0.0.1:" + get_free_port(container)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	var text string

	in := bufio.NewReader(os.Stdin)
	text, err = in.ReadString('\n')
	text = strings.TrimSpace(text)

	_, err = conn.Write([]byte(text))
	checkError(err)

	os.Exit(0)
}
