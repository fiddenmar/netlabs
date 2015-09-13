package main

import (
	client "../"
	"fmt"
	"time"
)

func main() {
	answers := make(chan string)
	var c client.Client
	go func() {
		c.Init("User", "127.0.0.1", 34310, answers)
		go c.Answer()
		c.Register()
		c.Message("Hello world!1")
		c.Message("Hello world!2")
		c.List()
		c.Leave()
	}()
	for {
		fmt.Println(<-answers)
	}
}