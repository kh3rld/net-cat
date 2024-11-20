package main

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_readMessages(t *testing.T) {
	type args struct {
		client Client
	}
	tests := []struct {
		name           string
		args           args
		initialClients []Client
		inputMessages  []string
		expectedOutput []string
	}{
		{
			name: "Client sends a regular message",
			args: args{
				client: Client{
					name:   "TestClient",
					writer: bufio.NewWriter(&strings.Builder{}),
				},
			},
			initialClients: []Client{},
			inputMessages:  []string{"Hello, everyone!"},
			expectedOutput: []string{
				"[TestClient]: Hello, everyone!\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the initial state
			clients = tt.initialClients

			// Mock the connection input and output
			input := strings.Join(tt.inputMessages, "\n")
			mockOutput := &strings.Builder{}
			reader := bufio.NewReader(strings.NewReader(input))
			writer := bufio.NewWriter(mockOutput)

			// Create a mock client
			client := tt.args.client
			client.conn = &mockConn{
				reader: reader,
				writer: mockOutput,
			}
			client.writer = writer

			// Add the client to the global list
			mu.Lock()
			clients = append(clients, client)
			mu.Unlock()

			// Synchronize goroutine completion
			var wg sync.WaitGroup
			wg.Add(1)

			// Run the function in a goroutine
			go func() {
				defer wg.Done()
				readMessages(client)
			}()

			// Write messages to the input
			for _, message := range tt.inputMessages {
				client.conn.(*mockConn).reader = bufio.NewReader(strings.NewReader(message + "\n"))
				time.Sleep(10 * time.Millisecond)
			}

			// Wait for the goroutine to complete
			wg.Wait()

			// Verify the output
			writer.Flush() // Ensure the writer flushes to the mockOutput
			output := mockOutput.String()
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain: %q, got: %q", expected, output)
				}
			}
		})
	}
}

// Mock connection to simulate net.Conn for testing
type mockConn struct {
	reader *bufio.Reader
	writer *strings.Builder
}

func (m *mockConn) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

func (m *mockConn) Write(p []byte) (int, error) {
	return m.writer.Write(p)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }
