package client

import (
	"fmt"
	"net"
	"strconv"
)

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}

type Client struct {
	sendIP string
	sendPort int
	answerPort int
	login string
	Answers chan string
}

func (client *Client) Init(_login string, _sendIP string, _sendPort int) {
	client.sendPort = _sendPort
	client.answerPort = client.sendPort + 1
	client.sendIP = _sendIP
	client.login = _login
	client.Answers = make(chan string)
}

func (client *Client) Connect() {
	client.Send(client.login)
}

func (client *Client) Send(message string) {
	ServerAddr,err := net.ResolveUDPAddr("udp", client.sendIP+":"+strconv.Itoa(client.sendPort))
	CheckError(err)
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
	defer Conn.Close()
	buf := []byte(message)
	_, err = Conn.Write(buf)
	CheckError(err)
}

func (client *Client) Answer() (message string) {
	AnswerAddr,err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(client.answerPort))
    CheckError(err)
    AnswerConn, err := net.ListenUDP("udp", AnswerAddr)
    CheckError(err)
    defer AnswerConn.Close()
    buf := make([]byte, 1024)

    for {
	    n, _, err := AnswerConn.ReadFromUDP(buf)
	    CheckError(err)
	    ans := string(buf[0:n])
	    client.Answers <- ans
	}
}