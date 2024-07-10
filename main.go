package main

import (
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
			configData := loadConfigData()
			configData["api-token"] = apiToken
			saveConfigData(configData)
		},
	}

	// Define the slack-id sub-command
	var slackIDCmd = &cobra.Command{
		Use:   "slack-id [id]",
		Short: "Set the Slack ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			slackID = args[0]
			configData := loadConfigData()
			configData["slack-id"] = slackID
			saveConfigData(configData)
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

	// Define the list sub-command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the latest session",
		Run: func(cmd *cobra.Command, args []string) {
			configData := loadConfigData()
			apiToken = configData["api-token"]
			slackID = configData["slack-id"]

			if apiToken == "" || slackID == "" {
				fmt.Println("Please set your API token and Slack ID using the configure command.")
				return
			}

			// Make the API request
			url := fmt.Sprintf("https://hackhour.hackclub.com/api/session/%s", slackID)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("Unable to create request: %v", err)
			}

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Unable to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Fatalf("Error: received status code %d", resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Unable to read response body: %v", err)
			}

			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Fatalf("Unable to unmarshal response: %v", err)
			}

			// Print the session data
			if result["ok"].(bool) {
				data, _ := json.MarshalIndent(result["data"], "", "  ")
				fmt.Println(string(data))
			} else {
				fmt.Println("Error: unable to fetch session data")
			}
		},
	}

	// Add the list sub-command to the session command
	sessionCmd.AddCommand(listCmd)

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
	if err != nil {
		log.Fatalf("Unable to get the current user: %v", err)
	}

	configFile := filepath.Join(usr.HomeDir, ".att", "config", "data.json")

	// Create a default empty config if the file doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return make(map[string]string)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	var configData map[string]string
	err = json.Unmarshal(data, &configData)
	if err != nil {
		log.Fatalf("Unable to unmarshal config data: %v", err)
	}

	return configData
}

// saveConfigData saves the configuration data in JSON format to ~/.att/config/data.json
func saveConfigData(configData map[string]string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to get the current user: %v", err)
	}

	configDir := filepath.Join(usr.HomeDir, ".att", "config")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			log.Fatalf("Unable to create config directory: %v", err)
		}
	}

	configFile := filepath.Join(configDir, "data.json")

	data, err := json.Marshal(configData)
	if err != nil {
		log.Fatalf("Unable to marshal config data: %v", err)
	}

	// Write the JSON data to the config file
	err = ioutil.WriteFile(configFile, data, 0600)
	if err != nil {
		log.Fatalf("Unable to write config file: %v", err)
	}

	fmt.Println("Configuration data saved successfully.")
}
