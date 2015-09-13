package main

import (
    "fmt"
    "net"
    "time"
    "os"
    "strconv"
    "strings"
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
	private chan string
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
	server.private = make(chan string, 100)
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
        	received := string(buf[:n])
        	hdr := received[:strings.Index(received, " ")]
			msg := received[strings.Index(received, " ")+1:]
	        var message string
	        switch {
	        	case hdr == "REGISTER" :
	        		if _, ok := server.userList[addr.IP.String()]; !ok {
		        		server.userList[addr.IP.String()] = msg
			        	message = "User " + msg + " (" + addr.IP.String() + ") joined at " + getCurrTime()
			        	fmt.Println(message)
			        	server.messages <- message
			        }
	        	case hdr == "MESSAGE" :
	        		fmt.Println("Received", msg, "from", server.userList[addr.IP.String()], "(", addr.IP.String() ,") at", getCurrTime())
	        		message = server.userList[addr.IP.String()] + " said at " + getCurrTime() + ": " + msg
	        		server.messages <- message
	        	case hdr == "LIST" :
	        		fmt.Println("List request from", server.userList[addr.IP.String()], "at", getCurrTime())
	        		for _, login := range server.userList {
	        			message += login + "\n"
	        		}
	        		server.messages <- message
	        	case hdr == "PRIVATE" :
	        		rcvr := msg[:strings.Index(msg, " ")]
					msg := msg[strings.Index(msg, " ")+1:]
	        		fmt.Println("Received private message to", rcvr, msg, "from", server.userList[addr.IP.String()], "(", addr.IP.String() ,") at", getCurrTime())
					message = rcvr + " " + server.userList[addr.IP.String()] + " privately said at " + getCurrTime() + ": " + msg
	        		server.private <- message
	        	case hdr == "LEAVE" :
	        		if _, ok := server.userList[addr.IP.String()]; ok {
			        	message = "User " + server.userList[addr.IP.String()] + " (" + addr.IP.String() + ") left at " + getCurrTime()
			        	fmt.Println(message)
			        	delete(server.userList, addr.IP.String())
			        	server.messages <- message
			        }
	        }
	        
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

func (server *Server) PrivateBroadcast() {
	for {
		message := <- server.private
		func(){
			rcvr := message[:strings.Index(message, " ")]
			msg := message[strings.Index(message, " ")+1:]
			RecvAddr,err := net.ResolveUDPAddr("udp", rcvr+":"+strconv.Itoa(server.broadcastPort))
			CheckError(err)
			LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
			CheckError(err)
			Conn, err := net.DialUDP("udp", LocalAddr, RecvAddr)
			CheckError(err)
			defer Conn.Close()
			buf := []byte(msg)
			_, err = Conn.Write(buf)
			CheckError(err);
		}()
	}
}