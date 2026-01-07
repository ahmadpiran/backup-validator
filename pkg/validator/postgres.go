package validator

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type PostgresValidator struct {
}

func (p *PostgresValidator) Validate(backupPath string) error {
	// 1. Start container
	fmt.Println("[Postgres] Starting container...")
	runCmd := exec.Command("docker", "run",
		"--rm", "-d",
		"-e", "POSTGRES_PASSWORD=mysecret",
		"-p", "5433:5432",
		"postgres:15-alpine",
	)
	out, err := runCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	containerID := strings.TrimSpace(string(out))
	if len(containerID) > 12 {
		containerID = containerID[:12]
	}
	fmt.Printf("[Postgres] Container ID: %s\n", containerID)
	defer func() {
		fmt.Println("[Postgres] Cleaning up container...")
		exec.Command("docker", "stop", containerID).Run()
	}()

	// 2. Wait for Readiness
	fmt.Println("[Postgres] Waiting for readiness...")
	ready := false
	for i := 0; i < 10; i++ {
		if exec.Command("docker", "exec", containerID, "pg_isreadey", "-U", "postgres").Run() == nil {
			ready = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	if !ready {
		return fmt.Errorf("database timed out")
	}

	// 3. Copy File
	fmt.Println("[Postgres] Copying backup file...")
	if err := exec.Command("docker", "cp", backupPath, containerID+"/tmp/backup.sql").Run(); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// 4. Restore & Verify
	fmt.Println("[Postgres] Running restore...")
	restoreCmd := exec.Command("docker", "exec", containerID, "psql", "-U", "postgres", "-v", "ON_ERROR_STOP=1", "-f", "/tmp/backup.sql")

	output, err := restoreCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("validation failed:\n%s", string(output))
	}

	fmt.Println("[Postgres] Restore successful!")
	return nil
}
