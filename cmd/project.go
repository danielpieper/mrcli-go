package cmd

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// ProjectCmd is exported
var ProjectCmd = &cobra.Command{
	Use:     "project",
	Aliases: []string{"p"},
	Short:   "Show pending merge requests by project",
	Long:    "This displays a list of pending merge requests by project",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, projectNames []string) {
		client := gitlab.NewClient(nil)

		pendingRequests, err := client.PendingRequests()
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}

		filteredPendingRequests := []gitlab.PendingRequest{}
		for _, pr := range pendingRequests {
			for _, p := range projectNames {
				if p == pr.Project.Name {
					filteredPendingRequests = append(filteredPendingRequests, *pr)
					break
				}
			}
		}

		fmt.Printf("%d Pending merge requests for projects %v:\n\n", len(filteredPendingRequests), projectNames)
		for _, pr := range filteredPendingRequests {
			fmt.Println(pr.Request.Author.Username)
			fmt.Printf("[%v] %v\n", pr.Project.Name, pr.Request.Title)
			fmt.Println(pr.Request.WebURL)
			fmt.Println("Created:", time.Since(*pr.Request.CreatedAt).Round(time.Duration(time.Hour)))
			fmt.Println("Updated:", time.Since(*pr.Request.UpdatedAt).Round(time.Duration(time.Hour)))
			fmt.Println("Approvers:", strings.Join(pr.ApproverNames(), ", "))
			fmt.Println()
		}
	},
}
