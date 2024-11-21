package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const maxConn = 10

type Client struct {
	name   string
	conn   net.Conn
	writer *bufio.Writer
}

var (
	clients []Client
	mu      sync.Mutex
	chats   []string
)

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
		if len(clients) >= maxConn {
			conn.Close()
			mu.Unlock()
			continue
		}
		mu.Unlock()

		go handleClient(conn)
	}
}

func broadcast(message string, excludeClient *Client) {
	mu.Lock()
	defer mu.Unlock()
	for _, client := range clients {
		// Skip broadcasting the message to the client that is excluded
		if excludeClient != nil && client == *excludeClient {
			continue
		}
		if client.conn == nil {
			continue
		}
		client.writer.WriteString(message)
		client.writer.Flush()
	}
}

// Handle the communication with a client
func handleClient(conn net.Conn) {
	client := Client{conn: conn, writer: bufio.NewWriter(conn)}
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
	client.name = strings.TrimSpace(scanner.Text())

	if client.name == "" {
		fmt.Fprintf(conn, "[ERROR] Name cannot be empty.\n")
		conn.Close()
		return
	}

	// Check for name uniqueness
	mu.Lock()
	for _, c := range clients {
		if client.name == c.name {
			fmt.Fprintf(conn, "[ERROR] Name already in use.\n")
			mu.Unlock()
			conn.Close()
			return
		}
	}
	clients = append(clients, client)
	mu.Unlock()

	for _, historyMessage := range chats {
		fmt.Fprint(client.conn, historyMessage)
	}
	broadcast(fmt.Sprintf("%s has joined the chat...\n", client.name), &client)
	go readMessages(client)
}

// Read messages from a client and broadcast them to everyone except the sender
func readMessages(client Client) {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		message := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.name, msg)
		chats = append(chats, message)
		broadcast(message, &client)
		saveLog(message)
	}
	mu.Lock()
	for i, c := range clients {
		if c.name == client.name {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()
	broadcast(fmt.Sprintf("%s has left the chat...\n", client.name), &client)
	client.conn.Close()
}

// Save the message to a log file
func saveLog(message string) {
	saveLog, err := os.OpenFile("/tmp/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	defer saveLog.Close()
	_, err = saveLog.WriteString(message + "\n")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
