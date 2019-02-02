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
	Use:     "overview",
	Aliases: []string{"o"},
	Short:   "Show pending merge requests overview",
	Long:    "This displays a table by project and approver sorted by the most pending merge requests",
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

		rankedApprovers := rankApproversByPendingRequests(pendingRequests)
		for _, a := range rankedApprovers {
			row := []string{a.username}
			for _, h := range rankedProjects {
				project, ok := a.projects.get(h.name)
				value := ""
				if ok {
					value = strconv.Itoa(project.value)
				}
				row = append(row, value)
			}
			table.Append(row)
		}

		table.Render()
	},
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

type projectRanking []*projectRank

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
	*p = append(*p, v)
}

func (p *projectRanking) get(name string) (*projectRank, bool) {
	for _, v := range *p {
		if v.name == name {
			return v, true
		}
	}
	return &projectRank{}, false
}

func rankApproversByPendingRequests(pendingRequests []*gitlab.PendingRequest) (approvers approverRanking) {
	for _, pr := range pendingRequests {
		project := pr.Project.Name
		for _, username := range pr.ApproverNames() {
			approvers.increment(username, "Total", 1)
			approvers.increment(username, project, 1)
		}
	}
	sort.Sort(sort.Reverse(approvers))
	return
}

type approverRank struct {
	username string
	projects *projectRanking
}

type approverRanking []*approverRank

func (a approverRanking) Len() int { return len(a) }
func (a approverRanking) Less(i, j int) bool {
	itemA, _ := a[i].projects.get("Total")
	itemB, _ := a[j].projects.get("Total")
	return itemA.value < itemB.value
}
func (a approverRanking) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a *approverRanking) increment(username string, project string, value int) {
	v, ok := a.get(username)
	if ok {
		v.projects.increment(project, value)
		return
	}
	v = &approverRank{username, &projectRanking{&projectRank{project, value}}}
	*a = append(*a, v)
}

func (a *approverRanking) get(username string) (*approverRank, bool) {
	for _, v := range *a {
		if v.username == username {
			return v, true
		}
	}
	return &approverRank{}, false
}
