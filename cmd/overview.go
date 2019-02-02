package cmd

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
	"github.com/spf13/cobra"
)

// OverviewCmd is exported
var OverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show pending merge requests overview",
	Long:  "This displays a table by project and approver sorted by the most pending merge requests",
	Run: func(cmd *cobra.Command, args []string) {
		client := gitlab.NewClient(nil)

		pendingRequests, err := client.PendingRequests()
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}
		rows := make(map[string]map[string]int)
		headers := make(map[string]int)
		for _, pr := range pendingRequests {
			project := pr.Project.Name
			headers[project]++
			for _, username := range pr.ApproverNames() {
				projectMap, ok := rows[username]
				if !ok {
					projectMap = make(map[string]int)
				}
				projectMap[project]++
				rows[username] = projectMap
			}
		}

		fmt.Println(rows)
	},
}
