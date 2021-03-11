package commands

import (
	"fmt"
	"github.com/bclicn/color"
	"github.com/danielpieper/mrcli-go/internal/gitlab"
	"github.com/spf13/cobra"
	"strings"
)

// ApproverCmd is exported
var ApproverCmd = &cobra.Command{
	Use:     "approver",
	Aliases: []string{"a"},
	Short:   "Show pending merge requests by user",
	Long:    "This displays a list of pending merge requests by user",
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, userNames []string) {
		client := gitlab.NewClient(nil)

		pendingRequests, err := client.PendingRequests()
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}

		if len(userNames) == 0 {
			user, err := client.AuthenticatedUser()
			if err != nil {
				fmt.Println("An error occured:", err)
				return
			}
			userNames = append(userNames, user.Username)
		}

		filteredPendingRequests := []gitlab.PendingRequest{}
		for _, pr := range pendingRequests {
		filterUser:
			for _, approverName := range pr.ApproverNames() {
				for _, userName := range userNames {
					if approverName == userName {
						filteredPendingRequests = append(filteredPendingRequests, *pr)
						break filterUser
					}
				}
			}
		}

		fmt.Printf("%d Pending merge requests for approvers %v:\n\n", len(filteredPendingRequests), userNames)
		for _, pr := range filteredPendingRequests {
			fmt.Println(color.LightGray(pr.Request.Author.Username))
			fmt.Printf("[%v] %v\n", pr.Project.Name, pr.Color(pr.Request.Title))
			fmt.Println(pr.Request.WebURL)
			fmt.Println(color.LightGray(fmt.Sprintf("%10s", "Created:")), pr.HumanReadableCreatedAtDiff())
			fmt.Println(color.LightGray(fmt.Sprintf("%10s", "Updated:")), pr.HumanReadableUpdatedAtDiff())
			fmt.Println(color.LightGray(fmt.Sprintf("%10s", "Approvers:")), strings.Join(pr.ApproverNames(), ", "))
			fmt.Println()
		}
	},
}
