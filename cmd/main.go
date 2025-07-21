package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . [command]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "db":
		handleDBCommand(os.Args[2:])
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func handleDBCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go run . db [--migrate|--rollback|--drop|--version]")
		return
	}

	// Build DSN from individual env vars
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)

	base := []string{
		"-path", "migrations",
		"-database", dsn,
	}

	var cmd *exec.Cmd

	switch args[0] {
	case "--migrate":
		cmd = exec.Command("migrate", append(base, "up")...)
	case "--rollback":
		cmd = exec.Command("migrate", append(base, "down", "1")...)
	case "--drop":
		cmd = exec.Command("migrate", append(base, "drop", "-f")...)
	case "--version":
		cmd = exec.Command("migrate", append(base, "version")...)
	default:
		fmt.Println("Unknown db command:", args[0])
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ migration command failed: %v", err)
	}
}
