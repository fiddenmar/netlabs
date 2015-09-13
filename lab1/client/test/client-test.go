package main

import (
	client "../"
	"fmt"
	"time"
	"os"
)

func main() {
	answers := make(chan string, 100)
	var c client.Client
	go func() {
		c.Init("Adam", "127.0.0.1", 34310, answers)
		go c.Answer()
		time.Sleep(1*time.Second)
		c.Register()
		time.Sleep(1*time.Second)
		c.Message("Hello world!1")
		time.Sleep(1*time.Second)
		c.Message("Hello world!2")
		time.Sleep(1*time.Second)
		c.Private("Adam", "Private hello!")
		time.Sleep(1*time.Second)
		c.List()
		time.Sleep(1*time.Second)
		c.Message("Hello world!3")
		time.Sleep(1*time.Second)
		c.Message("Hello world!4")
		time.Sleep(1*time.Second)
		c.Leave()
		time.Sleep(2*time.Second)
		os.Exit(0)
	}()
	for {
		fmt.Println(<-answers)
	}
}