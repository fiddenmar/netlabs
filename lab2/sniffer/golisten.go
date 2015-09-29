package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	args := os.Args[1:]
	messages := make(chan string)
	for _, arg := range args {
		switch {
		case arg == "ip":
			go listenIP(messages)
		case arg == "icmp":
			go listenICMP(messages, false)
		case arg == "udp":
			go listenUDP(messages, false)
		case arg == "tcp":
			go listenTCP(messages, false)
		}
	}
	for {
		msg:=<-messages
		fmt.Println(msg)
	}
}

func intToIP(a int) (string) {
	return fmt.Sprintf("%d.%d.%d.%d", byte(a>>24), byte(a>>16), byte(a>>8), byte(a))
}

func formatIP(in []byte) (out string) {
	version := int(in[0]>>4)
	ihl := int(in[0]&0x0f)
	dscp := int(in[1]>>2)
	ecn := int(in[1]&0x03)
	length := (int(in[2])<<8 + int(in[3]))
	identification := (int(in[4]))<<8 + int(in[5])
	flags := int(in[6]>>5)
	offset := (int(in[6]&0x1f) + int(in[7]))
	ttl := int(in[8])
	protocol := int(in[9])
	checksum := int(int(in[10])<<8 + int(in[11]))
	source := int(int(in[12])<<24 + int(in[13])<<16 + int(in[14])<<8 + int(in[15]))
	destination := int(int(in[16])<<24 + int(in[17])<<16 + int(in[17])<<8 + int(in[19]))
	data_hex := in[ihl:]
	data_str := string(data_hex)
	out = fmt.Sprintf("IP:\n version: % d\n ihl: % d\n", version, ihl) +
			fmt.Sprintf("dscp: % d\n ecn: % d\n length: % d\n", dscp, ecn, length) +
			fmt.Sprintf("identification: % d\n flags: % b\n", identification, flags) + 
			fmt.Sprintf("offset: % d\n ttl: % d\n protocol: % d\n", offset, ttl, protocol) +		
			fmt.Sprintf("checksum: % X\n source: % s\n destination: % s\n", checksum, intToIP(source), intToIP(destination)) + 	
			fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)	
	return
}

func formatICMP(in []byte) (out string) {
	ipHeaderSize := int(in[0]&0x0f)*4
	in = in[ipHeaderSize:]
	icmpType := int(in[0])
	code := int(in[1])
	checksum := int(in[2])<<8+int(in[3])
	rest_hex := in[4:8]
	rest_str := string(rest_hex)
	data_hex := in[8:]
	data_str := string(data_hex)
	out = fmt.Sprintf("ICMP:\n type: % d\n code: % d\n", icmpType, code) +
			fmt.Sprintf("checksum: % X\n rest (HEX): % X\n rest (STR): % s\n", checksum, rest_hex, rest_str) + 	
			fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)	
	return
}

func formatUDP(in []byte) (out string) {
	ipHeaderSize := int(in[0]&0x0f)*4
	in = in[ipHeaderSize:]
	var protocol int
	if (len(in)>9) {
		protocol = int(in[9])
	}
	if protocol != 17 {
		sourcePort := int(int(in[0])<<8 + int(in[1]))
		destPort := int(int(in[2])<<8 + int(in[3]))
		length := int(int(in[4])<<8 + int(in[5]))
		checksum := int(int(in[6])<<8 + int(in[7]))
		data_hex := in[8:]
		data_str := string(data_hex)
		out = fmt.Sprintf("UDP:\n source port: % d\n destination port: % d\n", sourcePort, destPort) + 
				fmt.Sprintf("length: % d\n checksum: % X\n", length, checksum) +			
				fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)
	} else {
		sourceIP := int(int(in[0])<<24 + int(in[1])<<16 + int(in[2])<<8 + int(in[3]))
		destIP := int(int(in[4])<<24 + int(in[5])<<16 + int(in[6])<<8 + int(in[7]))
		zeroes := int(in[8])
		udpLength := int(int(in[10])<<8 + int(in[11]))
		sourcePort := int(int(in[12])<<8 + int(in[13]))
		destPort := int(int(in[14])<<8 + int(in[15]))
		length := int(int(in[16])<<8 + int(in[17]))
		checksum := int(int(in[18])<<8 + int(in[19]))
		data_hex := in[20:]
		data_str := string(data_hex)
		out = fmt.Sprintf("UDP IPv4:\n sourceIP: % s\n destIP: % s\n", intToIP(sourceIP), intToIP(destIP)) +
				fmt.Sprintf("zeroes: % d\n protocol: % d\n udp length: % d\n", zeroes, protocol, udpLength) +
				fmt.Sprintf("source port: % d\n destination port: % d\n", sourcePort, destPort) + 
				fmt.Sprintf("length: % d\n checksum: % X\n", length, checksum) +			
				fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)	
	}
	return
}

func formatTCP(in []byte) (out string) {
	ipHeaderSize := int(in[0]&0x0f)*4
	in = in[ipHeaderSize:]
	protocol := int(in[9])
	if protocol != 6 {
		sourcePort := int(int(in[0])<<8 + int(in[1]))
		destPort := int(int(in[2])<<8 + int(in[3]))
		seq := int(int(in[4])<<24 + int(in[5])<<16 + int(in[6])<<8 + int(in[7]))
		ack := int(int(in[8])<<24 + int(in[9])<<16 + int(in[10])<<8 + int(in[11]))
		offset := int(in[12]>>4)
		reserved := int(in[12]&0x0f)
		flags := int(in[13])
		windowSize := int(in[14])<<8 + int(in[15])
		checksum := int(in[16])<<8 + int(in[17])
		urgent := int(in[18])<<8 + int(in[19])
		var options []byte
		if offset > 5 {
			options = in[20:offset*4]
		}
		data_hex := in[offset*4:]
		data_str := string(data_hex)
		out = fmt.Sprintf("TCP:\n source port: % d\n destination port: % d\n", sourcePort, destPort) + 
				fmt.Sprintf("seq: % d\n ack: % d\n", seq, ack) +			
				fmt.Sprintf("offset: % d\n reserved: % d\n flags: % b\n", offset, reserved, flags) +
				fmt.Sprintf("window size: % d\n checksum: % X\n urg: % d\n options: % X\n", windowSize, checksum, urgent, options) +		
				fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)
	} else {
		sourceIP := int(int(in[0])<<24 + int(in[1])<<16 + int(in[2])<<8 + int(in[3]))
		destIP := int(int(in[4])<<24 + int(in[5])<<16 + int(in[6])<<8 + int(in[7]))
		zeroes := int(in[8])
		tcpLength := int(int(in[10])<<8 + int(in[11]))
		sourcePort := int(int(in[12])<<8 + int(in[13]))
		destPort := int(int(in[14])<<8 + int(in[15]))
		seq := int(int(in[16])<<24 + int(in[17])<<16 + int(in[18])<<8 + int(in[19]))
		ack := int(int(in[20])<<24 + int(in[21])<<16 + int(in[22])<<8 + int(in[23]))
		offset := int(in[24]>>4)
		reserved := int(in[24]&0x0f)
		flags := int(in[25])
		windowSize := int(in[26])<<8 + int(in[27])
		checksum := int(in[28])<<8 + int(in[29])
		urgent := int(in[30])<<8 + int(in[31])
		var options []byte
		if offset >5 {
			options = in[32:offset*4+12]
		}
		data_hex := in[offset*4+12:]
		data_str := string(data_hex)
		out = fmt.Sprintf("TCP IPv4:\n sourceIP: % s\n destIP: % s\n", intToIP(sourceIP), intToIP(destIP)) +
				fmt.Sprintf("zeroes: % d\n protocol: % d\n tcp length: % d\n", zeroes, protocol, tcpLength) +
				fmt.Sprintf("UDP:\n source port: % d\n destination port: % d\n", sourcePort, destPort) + 
				fmt.Sprintf("seq: % d\n ack: % d\n", seq, ack) +			
				fmt.Sprintf("offset: % d\n reserved: % d\n flags: % b\n", offset, reserved, flags) +
				fmt.Sprintf("window size: % d\n checksum: % X\n urg: % d\n options: % X\n", windowSize, checksum, urgent, options) +		
				fmt.Sprintf("data (HEX): % X\n data (STR): % s\n\n", data_hex, data_str)	
	}
	return
}

func listenIP(messages chan string) {
	rawMessages := make(chan string)
	go listenICMP(rawMessages, true)
	go listenUDP(rawMessages, true)
	go listenTCP(rawMessages, true)
	for {
		msg := <-rawMessages
		fmtMsg := formatIP([]byte(msg))
		messages <- fmtMsg
	}
}

func listenICMP(messages chan string, raw bool) {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		tmpBuf := make([]byte, 1024)
		numRead, err := f.Read(tmpBuf)
		if err != nil {
			fmt.Println(err)
		}
		msg := tmpBuf[:numRead]
		if (raw) {
			messages <- string(msg)
		} else {
			fmtMsg := formatICMP(msg)
			messages <- fmtMsg
		}
	}
}

func listenUDP(messages chan string, raw bool) {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		tmpBuf := make([]byte, 1024)
		numRead, err := f.Read(tmpBuf)
		if err != nil {
			fmt.Println(err)
		}
		msg := tmpBuf[:numRead]
		if (raw) {
			messages <- string(msg)
		} else {
			fmtMsg := formatUDP(msg)
			messages <- fmtMsg
		}
	}
}

func listenTCP(messages chan string, raw bool) {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		tmpBuf := make([]byte, 1024)
		numRead, err := f.Read(tmpBuf)
		if err != nil {
			fmt.Println(err)
		}
		msg := tmpBuf[:numRead]
		if (raw) {
			messages <- string(msg)
		} else {
			fmtMsg := formatTCP(msg)
			messages <- fmtMsg
		}
	}
}