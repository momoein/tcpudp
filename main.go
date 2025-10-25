package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	mode := flag.String("mode", "server", "server or client")
	proto := flag.String("proto", "tcp", "protocol: tcp or udp")
	addr := flag.String("addr", "127.0.0.1:11000", "address to listen/connect")
	file := flag.String("file", "image.jpg", "image file to serve or save as")
	flag.Parse()

	switch *proto {
	case "tcp":
		if *mode == "server" {
			runTCPServer(*addr, *file)
		} else {
			runTCPClient(*addr, *file)
		}
	case "udp":
		if *mode == "server" {
			runUDPServer(*addr, *file)
		} else {
			runUDPClient(*addr, *file)
		}
	default:
		fmt.Println("Unknown protocol:", *proto)
	}
}

// --- TCP VERSION ---
func runTCPServer(addr, file string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("[TCP SERVER] Listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleTCPConnection(conn, file)
	}
}

func handleTCPConnection(conn net.Conn, file string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	if cmd == "GET IMAGE\n" {
		sendFileOverConn(conn, file)
	} else {
		conn.Write([]byte("Unknown command\n"))
	}
}

func runTCPClient(addr, file string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("GET IMAGE\n"))

	out, err := os.Create("received_tcp.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	start := time.Now()

	n, err := io.Copy(out, conn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[TCP CLIENT] Received %d bytes in %v\n", n, time.Since(start))
}

// --- UDP VERSION ---
func runUDPServer(addr, file string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("[UDP SERVER] Listening on", addr)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}
		cmd := string(buf[:n])
		if cmd == "GET IMAGE" {
			sendFileOverUDP(conn, clientAddr, file)
		} else {
			fmt.Println("Unknown UDP command:", cmd)
		}
	}
}

func runUDPClient(addr, file string) {

	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("GET IMAGE"))

	out, err := os.Create("received_udp.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	start := time.Now()

	buf := make([]byte, 2048)
	total := 0
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		if n == 3 && string(buf[:n]) == "EOF" {
			break
		}
		total += n
		out.Write(buf[:n])
	}

	fmt.Printf("[UDP CLIENT] Received %d bytes in %v\n", total, time.Since(start))
}

// --- FILE TRANSFER HELPERS ---
func sendFileOverConn(conn net.Conn, file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("File open error:", err)
		return
	}
	defer f.Close()

	n, err := io.Copy(conn, f)
	if err != nil {
		fmt.Println("Send error:", err)
		return
	}
	fmt.Printf("Sent %d bytes to %s\n", n, conn.RemoteAddr())
}

func sendFileOverUDP(conn *net.UDPConn, addr *net.UDPAddr, file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("File open error:", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	total := 0
	for {
		n, err := f.Read(buf)
		if n > 0 {
			conn.WriteToUDP(buf[:n], addr)
			total += n
			// time.Sleep(1 * time.Millisecond) // to avoid packet loss
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
	}

	for range 5 {
		conn.WriteToUDP([]byte("EOF"), addr)
	}

	fmt.Printf("Sent %d bytes over UDP to %s\n", total, addr)
}
