package main

import (
	"fmt"
	"net"
	"os"

	"github.com/maplelm/dwarfwars/pkg/logging"
)

func main() {
	fmt.Println("Testing Game Server...")

	var (
		addr string = "127.0.0.1"
		port int    = 3000
	)

	logging.Info("Connecting to Server")
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		logging.Error(err, "Failed to connect to server")
		os.Exit(1)
	}
	defer conn.Close()

	logging.Info("Sending data")

	wn, err := conn.Write([]byte("ping"))
	if err != nil {
		logging.Error(err, "Failed to Write data to socket")
		os.Exit(1)
	}

	logging.Info("Reading Reply")
	var data []byte = []byte{}
	rn, err := conn.Read(data)
	if err != nil {
		logging.Error(err, "Failed to read response from server")
		os.Exit(1)
	}

	if wn != rn {
		logging.Warningf("written data does not match read data, W: %d, R: %d", wn, rn)
	}

	logging.Infof("Server: %s", string(data))
}
