package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	addStart   string
	addEnd     string
	addProject string
)

var addCmd = &cobra.Command{
	Use:   "add <description>",
	Short: "Add a completed time entry",
	Long: `Create a time entry with specific start and end times.

Time formats:
  2026-09-16T14:00:00Z    date, 'T', time and timezone offset
  2026-09-16 14:00        date and time, local timezone
  14:00                   time only, today, local timezone

The date must be YYYY-MM-DD. The separator may be 'T' or a space, seconds are
optional, and a timezone offset is only allowed after a 'T' separator.

An --end without a date belongs to the day of --start, not to today. An entry
crossing midnight therefore needs --end with an explicit date.

Examples:
  soltty add "Meeting" --start "2026-03-31T14:00:00Z" --end "2026-03-31T15:30:00Z"
  soltty add "Code review" --start "14:00" --end "15:30"
  soltty add "Client call" --start "2026-09-16 09:00" --end "17:00"
  soltty add "Deploy" --start "2026-09-16 23:00" --end "2026-09-17 01:00"
  soltty add "Sprint planning" --start "10:00" --end "12:00" --project "Acme-Meetings"`,
	Args: cobra.ExactArgs(1),
	Run:  runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addStart, "start", "", "Start time (see 'soltty add --help' for formats) [required]")
	addCmd.Flags().StringVar(&addEnd, "end", "", "End time (see 'soltty add --help' for formats) [required]")
	addCmd.Flags().StringVarP(&addProject, "project", "p", "", "Project name")
	addCmd.MarkFlagRequired("start")
	addCmd.MarkFlagRequired("end")
}

// resolveEntryTimes resolves the start and end of a new entry. An end time
// without its own date belongs to the day of the start time, not to today. The
// date is never shifted forward: an end that still falls on or before the start
// is left as-is so the caller can reject it.
func resolveEntryTimes(start, end parsedTime) (time.Time, time.Time) {
	if end.HasDate {
		return start.Time, end.Time
	}
	return start.Time, stampOnDate(start.Time, end.Time)
}

func runAdd(cmd *cobra.Command, args []string) {
	description := args[0]

	// Parse times
	start, err := parseTime(addStart)
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	end, err := parseTime(addEnd)
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	startTime, endTime := resolveEntryTimes(start, end)

	// Validate: end must be after start
	if !endTime.After(startTime) {
		fmt.Println("Error: End time must be after start time")
		if !end.HasDate {
			fmt.Printf("  Does the entry cross midnight? Give --end with a date, e.g. --end \"%s\".\n",
				startTime.Local().AddDate(0, 0, 1).Format("2006-01-02")+endTime.Local().Format(" 15:04"))
		}
		return
	}

	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	// Resolve project if specified
	var projectID *string
	if addProject != "" {
		pid, err := c.FindProjectByName(addProject)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		projectID = pid
	}

	// Create the entry
	entry, err := c.CreateTimeEntry(description, startTime, endTime, projectID)
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	duration := formatDuration(entry.Duration)
	fmt.Printf("✓ Time entry added: \"%s\"\n", entry.Description)
	fmt.Printf("  Start: %s\n", formatEntryTime(entry.Start))
	fmt.Printf("  End: %s\n", formatEntryTime(*entry.End))
	fmt.Printf("  Duration: %s\n", duration)
	if addProject != "" {
		fmt.Printf("  Project: %s\n", addProject)
	}
	fmt.Printf("  Entry ID: %s\n", entry.ID)
}
