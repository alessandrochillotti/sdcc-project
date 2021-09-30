package utils

type Update struct {
	Source_address string
	Source_pid     int
	Message        string
}

type Ack struct {
	Source_address string
	Source_pid     int
	Message        string
}
