/*
	This is a frontend that allow a generic peer to attach to a specific container and send messages.
*/

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Welcome to multicast group!")
	fmt.Printf("Select one of available container:\n")

	var container int

	fmt.Scanf("%d", container)

	if container == 1 {
		fmt.Println("Forza Roma!")
	}
}
