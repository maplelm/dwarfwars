package main

import (
	"fmt"
	"github.com/maplelm/dwarfwars/pkg/server"
	"os"
	_ "time"
)

func main() {
	fmt.Println("Starting Server...")
	serv, err := server.New("0.0.0.0", "3000")
	if err != nil {
		fmt.Printf("Failed to create Server, %s\n", err)
		os.Exit(1)
	}
	go serv.Start()
	//time.Sleep(time.Duration(1) * time.Second)
	serv.Stop()
	fmt.Println("Closing Server...")
}
