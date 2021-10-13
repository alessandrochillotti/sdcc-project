/*
This file mantain the info to communicate into multicast group.
*/
package utils

import "net/rpc"

// This struct is used to mantain the network information to communicate with other peer
type Connection struct {
	Addresses          []string
	Peer               []*rpc.Client
	Channel_connection chan bool
}

func (c *Connection) Init_connection(nodes int) {
	c.Addresses = make([]string, nodes)
	c.Peer = make([]*rpc.Client, nodes)
	c.Channel_connection = make(chan bool)
}
