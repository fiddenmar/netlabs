package main

import (
	client "../"
	"fmt"
)

func main() {
	var c client.Client
	go func() {
		c.Init("User", "127.0.0.1", 34310)
		go c.Answer()
		c.Connect()
		c.Send("Hello message")
	}()
	
	fmt.Println(<-c.Answers)
}