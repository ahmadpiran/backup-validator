package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// Run the Postgres container
	fmt.Println("Starting PostgreSQL container...")

	// We use port 5433 on host to avoid conflict with local DBs
	// We use a simple password 'mysecret' for this validation env
	runCmd := exec.Command("docker", "run",
		"--rm",
		"-d",
		"-e", "POSTGRES_PASSWORD=mysecret",
		"-p", "5432:5432",
		"postgres:15-alpine",
	)

	out, dockerRunErr := runCmd.Output()

	if dockerRunErr != nil {
		fmt.Println("Error: failed to start container.")
		fmt.Printf("System error: %v\n", dockerRunErr)
		os.Exit(1)
	}

	// Docker adds a newline at the end of the ID, we trim it
	containerID := strings.TrimSpace(string(out))

	// We only take the first 12 chars of the ID for display, like Docker CLI does
	shortID := containerID
	if len(containerID) > 12 {
		shortID = containerID[:12]
	}

	fmt.Printf("Success: Container started. ID: %s\n", shortID)

	fmt.Println("Waiting for database to initialize...")

	maxRetries := 10
	isReady := false

	for i := 0; i < maxRetries; i++ {
		checkCmd := exec.Command("docker", "exec", containerID, "pg_isready", "-U", "postgres")

		if checkCmdErr := checkCmd.Run(); checkCmdErr == nil {
			isReady = true
			fmt.Println("Database is ready!")
			break
		}

		fmt.Printf("Database not ready yet... retrying (%d/%d)\n", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if !isReady {
		fmt.Println("Error: Database timed out and did not start.")
		exec.Command("docker", "stop", containerID).Run()
		os.Exit(1)
	}

	fmt.Println("Copying backup file to container...")
	destination := containerID + ":/tmp/backup.sql"

	cpCmd := exec.Command("docker", "cp", filePath, destination)

	if cpCmdErr := cpCmd.Run(); cpCmdErr != nil {
		fmt.Printf("Error: Failed to copy file to container: %v\n", cpCmdErr)
		exec.Command("docker", "stop", containerID).Run()
		os.Exit(1)
	}

	fmt.Println("Success: File copied to /tmp/backup.sql inside container.")

}
