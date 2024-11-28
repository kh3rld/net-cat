package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

var chatss []string

// Mock function for broadcasting
var broadcasts = func(message string, sender *Client) {
	if sender != nil {
		chatss = append(chatss, message)
	}
}

// Mock function for saving logs
var saveLogs = func(message string) {
	chatss = append(chatss, message)
}

// mock function to read messages from the client connection and broadcast them
func readMessagess(client Client, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		message := fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), client.name, msg)
		chats = append(chats, message)
		broadcasts(message, &client)
		saveLogs(message)
	}
}

// Test function for readMessages
func TestReadMessages(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	client := Client{
		name: "client1",
		conn: clientConn,
	}
	chats = nil

	var wg sync.WaitGroup
	wg.Add(1)
	go readMessagess(client, &wg)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(serverConn, "Hello everyone!")
		serverConn.Close()
	}()

	wg.Wait()

	expectedMessage := fmt.Sprintf("[%s][%s]: Hello everyone!\n", time.Now().Format("2006-01-02 15:04:05"), client.name)

	if len(chatss) < 2 {
		t.Fatalf("Expected at least 2 messages, got %d", len(chatss))
	}
	if chatss[0] != expectedMessage {
		t.Errorf("Expected message: %s, but got: %s", expectedMessage, chatss[0])
	}
	if len(chatss) < 2 || chatss[1] != expectedMessage {
		t.Errorf("Message not broadcasted correctly. Chats: %v", chatss)
	}
}

func TestStartServer(t *testing.T) {
	go startServer("8080")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Expected no error connecting to the server, but got %v", err)
	}
	defer conn.Close()

	t.Log("Server accepted the connection successfully.")
}
