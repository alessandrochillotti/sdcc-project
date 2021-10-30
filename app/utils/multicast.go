/*
This file mantain the info to communicate into multicast group.
*/
package utils

import "net/rpc"

// This struct is used to mantain the network information to communicate with other peer
type Connection struct {
	Addresses []string
	Username  map[string]string
	Index     map[string]int
	Peer      []*rpc.Client
}

func (c *Connection) Init_connection(nodes int) {
	c.Addresses = make([]string, nodes)
	c.Username = make(map[string]string)
	c.Index = make(map[string]int)
	c.Peer = make([]*rpc.Client, nodes)
}

func (c *Connection) GetUsername(ip string) string {
	return c.Username[ip]
}

func (c *Connection) GetIndex(ip string) int {
	return c.Index[ip]
}
