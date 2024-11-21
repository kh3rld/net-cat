# NetCat

This project is a simple implementation of a group chat system, modeled after the NetCat utility, using a Server-Client architecture in Go. The server listens on a specified port and accepts incoming connections from clients. Clients can send messages to the group, and all connected clients receive the messages in real-time. The system supports multiple clients, ensuring that no client is disconnected when another leaves.

## Features

- **TCP Connections**: A server that listens for incoming TCP connections on a specified port and supports communication with multiple clients.
 - **Client Management**: Each client is required to provide a name, and the server handles up to 10 concurrent client connections.
 - **Message Broadcasting**: Messages sent by clients are broadcasted to all other connected clients. Each message is prefixed with a timestamp and the sender's name, e.g., [2020-01-20 15:48:41][client.name]:[client.message].
 - **Client Join/Exit Notifications**: When a client joins or leaves the chat, all other clients are notified by the server.
 - **Message History**: New clients that join the chat are provided with the previous chat history.
 - **Error Handling**: Proper error handling for both the server and client sides.
 - **Port Flexibility**: Default port is 8989 if no port is specified. Otherwise, it listens on the provided port.
  - **Concurrency**: Utilizes Go routines and channels (or mutexes) to handle concurrency.

## Requirements

   - Go (Golang): Make sure Go is installed on your system.
   - Maximum 10 Clients: The server can handle up to 10 simultaneous client connections.

## Installation

Ensure you have go version 1.23 installed in your machine. Clone the repo and navigate to the project directory.

```bash
git clone 
cd netcat/cmd
```

## Usage

To run the server you can use two ways:

 1. ```go run .```  This runs the default port 8989
 2. ```go run 2525``` This is a case where a port is specified, remember you can specify any port you want to.

To run the client you can also use two ways: 

1. ```nc localhost 2525``` This starts a client in the host machine.
2. ```nc $IP 2525``` This you must specify the ip address of the host machine and the port it is running on.

## Testing

This program has got test files and to test run the following command in the right directory.

```go test -v```

## Error Handling

The server and client are designed to handle various errors, including:

 - Connection failures
 - Server-side disconnections
 - Invalid inputs (e.g., empty messages)
 - Maximum client connection limits

## Contributing

Feel free to submit issues or pull requests if you find any bugs or would like to contribute to the project. If you have suggestions for improvements, don't hesitate to create an issue in the repository.

## License

This project is open-source and available under the MIT License.