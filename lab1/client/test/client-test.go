package main

import (
	client "../"
)

func main() {
	var c client.Client
	c.Init("User", "127.0.0.1", 34310)
	go c.Answer()
	c.Connect()
	c.Send("Hello message")
}