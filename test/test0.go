package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	// "os"
)

func main() {
	// Define the pipePath flag
	var pipePath string
	flag.StringVar(&pipePath, "pipe-path", "/tmp/attd", "set the path for the pipe")
	flag.Parse()

	// Prepare the JSON payload
	payload := map[string]string{
		"work":    "work on att",
		"slack_id": "U05XXXXXXXX",
		"api_key":  "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	// Connect to the Unix domain socket
	conn, err := net.Dial("unix", pipePath)
	if err != nil {
		fmt.Printf("Failed to connect to socket: %v\n", err)
		return
	}
	defer conn.Close()

	// Send the JSON payload
	_, err = conn.Write(jsonPayload)
	if err != nil {
		fmt.Printf("Failed to write to socket: %v\n", err)
		return
	}

	// Receive the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read from socket: %v\n", err)
		return
	}

	response := string(buf[:n])
	fmt.Printf("Received response: %s\n", response)
}
