package main

/*
Connection Server:
This server will be used to facilitate connections between users and all the
other servers. This will also be the server that authenticates users before
they are allowed to connect to anything else.
*/

import (
	"fmt"
	"net"
)

const (
	connectionSecretLength int = 255
)

var (
	connections map[net.Addr][]byte
)

func main() {
	fmt.Println("Connection Server Starting...")
	fmt.Println("Connection Server Closing...")
}
