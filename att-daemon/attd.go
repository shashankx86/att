package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gen2brain/beeep"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

var pipePath string
const iconPath = "./assets/ico.png"

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
	notify("attd", "Arcade Time Tracker Daemon", fmt.Sprintf("Starting daemon with pipe path: %s", pipePath))

	// Ensure the pipe file does not already exist (Unix-like systems)
	if runtime.GOOS != "windows" {
		if _, err := os.Stat(pipePath); err == nil {
			fmt.Printf("Pipe file already exists, removing: %s\n", pipePath)
			notify("attd", "Arcade Time Tracker Daemon", fmt.Sprintf("Pipe file already exists, removing: %s", pipePath))
			os.Remove(pipePath)
		}
	}

	// Create a Named Pipe listener
	listener, err := net.Listen("unix", pipePath)
	if err != nil {
		fmt.Printf("Failed to listen on pipe: %v\n", err)
		notify("attd", "Arcade Time Tracker Daemon", fmt.Sprintf("Failed to listen on pipe: %v", err))
		return
	}
	defer listener.Close()

	// Ensure the pipe file is removed on exit (Unix-like systems)
	if runtime.GOOS != "windows" {
		defer os.Remove(pipePath)
	}

	fmt.Println("Daemon started and listening on", pipePath)
	notify("attd", "Arcade Time Tracker Daemon", fmt.Sprintf("Daemon started and listening on %s", pipePath))

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
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Failed to read from connection: %v\n", err)
		return
	}

	input := string(buf[:n])
	fmt.Printf("Received input: %s\n", input)

	// Parse the input as JSON
	var payload struct {
		Work    string `json:"work"`
		SlackID string `json:"slack_id"`
		APIKey  string `json:"api_key"`
	}
	err = json.Unmarshal([]byte(input), &payload)
	if err != nil {
		fmt.Printf("Failed to parse input as JSON: %v\n", err)
		conn.Write([]byte(fmt.Sprintf("Failed to parse input as JSON: %v\n", err)))
		return
	}

	// Perform the API GET request to get the latest session
	endTime, err := getSessionEndTime(payload.SlackID, payload.APIKey)
	if err != nil {
		fmt.Printf("Failed to get session end time: %v\n", err)
		conn.Write([]byte(fmt.Sprintf("Failed to get session end time: %v\n", err)))
		return
	}

	// Perform the API POST request to start a new session
	respStatus, respBody := postToAPI(payload.Work, payload.SlackID, payload.APIKey)
	response := fmt.Sprintf("Response Status: %s\nResponse Body: %s\n", respStatus, respBody)
	fmt.Println("API request made, sending response back to sender")

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Failed to write to connection: %v\n", err)
	}

	// Send a push notification based on the response
	handleNotification(respBody, payload.Work, endTime)
	time.Sleep(1 * time.Second) // Adding delay to ensure response is sent before the client closes the connection
}

func getSessionEndTime(slackID, apiKey string) (time.Time, error) {
	url := fmt.Sprintf("https://hackhour.hackclub.com/api/session/%s", slackID)

	// Create the HTTP request
	req, _ := http.NewRequest("GET", url, nil)

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to perform API request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, _ := io.ReadAll(resp.Body)

	// Parse the response
	var response struct {
		OK   bool `json:"ok"`
		Data struct {
			EndTime string `json:"endTime"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if the response is OK
	if !response.OK {
		return time.Time{}, fmt.Errorf("API response not OK")
	}

	// Parse the endTime
	endTime, err := time.Parse(time.RFC3339, response.Data.EndTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse endTime: %w", err)
	}

	return endTime, nil
}

func postToAPI(work, slackID, apiKey string) (string, string) {
	url := fmt.Sprintf("https://hackhour.hackclub.com/api/start/%s", slackID)

	// Prepare the JSON body
	body := map[string]string{
		"work": work,
	}
	jsonBody, _ := json.Marshal(body)

	// Create the HTTP request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	// Set the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Failed to perform API request", ""
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, _ := io.ReadAll(resp.Body)

	// Return the response status and body
	return resp.Status, string(respBody)
}

func handleNotification(respBody string, work string, endTime time.Time) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		notify("attd", "Arcade Time Tracker", respBody)
		return
	}

	ok, okPresent := response["ok"].(bool)
	if !okPresent || !ok {
		if errMsg, errPresent := response["error"].(string); errPresent {
			switch errMsg {
			case "Unauthorized":
				notify("attd", "Arcade Time Tracker", "Unauthorized: Invalid API Key or Slack ID")
			case "You already have an active session":
				notify("attd", "Arcade Time Tracker", "You already have an active session")
			default:
				notify("attd", "Arcade Time Tracker", respBody)
			}
		}
		return
	}

	message := fmt.Sprintf("Session started: %s", work)
	notify("attd", "Arcade Time Tracker", message)

	go setupNotifications(endTime)
}

func setupNotifications(endTime time.Time) {
	currentTime := time.Now()

	for currentTime.Before(endTime) {
		timeRemain := int(endTime.Sub(currentTime).Minutes())
		if timeRemain <= 0 {
			break
		}

		switch {
		case timeRemain > 10 && timeRemain%20 == 0:
			notify("attd", "Arcade Time Tracker", fmt.Sprintf("You have %d minutes left!", timeRemain))
			time.Sleep(20 * time.Minute)
		case timeRemain == 10:
			notify("attd", "Arcade Time Tracker", "Just 10 minutes left!")
			time.Sleep(10 * time.Minute)
		case timeRemain == 5:
			notify("attd", "Arcade Time Tracker", "The last 5 minutes!")
			time.Sleep(5 * time.Minute)
		default:
			time.Sleep(1 * time.Minute)
		}

		currentTime = time.Now()
	}
}

func notify(appName, title, message string) {
	err := beeep.Notify(title, message, iconPath)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
	}
}
