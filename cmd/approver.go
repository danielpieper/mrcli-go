package cmd

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// ApproverCmd is exported
var ApproverCmd = &cobra.Command{
	Use:   "approver",
	Short: "Show pending merge requests by user",
	Long:  "This displays a list of pending merge requests by user",
	Run: func(cmd *cobra.Command, args []string) {
		client := gitlab.NewClient(nil)

		pendingRequests, err := client.PendingRequests()
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}
		user, err := client.AuthenticatedUser()
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}

		filteredPendingRequests := []gitlab.PendingRequest{}
		for _, pr := range pendingRequests {
			for _, username := range pr.ApproverNames() {
				if username == user.Username {
					filteredPendingRequests = append(filteredPendingRequests, *pr)
					break
				}
			}
		}

		fmt.Printf("%d Pending merge requests for %v:\n\n", len(filteredPendingRequests), user.Username)
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
