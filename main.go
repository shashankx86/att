package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	var apiToken string
	var slackID string

	// Define the root command
	var rootCmd = &cobra.Command{
		Use:   "att",
		Short: "Arcade Time Tracker",
	}

	// Define the configure command
	var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure the CLI tool",
	}

	// Define the api-token sub-command
	var apiTokenCmd = &cobra.Command{
		Use:   "api-token [token]",
		Short: "Set the API token",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiToken = args[0]
			updateConfigData("api-token", apiToken)
		},
	}

	// Define the slack-id sub-command
	var slackIDCmd = &cobra.Command{
		Use:   "slack-id [id]",
		Short: "Set the Slack ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			slackID = args[0]
			updateConfigData("slack-id", slackID)
		},
	}

	// Add the sub-commands to the configure command
	configureCmd.AddCommand(apiTokenCmd)
	configureCmd.AddCommand(slackIDCmd)

	// Define the session command
	var sessionCmd = &cobra.Command{
		Use:   "session",
		Short: "Manage sessions",
	}

	// Define reusable function to fetch and print data from the API
	fetchAndPrintData := func(endpoint string) {
		configData := loadConfigData()
		apiToken = configData["api-token"]
		slackID = configData["slack-id"]

		if apiToken == "" || slackID == "" {
			fmt.Println("Please set your API token and Slack ID using the configure command.")
			return
		}

		// Make the API request
		url := fmt.Sprintf("https://hackhour.hackclub.com/api/%s/%s", endpoint, slackID)
		resp, err := makeAPIRequest("GET", url, nil, apiToken)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Error: received status code %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		handleError("Unable to read response body", err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		handleError("Unable to unmarshal response", err)

		// Print the response data
		if result["ok"].(bool) {
			data, _ := json.MarshalIndent(result["data"], "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Println("Error: unable to fetch data")
		}
	}

	// Define the list sub-command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the latest session",
		Run: func(cmd *cobra.Command, args []string) {
			fetchAndPrintData("session")
		},
	}

	// Define the stats sub-command
	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Get the stats for the user",
		Run: func(cmd *cobra.Command, args []string) {
			fetchAndPrintData("stats")
		},
	}

	// Define the goals sub-command
	var goalsCmd = &cobra.Command{
		Use:   "goals",
		Short: "Get the goals for the user",
		Run: func(cmd *cobra.Command, args []string) {
			fetchAndPrintData("goals")
		},
	}

	// Define the history sub-command
	var historyCmd = &cobra.Command{
		Use:   "history",
		Short: "Get the history for the user",
		Run: func(cmd *cobra.Command, args []string) {
			fetchAndPrintData("history")
		},
	}

	// Define the start sub-command
	var startCmd = &cobra.Command{
		Use:   "start [work]",
		Short: "Start a new session",
		Run: func(cmd *cobra.Command, args []string) {
			var work string
			if len(args) > 0 {
				work = args[0]
			} else {
				fmt.Print("Session Description: ")
				fmt.Scanln(&work)
			}
			startNewSession(work)
		},
	}

	// Define the pause sub-command
	var pauseCmd = &cobra.Command{
		Use:   "pause",
		Short: "Pause or resume the current session",
		Run: func(cmd *cobra.Command, args []string) {
			pauseOrResumeSession()
		},
	}

	// Define the cancel sub-command
	var cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel the current session",
		Run: func(cmd *cobra.Command, args []string) {
			cancelSession()
		},
	}

	// Add the sub-commands to the session command
	sessionCmd.AddCommand(listCmd)
	sessionCmd.AddCommand(statsCmd)
	sessionCmd.AddCommand(goalsCmd)
	sessionCmd.AddCommand(historyCmd)
	sessionCmd.AddCommand(startCmd)
	sessionCmd.AddCommand(pauseCmd)
	sessionCmd.AddCommand(cancelCmd)

	// Add the configure and session commands to the root command
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(sessionCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadConfigData loads the configuration data from ~/.att/config/data.json
func loadConfigData() map[string]string {
	usr, err := user.Current()
	handleError("Unable to get the current user", err)

	configFile := filepath.Join(usr.HomeDir, ".att", "config", "data.json")

	// Create a default empty config if the file doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return make(map[string]string)
	}

	data, err := ioutil.ReadFile(configFile)
	handleError("Unable to read config file", err)

	var configData map[string]string
	err = json.Unmarshal(data, &configData)
	handleError("Unable to unmarshal config data", err)

	return configData
}

// saveConfigData saves the configuration data in JSON format to ~/.att/config/data.json
func saveConfigData(configData map[string]string) {
	usr, err := user.Current()
	handleError("Unable to get the current user", err)

	configDir := filepath.Join(usr.HomeDir, ".att", "config")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0700)
		handleError("Unable to create config directory", err)
	}

	configFile := filepath.Join(configDir, "data.json")

	data, err := json.Marshal(configData)
	handleError("Unable to marshal config data", err)

	// Write the JSON data to the config file
	err = ioutil.WriteFile(configFile, data, 0600)
	handleError("Unable to write config file", err)

	fmt.Println("Configuration data saved successfully.")
}

// updateConfigData updates a specific key-value pair in the configuration data
func updateConfigData(key, value string) {
	configData := loadConfigData()
	configData[key] = value
	saveConfigData(configData)
}

// handleError is a reusable function to handle errors
func handleError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

// makeAPIRequest is a reusable function to make API requests
func makeAPIRequest(method, url string, payload []byte, apiToken string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("Unable to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request: %v", err)
	}

	return resp, nil
}

// startNewSession sends a POST request to start a new session
func startNewSession(work string) {
	configData := loadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("https://hackhour.hackclub.com/api/start/%s", slackID)
	payload := map[string]string{"work": work}
	payloadBytes, err := json.Marshal(payload)
	handleError("Unable to marshal request payload", err)

	resp, err := makeAPIRequest("POST", url, payloadBytes, apiToken)
	handleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	handleError("Unable to unmarshal response", err)

	// Check for active session
	if resp.StatusCode != http.StatusOK {
		if errorMessage, exists := result["error"].(string); exists && errorMessage == "You already have an active session" {
			fmt.Println("You already have an active session")
		} else {
			fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		}
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}

// pauseOrResumeSession sends a POST request to pause or resume the current session
func pauseOrResumeSession() {
	configData := loadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("https://hackhour.hackclub.com/api/pause/%s", slackID)
	resp, err := makeAPIRequest("POST", url, nil, apiToken)
	handleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	handleError("Unable to unmarshal response", err)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}

// cancelSession sends a POST request to cancel the current session
func cancelSession() {
	configData := loadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("https://hackhour.hackclub.com/api/cancel/%s", slackID)
	resp, err := makeAPIRequest("POST", url, nil, apiToken)
	handleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	handleError("Unable to unmarshal response", err)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}
