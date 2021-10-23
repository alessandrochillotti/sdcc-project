/*
	This is a file that allow to perform handshake phase in according to received frontend informations.
*/

package main

import (
	"net"
	"net/rpc"

	"alessandro.it/app/utils"
)

// This struct is used to build a new peer in according to received frontend informations
type Handshake struct {
	New_peer *Peer
}

// This function allow to manage the log file to frontend
func (h *Handshake) Handshake(request *utils.Hand_request, reply *utils.Hand_reply) error {
	h.New_peer = &Peer{Index: 0, Ip_address: getIpAddress(), Port: PORT, Username: request.Username}
	// p.username = request.Username
	reply.Ip_address = h.New_peer.Ip_address
	reply.Algorithm = conf.Algorithm
	conf.Test = request.Test

	channel_handshake <- true

	return nil
}

// This function make an handshake with frontend
func (h *Handshake) frontend_handshake() *net.Listener {
	// Register a new RPC server and the struct we created above
	frontend_handshake := rpc.NewServer()
	err := frontend_handshake.RegisterName("Handshake", h)
	utils.Check_error(err)

	// Listen for incoming messages on port 4444
	listener_frontend, err := net.Listen("tcp", ":4444")
	utils.Check_error(err)

	go frontend_handshake.Accept(listener_frontend)

	return &listener_frontend
}
