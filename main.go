package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	// "runtime"

	"github.com/shirou/gopsutil/mem"
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
			saveConfigData(apiToken, slackID)
		},
	}

	// Define the slack-id sub-command
	var slackIDCmd = &cobra.Command{
		Use:   "slack-id [id]",
		Short: "Set the Slack ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			slackID = args[0]
			saveConfigData(apiToken, slackID)
		},
	}

	// Add the sub-commands to the configure command
	configureCmd.AddCommand(apiTokenCmd)
	configureCmd.AddCommand(slackIDCmd)

	// Add the configure command to the root command
	rootCmd.AddCommand(configureCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// getTotalRAM retrieves the total RAM size in kilobytes
func getTotalRAM() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total / 1024, nil
}

// saveConfigData saves the API token and Slack ID in binary format to ~/.att/config/
func saveConfigData(token string, slackID string) {
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

	configFile := filepath.Join(configDir, "data.bin")

	// Get total RAM size in kilobytes
	totalRAM, err := getTotalRAM()
	if err != nil {
		log.Fatalf("Unable to get total RAM: %v", err)
	}

	// Convert total RAM size to bytes
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, totalRAM)

	// XOR the token bytes with the RAM size key to obfuscate the token
	tokenBytes := []byte(token)
	for i := range tokenBytes {
		tokenBytes[i] ^= key[i%len(key)]
	}

	// XOR the Slack ID bytes with the RAM size key to obfuscate the Slack ID
	slackIDBytes := []byte(slackID)
	for i := range slackIDBytes {
		slackIDBytes[i] ^= key[i%len(key)]
	}

	// Write the length of the token followed by the obfuscated token bytes and Slack ID bytes
	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Unable to create config file: %v", err)
	}
	defer file.Close()

	// Write the length of the token as a 4-byte integer in little-endian format
	tokenLength := int32(len(tokenBytes))
	err = binary.Write(file, binary.LittleEndian, tokenLength)
	if err != nil {
		log.Fatalf("Unable to write token length: %v", err)
	}

	// Write the actual obfuscated token bytes
	_, err = file.Write(tokenBytes)
	if err != nil {
		log.Fatalf("Unable to write API token: %v", err)
	}

	// Write the length of the Slack ID as a 4-byte integer in little-endian format
	slackIDLength := int32(len(slackIDBytes))
	err = binary.Write(file, binary.LittleEndian, slackIDLength)
	if err != nil {
		log.Fatalf("Unable to write Slack ID length: %v", err)
	}

	// Write the actual obfuscated Slack ID bytes
	_, err = file.Write(slackIDBytes)
	if err != nil {
		log.Fatalf("Unable to write Slack ID: %v", err)
	}

	fmt.Println("Configuration data saved successfully.")
}
