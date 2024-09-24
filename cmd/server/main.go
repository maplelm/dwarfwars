package main

import (
	"flag"
	"fmt"
)

func main() {
	var (
		headless   *bool
		configPath *string
		savepath   *string
	)
	fmt.Println("Starting Dwarf Wars Server")

	headless = flag.Bool("headless", false, "if true the server will not use a tui so it can be automated")
	configPath = flag.String("config", "./config/", "location of server settings and configuration files")
	savepath = flag.String("world", "./saves/World/", "location of all the saved data for a world/game")
	flag.Parse()

	if *headless {
	} else {
	}
}

func TuiMode() {}
func CliMode() {}
