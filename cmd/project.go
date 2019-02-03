package cmd

import (
	"fmt"
	"github.com/bclicn/color"
	"github.com/danielpieper/mrcli-go/gitlab"
	"github.com/spf13/cobra"
	"strings"
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
			fmt.Println(color.LightGray(pr.Request.Author.Username))
			fmt.Printf("[%v] %v\n", pr.Project.Name, pr.Color(pr.Request.Title))
			fmt.Println(pr.Request.WebURL)
			fmt.Println(color.LightGray("Created:"), pr.HumanReadableCreatedAtDiff())
			fmt.Println(color.LightGray("Updated:"), pr.HumanReadableUpdatedAtDiff())
			fmt.Println(color.LightGray("Approvers:"), strings.Join(pr.ApproverNames(), ", "))
			fmt.Println()
		}
	},
}
