package main

import (
	"fmt"
	"os"
)

const defaultPort = "8989"

func main() {
	if len(os.Args) > 3 || (len(os.Args) == 2 && os.Args[1] == "help") {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	port := defaultPort
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	if len(os.Args) == 3 {
		port = os.Args[2]
	}

	startServer(port)
}
