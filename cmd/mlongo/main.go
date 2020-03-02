package main

import (
	"fmt"
	"os"
)

const (
	DEFAULT_HOST        = "localhost"
	DEFAULT_PORT uint64 = 27017
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("url empty or invalid")
		os.Exit(1)
	}

	fmt.Println(args)

}
