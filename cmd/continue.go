package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var continueYes bool

var continueCmd = &cobra.Command{
	Use:   "continue <entry-id>",
	Short: "Start a new timer using an existing entry as template",
	Long: `Start a new timer by copying the description and project from an existing entry.

The entry ID can be a short ID (6-36 characters) or full UUID.
Use 'soltty list' to see 8-character short IDs for recent entries.

Examples:
  soltty continue 985d7cb2
  soltty continue 985d7cb2-cb20
  soltty continue 985d7cb2-cb20-40a4-ad9a-627ffa5cdc77`,
	Args: cobra.ExactArgs(1),
	Run:  runContinue,
}

func init() {
	continueCmd.Flags().BoolVarP(&continueYes, "yes", "y", false, "Skip confirmation when a timer is already running")
}

func runContinue(cmd *cobra.Command, args []string) {
	entryID := args[0]

	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	// Look up the entry by short ID
	entry, err := c.FindEntryByShortID(entryID)
	if err != nil {
		fmt.Println(formatError(err))
		os.Exit(1)
	}

	// Check if timer is already running
	current, err := c.GetCurrentTimeEntry()
	if err != nil {
		fmt.Println(formatError(err))
		os.Exit(2)
	}

	// If timer is running, ask user if they want to stop it
	if current != nil {
		elapsed := formatElapsedTime(current.Start)
		fmt.Printf("A timer is currently running: \"%s\" (started %s ago)\n", current.Description, elapsed)

		if !continueYes {
			fmt.Print("Stop this timer and start a new one? [y/N]: ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(formatError(fmt.Errorf("failed to read input: %w", err)))
				os.Exit(2)
			}

			input = strings.TrimSpace(strings.ToLower(input))

			if input != "y" && input != "yes" {
				fmt.Println("Keeping current timer running. No new timer started.")
				os.Exit(0)
			}
		}

		stoppedEntry, err := c.StopTimeEntry(current.ID)
		if err != nil {
			fmt.Println(formatError(fmt.Errorf("failed to stop current timer: %w", err)))
			os.Exit(2)
		}

		duration := formatDuration(stoppedEntry.Duration)
		fmt.Printf("✓ Stopped: \"%s\" (duration: %s)\n", stoppedEntry.Description, duration)
	}

	// Start a new timer with the same description and project
	newEntry, err := c.StartTimeEntry(entry.Description, entry.ProjectID, nil)
	if err != nil {
		fmt.Println(formatError(err))
		os.Exit(2)
	}

	// Get project name for display
	projectName := "No project"
	if entry.ProjectID != nil {
		projects, err := c.GetProjects()
		if err == nil {
			for _, p := range projects {
				if p.ID == *entry.ProjectID {
					projectName = p.Name
					break
				}
			}
		}
	}

	fmt.Printf("✓ Timer started: \"%s\"\n", newEntry.Description)
	fmt.Printf("  Start time: %s\n", newEntry.Start.Local().Format("15:04"))
	fmt.Printf("  Project: %s\n", projectName)
	fmt.Printf("  Entry ID: %s\n", newEntry.ID)
}
