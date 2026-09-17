// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Wouter van der Toorren

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	startProject string
	startTime    string
	startYes     bool
)

var startCmd = &cobra.Command{
	Use:   "start <description>",
	Short: "Start a new timer",
	Long: `Start a new time tracking entry.

Examples:
  soltty start "Working on feature X"
  soltty start "Bug fix" --project "Example-Project"
  soltty start "Forgot to start" --time "09:00"
  soltty start "Task" --time "2026-03-31T08:00:00Z"
  soltty start "Backdated" --time "2026-09-16 09:00"

Time formats: 2026-09-16T14:00:00Z, 2026-09-16 14:00 or 14:00. The date must be
YYYY-MM-DD. The separator may be 'T' or a space, seconds are optional, and a
timezone offset is only allowed after a 'T' separator. A bare time means today,
local timezone.`,
	Args: cobra.ExactArgs(1),
	Run:  runStart,
}

func init() {
	startCmd.Flags().StringVarP(&startProject, "project", "p", "", "Project name")
	startCmd.Flags().StringVarP(&startTime, "time", "t", "", "Custom start time (see 'soltty start --help' for formats)")
	startCmd.Flags().BoolVarP(&startYes, "yes", "y", false, "Skip confirmation when a timer is already running")
}

func runStart(cmd *cobra.Command, args []string) {
	description := args[0]

	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
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

		if !startYes {
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

	// Resolve project if specified
	var projectID *string
	if startProject != "" {
		pid, err := c.FindProjectByName(startProject)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		projectID = pid
	}

	// Parse custom start time if specified
	var customStart *time.Time
	if startTime != "" {
		pt, err := parseTime(startTime)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		customStart = &pt.Time
	}

	// Start the timer
	entry, err := c.StartTimeEntry(description, projectID, customStart)
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	fmt.Printf("✓ Timer started: \"%s\"\n", entry.Description)
	if customStart != nil {
		fmt.Printf("  Start time: %s (custom)\n", entry.Start.Local().Format("15:04"))
	} else {
		fmt.Printf("  Start time: %s\n", entry.Start.Local().Format("15:04"))
	}
	if startProject != "" {
		fmt.Printf("  Project: %s\n", startProject)
	}
	fmt.Printf("  Entry ID: %s\n", entry.ID)
}
