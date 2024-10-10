package main

import (
	"fmt"
	"net"
	"os"
	//"github.com/maplelm/dwarfwars/pkg/logging"
)

func main() {
	fmt.Println("Testing Game Server...")

	var (
		addr string = "127.0.0.1"
		port int    = 3000
	)

	fmt.Println("Connecting to Server")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Printf("Failed to connect to server: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Sending data: 'ping'")

	wn, err := conn.Write([]byte("ping"))
	if err != nil {
		fmt.Printf("Failed to Write data to socket: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Reading Reply")
	var data []byte = make([]byte, 2000)
	rn, err := conn.Read(data)
	if err != nil {
		fmt.Printf("Failed to read response from server: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Response: %s\n", string(data))

	if wn != rn {
		fmt.Printf("written data does not match read data, W: %d, R: %d\n", wn, rn)
	}

	fmt.Printf("Server: %s\n", string(data))
}
