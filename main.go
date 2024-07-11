package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
	"strings"
    "att/handler"
)

var VERSION = "0.0.1"

func PrintASCII() {
	fmt.Println("      _    _   ")
	fmt.Println(" ___ | |_ | |_ ")
	fmt.Println("| .'||  _||  _|")
  	fmt.Println("|__,||_|  |_|  \n")
}

func main() {
    var apiToken string
    var slackID string

    // Define the root command
    var rootCmd = &cobra.Command{
        Use:   "att",
        Short: "Arcade Time Tracker",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
            if cmd.Name() == "help" || len(os.Args) == 1 {
                PrintASCII()
            }
        },
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
		Use:   "start [work...]",
		Short: "Start a new session",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			work := strings.Join(args, " ")
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

    // Define the ping command
    var pingCmd = &cobra.Command{
        Use:   "ping",
        Short: "Ping the server",
        Run: func(cmd *cobra.Command, args []string) {
            handler.PingServer()
        },
    }

    // Define the status command
    var statusCmd = &cobra.Command{
        Use:   "status",
        Short: "Get the status of hack hour",
        Run: func(cmd *cobra.Command, args []string) {
            handler.FetchAndPrintStatus()
        },
    }

	// CLI Version
	var versionCmd = &cobra.Command{
        Use:   "version",
        Short: "Print CLI version and quit",
        Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(VERSION)
        },
    }

    // Add the configure, session, ping, and status commands to the root command
    rootCmd.AddCommand(configureCmd)
    rootCmd.AddCommand(sessionCmd)
    rootCmd.AddCommand(pingCmd)
    rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)

    // Execute the root command
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
