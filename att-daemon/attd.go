package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Define the path to the Unix domain socket (change as needed)
const socketPath = "/tmp/timer.sock"

// Timer struct to manage the timer state
type Timer struct {
	mu      sync.Mutex
	running bool
}

func main() {
	// Ensure the socket file does not already exist
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// Create a Unix domain socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("Failed to listen on socket: %v\n", err)
		return
	}
	defer listener.Close()

	// Ensure the socket file is removed on exit
	defer os.Remove(socketPath)

	// Initialize the timer
	timer := &Timer{}

	fmt.Println("Daemon started and listening on", socketPath)

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
