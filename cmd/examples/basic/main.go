// Package main provides a basic example of using the Bishoujo-Huntress Go client to interact with the Huntress API.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func main() {
	apiKey := os.Getenv("HUNTRESS_API_KEY")
	apiSecret := os.Getenv("HUNTRESS_API_SECRET")
	baseURL := os.Getenv("HUNTRESS_BASE_URL") // Optional, defaults to production

	if apiKey == "" || apiSecret == "" {
		log.Fatal("HUNTRESS_API_KEY and HUNTRESS_API_SECRET environment variables must be set")
	}

	// Use the public client constructor from pkg/huntress
	client := huntress.New(
		huntress.WithCredentials(apiKey, apiSecret),
		huntress.WithTimeout(60*time.Second),
	)
	if baseURL != "" {
		client = huntress.New(
			huntress.WithCredentials(apiKey, apiSecret),
			huntress.WithTimeout(60*time.Second),
			huntress.WithBaseURL(baseURL),
		)
	}

	ctx := context.Background()

	// Get current account details
	fmt.Println("Fetching account details...")
	account, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("Error fetching account details: %v", err)
	}
	fmt.Printf("Account Name: %s (ID: %s)\n", account.Name, account.ID)

	// List organizations
	fmt.Println("\nFetching organizations...")
	orgs, pagination, err := client.Organization.List(ctx, &huntress.OrganizationListOptions{
		ListParams: huntress.ListParams{
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
		fmt.Printf("%d. %s (ID: %s, Status: %s)\n", i+1, org.Name, org.ID, org.Status)

		// For the first organization, get its agents
		if i == 0 && len(orgs) > 0 {
			fmt.Printf("\nFetching agents for organization '%s'...\n", org.Name)

			orgIDInt, err := strconv.Atoi(org.ID)
			if err != nil {
				fmt.Printf("Error converting organization ID to int: %v\n", err)
				continue
			}

			agents, _, err := client.Agent.List(ctx, &huntress.AgentListOptions{
				ListParams: huntress.ListParams{
					Page:    1,
					PerPage: 5,
				},
				OrganizationID: orgIDInt,
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
				ListOptions: huntress.ListOptions{
					Page:    1,
					PerPage: 5,
				},
				Organization: orgIDInt,
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
	stats, err := client.Account.GetStats(ctx)
	if err != nil {
		log.Fatalf("Error fetching account statistics: %v", err)
	}

	fmt.Printf("Account statistics:\n")
	fmt.Printf("  Organization Count: %d\n", stats.OrganizationCount)
	fmt.Printf("  Agent Count: %d\n", stats.AgentCount)
	fmt.Printf("  Incident Count: %d\n", stats.IncidentCount)
	fmt.Printf("  User Count: %d\n", stats.UserCount)
}
