package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var currentJSON bool

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current running timer",
	Long: `Display information about the currently running time entry.

Example:
  soltty current
  soltty current --json`,
	Run: runCurrent,
}

func init() {
	currentCmd.Flags().BoolVar(&currentJSON, "json", false, "Output as JSON")
}

type currentJSONOutput struct {
	Running     bool    `json:"running"`
	ID          string  `json:"id,omitempty"`
	Description string  `json:"description,omitempty"`
	ProjectID   *string `json:"project_id,omitempty"`
	Project     *string `json:"project,omitempty"`
	Elapsed     string  `json:"elapsed,omitempty"`
	Start       string  `json:"start,omitempty"`
}

func runCurrent(cmd *cobra.Command, args []string) {
	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	current, err := c.GetCurrentTimeEntry()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	if currentJSON {
		if current == nil {
			fmt.Println(`{"running":false}`)
			return
		}

		out := currentJSONOutput{
			Running:     true,
			ID:          current.ID,
			Description: current.Description,
			ProjectID:   current.ProjectID,
			Elapsed:     formatElapsedTime(current.Start),
			Start:       current.Start.UTC().Format("2006-01-02T15:04:05Z"),
		}

		if current.ProjectID != nil {
			projects, err := c.GetProjects()
			if err == nil {
				for _, p := range projects {
					if p.ID == *current.ProjectID {
						name := p.Name
						out.Project = &name
						break
					}
				}
			}
		}

		b, err := json.Marshal(out)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		fmt.Println(string(b))
		return
	}

	if current == nil {
		fmt.Println("No timer is currently running")
		return
	}

	elapsed := formatElapsedTime(current.Start)
	fmt.Printf("Timer running: \"%s\"\n", current.Description)
	fmt.Printf("  Started: %s\n", current.Start.Local().Format("15:04"))
	fmt.Printf("  Elapsed: %s\n", elapsed)
	if current.ProjectID != nil {
		fmt.Printf("  Project ID: %s\n", *current.ProjectID)
	}
	fmt.Printf("  Entry ID: %s\n", current.ID)
}
