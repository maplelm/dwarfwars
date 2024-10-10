package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/maplelm/dwarfwars/pkg/command"
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

	_, err = conn.Write(command.Command{
		Version: 1,
		Type:    1,
		Size:    uint16(len([]byte("This is a command"))),
		Data:    []byte("This is a command"),
	}.Marshal())
	if err != nil {
		fmt.Printf("Failed to Send command to server: %s\n", err)
	}

	fmt.Println("Reading Reply")
	var data []byte = make([]byte, 2024)
	var fullData []byte
	readCount := 0
	for {
		n, err := conn.Read(data)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			fmt.Printf("Failed to read data from connection: %s\n", err)
			os.Exit(1)
		}
		readCount += n
		fullData = append(fullData, data...)
	}
	//fullData = fullData[:readCount]

	fmt.Printf("Unmarshaling: ")
	for _, v := range fullData {
		fmt.Printf("%b ", v)
	}
	fmt.Println()
	cmd, err := command.Unmarshal(fullData)
	if err != nil {
		fmt.Printf("Failed to Unmarshal Command: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response:\n\tVersion: %d\n\tType: %d\n\tSize: %d\n\tData: %s\n", cmd.Version, cmd.Type, cmd.Size, string(data))
}
