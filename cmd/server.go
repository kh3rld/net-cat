package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

const maxConnections = 10

type Client struct {
	name   string
	conn   net.Conn
	writer *bufio.Writer
}

var clients []Client
var mu sync.Mutex
var msgChannel = make(chan string)

// Server function to handle incoming connections
func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer ln.Close()
	fmt.Println("Server started on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		mu.Lock()
		if len(clients) >= maxConnections {
			conn.Close()
			mu.Unlock()
			continue
		}
		mu.Unlock()

		go handleClient(conn)
	}
}
