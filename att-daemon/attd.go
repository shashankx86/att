package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
)

var pipePath string

func init() {
	if runtime.GOOS == "windows" {
		pipePath = `\\.\pipe\attd`
	} else {
		pipePath = "/tmp/attd"
	}
}

func main() {
	// Define the pipePath flag
	var pipePathFlag string
	flag.StringVar(&pipePathFlag, "pipe-path", pipePath, "set the path for the pipe")
	flag.Parse()

	// If the flag is provided, update the pipePath
	if pipePathFlag != "" {
		pipePath = pipePathFlag
	}

	fmt.Printf("Starting daemon with pipe path: %s\n", pipePath)

	// Ensure the pipe file does not already exist (Unix-like systems)
	if runtime.GOOS != "windows" {
		if _, err := os.Stat(pipePath); err == nil {
			fmt.Printf("Pipe file already exists, removing: %s\n", pipePath)
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

	fmt.Println("Daemon started and listening on", pipePath)

	for {
		// Accept new connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		fmt.Println("New connection accepted")

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling new connection")

	// Read the command from the client
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read from connection: %v\n", err)
		return
	}

	input := string(buf[:n])
	fmt.Printf("Received input: %s\n", input)

	// Split the input into work, SLACK_ID, and API_KEY
	var work, slackID, apiKey string
	fmt.Sscanf(input, "%s %s %s", &work, &slackID, &apiKey)

	// Perform the API POST request
	err = postToAPI(work, slackID, apiKey)
	if err != nil {
		fmt.Printf("Failed to perform API request: %v\n", err)
		conn.Write([]byte("Failed to perform API request\n"))
	} else {
		conn.Write([]byte("API request successful\n"))
	}
}

func postToAPI(work, slackID, apiKey string) error {
	url := fmt.Sprintf("https://hackhour.hackclub.com/api/start/%s", slackID)

	// Prepare the JSON body
	body := map[string]string{
		"work": work,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	return nil
}
