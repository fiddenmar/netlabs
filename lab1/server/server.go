package main

import (
    "fmt"
    "net"
    "time"
    "os"
    "strconv"
)

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func getCurrTime() string {
	return time.Now().Format(time.RFC850)
}

type Server struct {
	listenPort int
	broadcastPort int
	userList map[string]string //ip->username
	messages chan string
}

func main() {   
	var server Server
	server.Init(34310)
	go server.Broadcast()
	server.Listen()
	for {

	}
}

func (server *Server) Init(_listenPort int) {
	server.listenPort = _listenPort
	server.broadcastPort = server.listenPort + 1
	server.userList = nil
	server.userList = make(map[string]string)
	server.messages = make(chan string, 100)
}

func (server *Server) Listen() {
	ServerAddr,err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(server.listenPort))
    CheckError(err)
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
    buf := make([]byte, 1024)

    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        CheckError(err)
        if n > 0 {
	        var message string
	        if _, ok := server.userList[addr.IP.String()]; ok {
	        	fmt.Println("Received", string(buf[0:n]), "from", server.userList[addr.IP.String()], "(", addr.IP.String() ,") at", getCurrTime())
	        	message = server.userList[addr.IP.String()] + " said at " + getCurrTime() + ": " + string(buf[0:n])
	        } else {
	        	server.userList[addr.IP.String()] = string(buf[0:n])
	        	message = "User " + string(buf[0:n]) + " (" + addr.IP.String() + ") joined at " + getCurrTime()
	        	fmt.Println(message)
	        }
	        server.messages <- message
	    }
    }
}

func (server *Server) Broadcast() {
	for {
		msg := <- server.messages
		for ip, _ := range server.userList {
			RecvAddr,err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(server.broadcastPort))
			CheckError(err)
			LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
			CheckError(err)
			Conn, err := net.DialUDP("udp", LocalAddr, RecvAddr)
			CheckError(err)
			defer Conn.Close()
			buf := []byte(msg)
			_, err = Conn.Write(buf)
			CheckError(err);
		}
	}
}