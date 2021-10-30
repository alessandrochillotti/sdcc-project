/*
This file mantain the info to communicate into multicast group.
*/
package utils

import "net/rpc"

// This struct is used to mantain the network information to communicate with other peer
type Connection struct {
	Addresses []string
	Username  map[string]string
	Peer      []*rpc.Client
}

func (c *Connection) Init_connection(nodes int) {
	c.Addresses = make([]string, nodes)
	c.Username = make(map[string]string)
	c.Peer = make([]*rpc.Client, nodes)
}

func (c *Connection) GetUsername(ip string) string {
	return c.Username[ip]
}
