package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/torreirow/soltty/internal/client"
)

var (
	listProjectsClientFilter string
)

var listProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List all projects",
	Long:  `Display all active projects with their client names.`,
	Run:   runListProjects,
}

type listProjectJSON struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ClientID *string `json:"client_id"`
	Client   *string `json:"client"`
}

func init() {
	listProjectsCmd.Flags().StringVarP(&listProjectsClientFilter, "client", "c", "", "Filter projects by client name (partial match)")
}

func runListProjects(cmd *cobra.Command, args []string) {
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

	clientMap := make(map[string]string)
	for _, cl := range clients {
		clientMap[cl.ID] = cl.Name
	}

	projects, err := c.GetProjects()
	if err != nil {
		fmt.Println(formatError(fmt.Errorf("failed to fetch projects: %w", err)))
		return
	}

	var activeProjects []client.Project
	for _, p := range projects {
		if !p.IsArchived {
			activeProjects = append(activeProjects, p)
		}
	}

	if listProjectsClientFilter != "" {
		var filteredProjects []client.Project
		filterLower := strings.ToLower(listProjectsClientFilter)

		for _, p := range activeProjects {
			if p.ClientID != nil {
				clientName := clientMap[*p.ClientID]
				if strings.Contains(strings.ToLower(clientName), filterLower) {
					filteredProjects = append(filteredProjects, p)
				}
			}
		}

		activeProjects = filteredProjects

		if len(activeProjects) == 0 && !listJSON {
			fmt.Printf("No projects found for client: %s\n", listProjectsClientFilter)
			return
		}
	}

	sort.Slice(activeProjects, func(i, j int) bool {
		clientNameI := "(no client)"
		clientNameJ := "(no client)"

		if activeProjects[i].ClientID != nil {
			if name, ok := clientMap[*activeProjects[i].ClientID]; ok {
				clientNameI = name
			} else {
				clientNameI = "(unknown client)"
			}
		}

		if activeProjects[j].ClientID != nil {
			if name, ok := clientMap[*activeProjects[j].ClientID]; ok {
				clientNameJ = name
			} else {
				clientNameJ = "(unknown client)"
			}
		}

		if clientNameI != clientNameJ {
			return clientNameI < clientNameJ
		}
		return activeProjects[i].Name < activeProjects[j].Name
	})

	if listJSON {
		result := make([]listProjectJSON, 0, len(activeProjects))
		for _, p := range activeProjects {
			entry := listProjectJSON{
				ID:       p.ID,
				Name:     p.Name,
				ClientID: p.ClientID,
			}
			if p.ClientID != nil {
				if name, ok := clientMap[*p.ClientID]; ok {
					entry.Client = &name
				}
			}
			result = append(result, entry)
		}
		b, err := json.Marshal(result)
		if err != nil {
			fmt.Println(formatError(err))
			return
		}
		fmt.Println(string(b))
		return
	}

	if len(activeProjects) == 0 {
		fmt.Println("No projects found")
		return
	}

	fmt.Println("Client            | Project")
	fmt.Println(strings.Repeat("-", 18) + "|" + strings.Repeat("-", 30))

	for _, p := range activeProjects {
		clientName := "(no client)"
		if p.ClientID != nil {
			if name, ok := clientMap[*p.ClientID]; ok {
				clientName = name
			} else {
				clientName = "(unknown client)"
			}
		}

		fmt.Printf("%-17s | %s\n", clientName, p.Name)
	}
}
