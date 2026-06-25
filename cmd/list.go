package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listShowID bool
	listJSON   bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List time entries, clients, or projects",
	Long: `Display time entries, clients, or projects.

Subcommands:
  soltty list              # List recent time entries (default)
  soltty list clients      # List all clients with project counts
  soltty list projects     # List all projects with client names
  soltty list projects -c <client>  # Filter projects by client

Examples:
  soltty list
  soltty list --limit 5
  soltty list --id           # Show entry IDs for deletion
  soltty list --json         # Machine-readable JSON output
  soltty list clients
  soltty list clients --json
  soltty list projects
  soltty list projects --json
  soltty list projects -c Acme`,
	Run: runList,
}

type listEntryJSON struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	ProjectID   *string `json:"project_id"`
	Project     *string `json:"project"`
	Date        string  `json:"date"`
	Start       string  `json:"start"`
	End         *string `json:"end"`
	Duration    int     `json:"duration"`
	Running     bool    `json:"running"`
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 10, "Number of entries to show")
	listCmd.Flags().BoolVar(&listShowID, "id", false, "Show entry IDs")
	listCmd.PersistentFlags().BoolVar(&listJSON, "json", false, "Output as JSON")

	// Register subcommands
	listCmd.AddCommand(listClientsCmd)
	listCmd.AddCommand(listProjectsCmd)
}

func runList(cmd *cobra.Command, args []string) {
	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	entries, err := c.ListTimeEntries(listLimit)
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	// Fetch projects for name lookup
	projects, err := c.GetProjects()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}
	projectMap := make(map[string]string)
	for _, p := range projects {
		projectMap[p.ID] = p.Name
	}

	if listJSON {
		result := make([]listEntryJSON, 0, len(entries))
		for _, entry := range entries {
			e := listEntryJSON{
				ID:          entry.ID,
				Description: entry.Description,
				ProjectID:   entry.ProjectID,
				Date:        entry.Start.Local().Format("2006-01-02"),
				Start:       entry.Start.UTC().Format("2006-01-02T15:04:05Z"),
				Duration:    entry.Duration,
				Running:     entry.End == nil,
			}
			if entry.ProjectID != nil {
				if name, ok := projectMap[*entry.ProjectID]; ok {
					e.Project = &name
				}
			}
			if entry.End != nil {
				s := entry.End.UTC().Format("2006-01-02T15:04:05Z")
				e.End = &s
			}
			result = append(result, e)
		}
		b, err := json.Marshal(result)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		fmt.Println(string(b))
		return
	}

	if len(entries) == 0 {
		fmt.Println("No time entries found")
		return
	}

	// Print header
	if listShowID {
		fmt.Println("ID                                   | Date       | Start | Duration | Project        | Description")
		fmt.Println(strings.Repeat("-", 120))
	} else {
		fmt.Println("ID       | Date       | Start | Duration | Project        | Description")
		fmt.Println(strings.Repeat("-", 90))
	}

	// Print entries
	for _, entry := range entries {
		localStart := entry.Start.Local()
		date := localStart.Format("2006-01-02")
		startTime := localStart.Format("15:04")

		var duration string
		if entry.End == nil {
			duration = "running"
		} else {
			duration = formatDuration(entry.Duration)
		}

		projectName := "No project"
		if entry.ProjectID != nil {
			if name, ok := projectMap[*entry.ProjectID]; ok {
				projectName = name
			}
		}

		if listShowID {
			fmt.Printf("%-36s | %-10s | %-5s | %-8s | %-14s | %s\n",
				entry.ID, date, startTime, duration, projectName, entry.Description)
		} else {
			fmt.Printf("%-8s | %-10s | %-5s | %-8s | %-14s | %s\n",
				entry.ID[:8], date, startTime, duration, projectName, entry.Description)
		}
	}
}
