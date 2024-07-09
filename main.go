package main

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/binary"
	"github.com/spf13/cobra"
	"log"
	"os/user"
	"syscall"
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

// getTotalRAM retrieves the total RAM size in kilobytes
func getTotalRAM() (uint64, error) {
	var sysinfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysinfo)
	if err != nil {
		return 0, err
	}
	return sysinfo.Totalram * uint64(syscall.Getpagesize()) / 1024, nil
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

	// Write the length of the token followed by the obfuscated token bytes
	length := int32(len(tokenBytes))

	// Write the length as a 4-byte integer in little-endian format
	file, err := os.Create(configFile)
	if err != nil {
		log.Fatalf("Unable to create config file: %v", err)
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, length)
	if err != nil {
		log.Fatalf("Unable to write length: %v", err)
	}

	// Write the actual obfuscated token bytes
	_, err = file.Write(tokenBytes)
	if err != nil {
		log.Fatalf("Unable to write API token: %v", err)
	}

	fmt.Println("API token saved successfully.")
}

