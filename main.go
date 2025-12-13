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

	fi, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: file not found at %s\n", filePath)
			os.Exit(1)
		}

		fmt.Printf("Error checking file %v\n", err)
		os.Exit(1)
	}

	if fi.IsDir() {
		fmt.Printf("Error: Path %s is a directory, not a file.\n", filePath)
		os.Exit(1)
	}
	// Just print it for now
	fmt.Printf("Success: Found valid file: %s (Size: %d bytes)\n", filePath, fi.Size())
}
