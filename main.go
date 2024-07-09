package main

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/binary"
	"github.com/spf13/cobra"
	"log"
	"os/user"
)

func main() {
	var apiToken string

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
			saveAPIToken(apiToken)
		},
	}

	// Add the sub-command to the configure command
	configureCmd.AddCommand(apiTokenCmd)
	// Add the configure command to the root command
	rootCmd.AddCommand(configureCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// saveAPIToken saves the API token in binary format to ~/.att/config/
func saveAPIToken(token string) {
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

	configFile := filepath.Join(configDir, "api-token")

	// Write the token in binary format
	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Unable to create config file: %v", err)
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, []byte(token))
	if err != nil {
		log.Fatalf("Unable to write API token: %v", err)
	}

	fmt.Println("API token saved successfully.")
}
