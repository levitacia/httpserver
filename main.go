package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "localhost:6969")
	if err != nil {
		fmt.Println("Failed to bind to port 6969")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	fmt.Println(req)
	if strings.HasPrefix(req, "GET /user-agent") {
		userAgent := ""
		for {
			line, err := reader.ReadString('\n')
			if err != nil || line == "\r\n" {
				break
			}
			if strings.HasPrefix(line, "User-Agent:") {
				userAgent = strings.TrimPrefix(line, "User-Agent: ")
				userAgent = strings.TrimSpace(userAgent)
				break
			}
		}
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		conn.Write([]byte(response))
		return
	} else if strings.HasPrefix(req, "GET /echo/") {
		parts := strings.Split(req, " ")
		if len(parts) > 1 {
			urlPart := parts[1]
			urlParts := strings.Split(urlPart, "/")
			if len(urlParts) > 2 {
				str := urlParts[2]
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
				conn.Write([]byte(response))
				return
			}
		}
	} else if strings.HasPrefix(req, "GET / HTTP/1.1") {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}
