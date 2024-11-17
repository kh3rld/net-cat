package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
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

// Broadcast messages to all clients
func broadcast(message string) {
	mu.Lock()
	defer mu.Unlock()
	for _, client := range clients {
		client.writer.WriteString(message)
		client.writer.Flush()
	}
}

// Read messages from client and broadcast
func readMessages(client Client) {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			message := fmt.Sprintf("[%s][%s]: %s\n", timestamp, client.name, msg)
			broadcast(message)
		}
	}
	mu.Lock()
	for i, c := range clients {
		if c.name == client.name {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()

	broadcast(fmt.Sprintf("%s has left the chat...\n", client.name))
	client.conn.Close()
}
