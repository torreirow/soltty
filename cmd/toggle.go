package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	toggleDescription string
	toggleProject     string
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle the running timer (stop if running, continue last if stopped)",
	Long: `Toggle acts as a non-interactive play/pause for bar widgets and scripts.

Behaviour:
  - Timer running  → stop it
  - Timer stopped, last entry available → continue last entry
  - Timer stopped, no entries, --description provided → start new timer
  - Timer stopped, no entries, no --description → exit 1

Examples:
  soltty toggle
  soltty toggle --description "General. Tasks" --project "TMCS-General"`,
	Run: runToggle,
}

func init() {
	toggleCmd.Flags().StringVarP(&toggleDescription, "description", "d", "", "Description for new timer (used only when no previous entry exists)")
	toggleCmd.Flags().StringVarP(&toggleProject, "project", "p", "", "Project for new timer (used only when no previous entry exists)")
}

func runToggle(cmd *cobra.Command, args []string) {
	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		os.Exit(2)
	}

	current, err := c.GetCurrentTimeEntry()
	if err != nil {
		fmt.Println(formatError(err))
		os.Exit(2)
	}

	// Timer running → stop it
	if current != nil {
		stoppedEntry, err := c.StopTimeEntry(current.ID)
		if err != nil {
			fmt.Println(formatError(fmt.Errorf("failed to stop timer: %w", err)))
			os.Exit(2)
		}
		duration := formatDuration(stoppedEntry.Duration)
		fmt.Printf("✓ Stopped: \"%s\" (duration: %s)\n", stoppedEntry.Description, duration)
		return
	}

	// Timer stopped → try to continue last entry
	entries, err := c.ListTimeEntries(1)
	if err != nil {
		fmt.Println(formatError(fmt.Errorf("failed to fetch entries: %w", err)))
		os.Exit(2)
	}

	if len(entries) > 0 {
		last := entries[0]
		newEntry, err := c.StartTimeEntry(last.Description, last.ProjectID, nil)
		if err != nil {
			fmt.Println(formatError(err))
			os.Exit(2)
		}
		fmt.Printf("✓ Timer started: \"%s\"\n", newEntry.Description)
		fmt.Printf("  Start time: %s\n", newEntry.Start.Local().Format("15:04"))
		return
	}

	// No entries → start fresh if description provided
	if toggleDescription != "" {
		var projectID *string
		if toggleProject != "" {
			pid, err := c.FindProjectByName(toggleProject)
			if err != nil {
				fmt.Println(formatError(err))
				os.Exit(2)
			}
			projectID = pid
		}

		newEntry, err := c.StartTimeEntry(toggleDescription, projectID, nil)
		if err != nil {
			fmt.Println(formatError(err))
			os.Exit(2)
		}
		fmt.Printf("✓ Timer started: \"%s\"\n", newEntry.Description)
		fmt.Printf("  Start time: %s\n", newEntry.Start.Local().Format("15:04"))
		if toggleProject != "" {
			fmt.Printf("  Project: %s\n", toggleProject)
		}
		return
	}

	fmt.Println("No timer running and no previous entry to continue.")
	os.Exit(1)
}
