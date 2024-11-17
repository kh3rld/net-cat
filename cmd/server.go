package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
func sendHistory(client Client) {
	history := "Previous messages:\n"
	for _, c := range clients {
		history += fmt.Sprintf("[%s]: Hello from %s\n", time.Now().Format("2006-01-02 15:04:05"), c.name)
	}
	client.writer.WriteString(history)
	client.writer.Flush()
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

// Handle communication with a client
func handleClient(conn net.Conn) {
	fmt.Fprintf(conn, "Welcome to TCP-Chat!\n")
	fmt.Fprintf(conn, "         _nnnn_\n")
	fmt.Fprintf(conn, "        dGGGGMMb\n")
	fmt.Fprintf(conn, "       @p~qp~~qMb\n")
	fmt.Fprintf(conn, "       M|@||@) M|\n")
	fmt.Fprintf(conn, "       @,----.JM|\n")
	fmt.Fprintf(conn, "      JS^\\__/  qKL\n")
	fmt.Fprintf(conn, "     dZP        qKRb\n")
	fmt.Fprintf(conn, "    dZP          qKKb\n")
	fmt.Fprintf(conn, "   fZP            SMMb\n")
	fmt.Fprintf(conn, "   HZM            MMMM\n")
	fmt.Fprintf(conn, "   FqM            MMMM\n")
	fmt.Fprintf(conn, " __| \".        |\\dS\"qML\n")
	fmt.Fprintf(conn, " |    `.       | `' \\Zq\n")
	fmt.Fprintf(conn, "_)      \\.___.,|     .'\n")
	fmt.Fprintf(conn, "\\____   )MMMMMP|   .'\n")
	fmt.Fprintf(conn, "     `-'       `--'\n")
	fmt.Fprintf(conn, "[ENTER YOUR NAME]: ")

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	clientName := scanner.Text()
	clientName = strings.TrimSpace(clientName)

	if clientName == "" {
		fmt.Fprintf(conn, "[ERROR] Name cannot be empty.\n")
		conn.Close()
		return
	}
	client := Client{
		name:   clientName,
		conn:   conn,
		writer: bufio.NewWriter(conn),
	}

	mu.Lock()
	clients = append(clients, client)
	mu.Unlock()

	broadcast(fmt.Sprintf("%s has joined the chat...\n", client.name))

	sendHistory(client)

	go readMessages(client)
	select {}
}
