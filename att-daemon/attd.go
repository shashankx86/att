package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

var pipePath string

func init() {
	if runtime.GOOS == "windows" {
		pipePath = `\\.\pipe\attd`
	} else {
		pipePath = "/tmp/attd"
	}
}

// Timer struct to manage the timer state
type Timer struct {
	mu      sync.Mutex
	running bool
}

func main() {
	// Ensure the pipe file does not already exist (Unix-like systems)
	if runtime.GOOS != "windows" {
		if _, err := os.Stat(pipePath); err == nil {
			os.Remove(pipePath)
		}
	}

	// Create a Named Pipe listener
	listener, err := net.Listen("unix", pipePath)
	if err != nil {
		fmt.Printf("Failed to listen on pipe: %v\n", err)
		return
	}
	defer listener.Close()

	// Ensure the pipe file is removed on exit (Unix-like systems)
	if runtime.GOOS != "windows" {
		defer os.Remove(pipePath)
	}

	// Initialize the timer
	timer := &Timer{}

	fmt.Println("Daemon started and listening on", pipePath)

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn, timer)
	}
}

func handleConnection(conn net.Conn, timer *Timer) {
	defer conn.Close()

	// Read the command from the client
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read from connection: %v\n", err)
		return
	}

	command := string(buf[:n])
	switch command {
	case "start":
		startTimer(conn, timer)
	default:
		conn.Write([]byte("Unknown command\n"))
	}
}

func startTimer(conn net.Conn, timer *Timer) {
	timer.mu.Lock()
	defer timer.mu.Unlock()

	if timer.running {
		conn.Write([]byte("Timer already running\n"))
	} else {
		timer.running = true
		conn.Write([]byte("Timer started for 60 seconds\n"))

		// Start the timer in a new goroutine
		go func() {
			time.Sleep(60 * time.Second)
			timer.mu.Lock()
			timer.running = false
			timer.mu.Unlock()
		}()
	}
}
