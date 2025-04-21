package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func main() {
	// Get API credentials from environment variables
	apiKey := os.Getenv("HUNTRESS_API_KEY")
	apiSecret := os.Getenv("HUNTRESS_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Missing required environment variables: HUNTRESS_API_KEY and/or HUNTRESS_API_SECRET")
	}

	// Create a new Huntress client with options
	client := huntress.New(
		huntress.WithCredentials(apiKey, apiSecret),
		huntress.WithTimeout(30*time.Second),
		huntress.WithRetryConfig(3, 500*time.Millisecond, 5*time.Second),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get the current account
	fmt.Println("Fetching current account...")
	account, err := client.Account.GetCurrent(ctx)
	if err != nil {
		log.Fatalf("Error fetching account: %v", err)
	}

	fmt.Printf("Account: %s (ID: %d)\n", account.Name, account.ID)

	// List organizations for the account
	fmt.Println("\nFetching organizations...")
	orgs, pagination, err := client.Organization.List(ctx, &huntress.OrganizationListOptions{
		ListOptions: huntress.ListOptions{
			Page:     1,
			PerPage:  10,
			SortBy:   "name",
			SortDesc: false,
		},
	})
	if err != nil {
		log.Fatalf("Error listing organizations: %v", err)
	}

	fmt.Printf("Found %d organizations (page %d of %d, total items: %d):\n",
		len(orgs), pagination.CurrentPage, pagination.TotalPages, pagination.TotalItems)

	for i, org := range orgs {
		fmt.Printf("%d. %s (ID: %d, Status: %s)\n", i+1, org.Name, org.ID, org.Status)

		// For the first organization, get its agents
		if i == 0 && len(orgs) > 0 {
			fmt.Printf("\nFetching agents for organization '%s'...\n", org.Name)
			agents, _, err := client.Agent.List(ctx, &huntress.AgentListOptions{
				OrganizationID: org.ID,
				ListOptions: huntress.ListOptions{
					Page:    1,
					PerPage: 5,
				},
			})
			if err != nil {
				fmt.Printf("Error listing agents: %v\n", err)
				continue
			}

			fmt.Printf("Found %d agents:\n", len(agents))
			for j, agent := range agents {
				fmt.Printf("  %d. %s (ID: %s, OS: %s, Status: %s)\n",
					j+1, agent.Hostname, agent.ID, agent.OS, agent.Status)
			}

			fmt.Printf("\nFetching incidents for organization '%s'...\n", org.Name)
			incidents, _, err := client.Incident.List(ctx, &huntress.IncidentListOptions{
				Organization: org.ID,
				ListOptions: huntress.ListOptions{
					Page:    1,
					PerPage: 5,
				},
			})
			if err != nil {
				fmt.Printf("Error listing incidents: %v\n", err)
				continue
			}

			fmt.Printf("Found %d incidents:\n", len(incidents))
			for j, incident := range incidents {
				fmt.Printf("  %d. %s (ID: %s, Status: %s, Severity: %s)\n",
					j+1, incident.Title, incident.ID, incident.Status, incident.Severity)
			}
		}
	}

	// Get account statistics
	fmt.Println("\nFetching account statistics...")
	stats, err := client.Account.GetStatistics(ctx, account.ID)
	if err != nil {
		log.Fatalf("Error fetching account statistics: %v", err)
	}

	fmt.Printf("Account statistics:\n")
	fmt.Printf("  Organization Count: %d\n", stats.OrganizationCount)
	fmt.Printf("  Active Agent Count: %d\n", stats.ActiveAgentCount)
	fmt.Printf("  Total Agent Count: %d\n", stats.TotalAgentCount)
	fmt.Printf("  Open Incident Count: %d\n", stats.OpenIncidentCount)
	fmt.Printf("  Total Incident Count: %d\n", stats.TotalIncidentCount)

	fmt.Println("\nExample completed successfully")
}
