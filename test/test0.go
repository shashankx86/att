package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
)

func main() {
	// Define the pipePath flag
	var pipePath string
	flag.StringVar(&pipePath, "pipe-path", "/tmp/attd", "set the path for the pipe")
	flag.Parse()

	// Define which command to test: "start" or "track"
	var command string
	flag.StringVar(&command, "command", "track", "specify the command to send (start or track)")
	flag.Parse()

	// Prepare the data payload based on the command
	var payload map[string]interface{}
	if command == "start" {
		payload = map[string]interface{}{
			"command": "start",
			"data": map[string]interface{}{
				"work":     "work on att",
				"slack_id": "xxxxxxxxx",
				"api_key":  "xxxxxxxxx",
			},
		}
	} else if command == "track" {
		payload = map[string]interface{}{
			"command": "track",
			"data": map[string]interface{}{
				"slack_id": "xxxxxxxxx",
				"api_key":  "xxxxxxxxxxx",
			},
		}
	} else {
		fmt.Printf("Invalid command specified: %s\n", command)
		return
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
