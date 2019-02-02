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

		rankedProjects := rankProjectsByPendingRequests(pendingRequests)
		headers := []string{"Approver"}
		for _, h := range rankedProjects {
			headers = append(headers, h.name)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)

		approverMap := createApproverMap(pendingRequests)
		for username, projectMap := range approverMap {
			row := []string{username}
			for _, h := range rankedProjects {
				pendingRequestCount, ok := projectMap[h.name]
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
			projectMap["Total"]++
			approverMap[username] = projectMap
		}
	}
	return approverMap
}

func rankProjectsByPendingRequests(pendingRequests []*gitlab.PendingRequest) (projects projectRanking) {
	for _, pr := range pendingRequests {
		project := pr.Project.Name
		projects.increment("Total", 1)
		projects.increment(project, 1)
	}
	sort.Sort(sort.Reverse(projects))
	return
}

type projectRank struct {
	name  string
	value int
}

type projectRanking []projectRank

func (p projectRanking) Len() int           { return len(p) }
func (p projectRanking) Less(i, j int) bool { return p[i].value < p[j].value }
func (p projectRanking) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p *projectRanking) increment(name string, value int) {
	v, ok := p.get(name)
	if ok {
		v.value += value
		return
	}
	v = &projectRank{name, value}
	*p = append(*p, *v)
}

func (p projectRanking) get(key string) (*projectRank, bool) {
	for _, v := range p {
		if v.name == key {
			return &v, true
		}
	}
	return &projectRank{}, false
}
