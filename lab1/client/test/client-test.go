package main

import (
	client "../"
	"fmt"
	"time"
)

func main() {
	var c client.Client
	go func() {
		c.Init("User", "127.0.0.1", 34310)
		go c.Answer()
		c.Register()
		c.Message("Hello world!1")
		c.Register()
		c.Message("Hello world!2")
		c.List()
		time.Sleep(2*time.Second)
		c.Leave()
	}()
	for {
		fmt.Println(<-c.Answers)
	}
}