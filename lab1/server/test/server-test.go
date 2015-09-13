package main

import (
	server "../"
)

func main() {   
	var s server.Server
	s.Init(34310)
	go s.Broadcast()
	go s.PrivateBroadcast()
	s.Listen()
}