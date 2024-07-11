package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"att/handler"
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
			handler.UpdateConfigData("api-token", apiToken)
		},
	}

	// Define the slack-id sub-command
	var slackIDCmd = &cobra.Command{
		Use:   "slack-id [id]",
		Short: "Set the Slack ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			slackID = args[0]
			handler.UpdateConfigData("slack-id", slackID)
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
			handler.FetchAndPrintData("session")
		},
	}

	// Define the stats sub-command
	var statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Get the stats for the user",
		Run: func(cmd *cobra.Command, args []string) {
			handler.FetchAndPrintData("stats")
		},
	}

	// Define the goals sub-command
	var goalsCmd = &cobra.Command{
		Use:   "goals",
		Short: "Get the goals for the user",
		Run: func(cmd *cobra.Command, args []string) {
			handler.FetchAndPrintData("goals")
		},
	}

	// Define the history sub-command
	var historyCmd = &cobra.Command{
		Use:   "history",
		Short: "Get the history for the user",
		Run: func(cmd *cobra.Command, args []string) {
			handler.FetchAndPrintData("history")
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
			handler.StartNewSession(work)
		},
	}

	// Define the pause sub-command
	var pauseCmd = &cobra.Command{
		Use:   "pause",
		Short: "Pause or resume the current session",
		Run: func(cmd *cobra.Command, args []string) {
			handler.PauseOrResumeSession()
		},
	}

	// Define the cancel sub-command
	var cancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel the current session",
		Run: func(cmd *cobra.Command, args []string) {
			handler.CancelSession()
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