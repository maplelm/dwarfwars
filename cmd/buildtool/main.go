package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	buildOnly := flag.Bool("b", false, "") // Default behavior is build and run
	target := flag.String("t", "server", "")
	output := flag.String("o", "", "")
	exeName := flag.String("n", "exe", "")

	flag.Parse()

	fmt.Println("Dwarf Wars Build Tool...")
	fmt.Printf("\tBuild Only: %t\n", *buildOnly)
	fmt.Printf("\tTarget: %s\n", *target)
	if len(*output) > 0 {
		fmt.Printf("\tOutput Dir: %s\n", *output)
	} else {
		fmt.Printf("\tOutput Dir: Default\n")
	}

	switch strings.ToLower(*target) {
	case "server":
		if len(*output) == 0 {
			*output = "./bin/server/"
		}
		fmt.Printf("Running Command (%s): %s\n", *target, fmt.Sprintf("go build -o %s ./cmd/server/", filepath.Join(*output, *exeName)))
		c := exec.Command("go", fmt.Sprintf("build -o %s ./cmd/server/", filepath.Join(*output, *exeName)))
		if o, err := c.Output(); err != nil {
			fmt.Printf("Failed To Build Go Target (%s): %s\n\n%s\n", *target, err, o)
			os.Exit(1)
		} else {
			fmt.Printf("Target Built (%s): %s\n", *target, o)
		}
		if !*buildOnly {
			fmt.Printf("Running Command (%s): %s\n", *target, fmt.Sprintf("./%s", *exeName))
			c = exec.Command(fmt.Sprintf("./%s", *exeName))
			c.Dir = *output
			if o, err := c.Output(); err != nil {
				fmt.Printf("Failed To Run Go Target (%s): %s\n\t%s\n", *target, err, o)
				os.Exit(1)
			} else {
				fmt.Printf("Ran Target (%s): %s\n", *target, o)
			}
		}
	case "client":
		if len(*output) == 0 {
			*output = "./bin/client/"
		}
		fmt.Printf("Running Command (%s): %s\n", *target, fmt.Sprintf("go build -o %s ./cmd/client/", filepath.Join(*output, *exeName)))
		c := exec.Command("go", "build", "-o", filepath.Join(*output, *exeName), "./cmd/client/")
		if o, err := c.Output(); err != nil {
			fmt.Printf("Failed to Build Go Target: %s\n\t%s\n", *output, err)
			os.Exit(1)
		} else {
			fmt.Printf("Target Built (%s): %s\n", *target, o)
		}
		if !*buildOnly {
			fmt.Printf("Running Command (%s): %s\n", *target, fmt.Sprintf("./%s", *exeName))
			c = exec.Command(fmt.Sprintf("./%s", *exeName))
			c.Dir = *output
			if err := c.Run(); err != nil {
				fmt.Printf("Failed To Run Go Target: %s\n\t%s\n", *output, err)
				os.Exit(1)
			}
		}
	default:
		fmt.Printf("Target Unrecognized: %s\n", strings.ToLower(*target))
		fmt.Println("Valid Targets:")
		for _, v := range []string{"server", "client"} {
			fmt.Printf("\t> %s\n", v)
		}
	}
}
