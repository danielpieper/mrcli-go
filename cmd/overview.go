package cmd

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strconv"
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

		projectMap := createProjectMap(pendingRequests)
		rankedHeaders := rankByMergeRequestCount(projectMap)
		headers := []string{"Approver"}
		for _, h := range rankedHeaders {
			headers = append(headers, h.Key)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)

		approverMap := createApproverMap(pendingRequests)
		for username, projectMap := range approverMap {
			row := []string{username}
			for _, h := range rankedHeaders {
				pendingRequestCount, ok := projectMap[h.Key]
				value := "0"
				if ok {
					value = strconv.Itoa(pendingRequestCount)
				}
				row = append(row, value)
			}
			table.Append(row)
		}

		table.Render()
	},
}

func createApproverMap(pendingRequests []*gitlab.PendingRequest) map[string]map[string]int {
	approverMap := make(map[string]map[string]int)
	for _, pr := range pendingRequests {
		project := pr.Project.Name
		for _, username := range pr.ApproverNames() {
			projectMap, ok := approverMap[username]
			if !ok {
				projectMap = make(map[string]int)
			}
			projectMap[project]++
			approverMap[username] = projectMap
		}
	}
	return approverMap
}

func createProjectMap(pendingRequests []*gitlab.PendingRequest) map[string]int {
	projectMap := map[string]int{}
	for _, pr := range pendingRequests {
		project := pr.Project.Name
		_, ok := projectMap[project]
		if !ok {
			projectMap[project] = 0
		}
		projectMap[project]++
	}
	return projectMap
}

func rankByMergeRequestCount(data map[string]int) PairList {
	pl := make(PairList, len(data))
	i := 0
	for k, v := range data {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// Pair is exported
type Pair struct {
	Key   string
	Value int
}

// PairList is exported
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
