package client

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"runtime"
)

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}

type MessageType int
const (
	REGISTER MessageType = iota
	MESSAGE
	LIST
	PRIVATE
	LEAVE
)
var messageTypes = [...]string {
	"REGISTER",
	"MESSAGE",
	"LIST",
	"PRIVATE",
	"LEAVE",
}
type MessageHeader struct {
	Type MessageType
	Extra string
}
func (m *MessageHeader) ToString() (message string) {
	return messageTypes[m.Type] + " " + m.Extra
}

type Client struct {
	sendIP string
	sendPort int
	answerPort int
	login string
	pending chan bool
	Answers chan string
}

func (client *Client) Init(_login string, _sendIP string, _sendPort int, _answers chan string) {
	client.sendPort = _sendPort
	client.answerPort = client.sendPort + 1
	client.sendIP = _sendIP
	client.login = _login
	client.pending = make(chan bool)
	client.Answers = _answers
}

func (client *Client) Register() {
	var header MessageHeader
	header.Type = REGISTER
	client.send(header, client.login)
}

func (client *Client) Message(message string) {
	var header MessageHeader
	header.Type = MESSAGE
	client.send(header, message)
}

func (client *Client) List() {
	var header MessageHeader
	header.Type = LIST
	client.send(header, "LIST")
}

func (client *Client) Private(receiver string, message string) {
	var header MessageHeader
	header.Type = PRIVATE
	header.Extra = receiver
	client.send(header, " "+message)
}

func (client *Client) Leave() {
	var header MessageHeader
	header.Type = LEAVE
	client.send(header, "LEAVE")
}

func (client *Client) send(header MessageHeader, message string) {
	ServerAddr,err := net.ResolveUDPAddr("udp", client.sendIP+":"+strconv.Itoa(client.sendPort))
	CheckError(err)
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
	defer Conn.Close()
	buf := []byte(header.ToString() + message)
	_, err = Conn.Write(buf)
	CheckError(err)
	client.pending <- true
}

func (client *Client) Answer() {
	runtime.LockOSThread()
	AnswerAddr,err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(client.answerPort))
    CheckError(err)
    AnswerConn, err := net.ListenUDP("udp", AnswerAddr)
    CheckError(err)
    defer AnswerConn.Close()
    buf := make([]byte, 1024)

    for {
    	<- client.pending
    	var ans string
    	AnswerConn.SetDeadline(time.Now().Add(3*time.Second))
	    n, _, err := AnswerConn.ReadFromUDP(buf)
	    switch e := err.(type) {
		case net.Error:
		    if e.Timeout() {
		        ans = "Timeout error: check your internet connection"
		    }
		default:
			CheckError(err)
		    ans = string(buf[0:n])
		}
	    client.Answers <- ans
	}
}