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
)

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
		if len(clients) >= maxConn {
			conn.Close()
			mu.Unlock()
			continue
		}
		mu.Unlock()

		go handleClient(conn)
	}
}

// Broadcast messages to all clients
func broadcast(message string, excludeClient *Client) {
	mu.Lock()
	defer mu.Unlock()
	for _, client := range clients {
		if excludeClient != nil && client == *excludeClient {
			continue
		}
		if client.conn == nil {
			continue // Skip clients whose connection has been closed
		}
		client.writer.WriteString(message)
		client.writer.Flush()
	}
}

func sendHistory(client Client) {
	history := ""
	for _, c := range clients {
		history += fmt.Sprintf("[%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), c.name)
	}
	client.writer.WriteString(history)
	client.writer.Flush()
}

// Read messages from client and broadcast
func readMessages(client Client) {
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "/rename ") {
			// Handle name change command
			newName := strings.TrimSpace(strings.TrimPrefix(msg, "/rename "))
			if newName == "" {
				client.writer.WriteString("[ERROR] New name cannot be empty.\n")
				client.writer.Flush()
				continue
			}

			mu.Lock()
			// Check if the new name is already taken
			nameExists := false
			for _, c := range clients {
				if c.name == newName {
					nameExists = true
					break
				}
			}
			if nameExists {
				client.writer.WriteString("[ERROR] [Name already in use. Choose a different name]: ")
				client.writer.Flush()
				mu.Unlock()
				continue
			}

			oldName := client.name
			client.name = newName
			for i := range clients {
				if clients[i].conn == client.conn {
					clients[i].name = newName
					break
				}
			}
			mu.Unlock()

			broadcast(fmt.Sprintf("%s changed their name to %s.\n", oldName, newName), nil)
			continue
		}

		if msg != "" {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			message := fmt.Sprintf("[%s][%s]: %s\n", timestamp, client.name, msg)
			broadcast(message, nil)
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
	}

	mu.Lock()
	for i, c := range clients {
		if c.conn == client.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()

	broadcast(fmt.Sprintf("%s has left the chat...\n", client.name), nil)
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
	clientName := strings.TrimSpace(scanner.Text())
	for _, c := range clients {
		if clientName == c.name {
			fmt.Fprintf(conn, "[ERROR] Name already in use. Choose a different name.\n Press any key then Enter to exit: ")
			conn.Close()
			return
		}
	}

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
	for _, c := range clients {
		if c.name == client.name {
			// Broadcast that a new client has joined the chat, excluding the joining client
			broadcast(fmt.Sprintf("%s has joined the chat...\n", client.name), &client)
		}
	}

	sendHistory(client)

	go readMessages(client)
	select {}
}
