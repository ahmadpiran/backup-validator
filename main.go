package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide a path to the .sql file")
		fmt.Println("Usage: go run main.go <path-to-sql-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Just print it for now
	fmt.Printf("Received backup file path: %s\n", filePath)
}
