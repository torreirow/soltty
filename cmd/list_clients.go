package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/torreirow/soltty/internal/client"
)

var listClientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List all clients",
	Long:  `Display all active clients with their project counts.`,
	Run:   runListClients,
}

type listClientJSON struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ProjectCount int    `json:"project_count"`
}

func runListClients(cmd *cobra.Command, args []string) {
	c, err := getClient()
	if err != nil {
		fmt.Println(formatError(err))
		return
	}

	clients, err := c.GetClients()
	if err != nil {
		fmt.Println(formatError(fmt.Errorf("failed to fetch clients: %w", err)))
		return
	}

	var activeClients []client.SolidtimeClient
	for _, cl := range clients {
		if !cl.IsArchived {
			activeClients = append(activeClients, cl)
		}
	}

	projects, err := c.GetProjects()
	if err != nil {
		fmt.Println(formatError(fmt.Errorf("failed to fetch projects: %w", err)))
		return
	}

	projectCounts := make(map[string]int)
	for _, p := range projects {
		if !p.IsArchived && p.ClientID != nil {
			projectCounts[*p.ClientID]++
		}
	}

	sort.Slice(activeClients, func(i, j int) bool {
		return activeClients[i].Name < activeClients[j].Name
	})

	if listJSON {
		result := make([]listClientJSON, 0, len(activeClients))
		for _, cl := range activeClients {
			result = append(result, listClientJSON{
				ID:           cl.ID,
				Name:         cl.Name,
				ProjectCount: projectCounts[cl.ID],
			})
		}
		b, err := json.Marshal(result)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		fmt.Println(string(b))
		return
	}

	if len(activeClients) == 0 {
		fmt.Println("No clients found")
		return
	}

	for _, cl := range activeClients {
		count := projectCounts[cl.ID]
		if count == 1 {
			fmt.Printf("%s (1 project)\n", cl.Name)
		} else {
			fmt.Printf("%s (%d projects)\n", cl.Name, count)
		}
	}
}
