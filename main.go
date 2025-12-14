package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Error: Please provide a path to the .sql file")
		fmt.Println("Usage: go run main.go <path-to-sql-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	fi, fileStatErr := os.Stat(filePath)

	if fileStatErr != nil {
		if os.IsNotExist(fileStatErr) {
			fmt.Printf("Error: file not found at %s\n", filePath)
			os.Exit(1)
		}

		fmt.Printf("Error checking file %v\n", fileStatErr)
		os.Exit(1)
	}

	if fi.IsDir() {
		fmt.Printf("Error: Path %s is a directory, not a file.\n", filePath)
		os.Exit(1)
	}
	// Just print it for now
	fmt.Printf("Success: Found valid file: %s (Size: %d bytes)\n", filePath, fi.Size())

	// Verify Docker is running
	fmt.Println("Checking Docker connectivity...")

	cmd := exec.Command("docker", "info")
	cmd.Stdout = nil
	cmd.Stderr = nil

	dockerInfoErr := cmd.Run()

	if dockerInfoErr != nil {
		fmt.Println("Error: Docker does not seem to be running or installed.")
		fmt.Printf("System error: %v\n", dockerInfoErr)
		fmt.Println("Please ensure Docker Engine/Desktop is started.")
		os.Exit(1)
	}
}
