package main

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
	g "github.com/xanzy/go-gitlab"
	"strings"
	"time"
)

func main() {
	client := gitlab.NewClient(nil)

	pendingRequests, err := client.PendingRequests()
	if err != nil {
		fmt.Println("An error occured:", err)
		return
	}
	// user, err := client.AuthenticatedUser()
	// if err != nil {
	// 	fmt.Println("An error occured:", err)
	// 	return
	// }

	// list(pendingRequests, user)
	overview(pendingRequests)
}

func list(pendingRequests []*gitlab.PendingRequest, user *g.User) {
	fmt.Printf("%d Pending merge requests for %v:\n\n", len(pendingRequests), user.Username)
	for _, pr := range pendingRequests {
		fmt.Println(pr.Request.Author.Username)
		fmt.Printf("[%v] %v\n", pr.Project.Name, pr.Request.Title)
		fmt.Println(pr.Request.WebURL)
		fmt.Println("Created:", time.Since(*pr.Request.CreatedAt).Round(time.Duration(time.Hour)))
		fmt.Println("Updated:", time.Since(*pr.Request.UpdatedAt).Round(time.Duration(time.Hour)))
		fmt.Println("Approvers:", strings.Join(pr.ApproverNames(), ", "))
		fmt.Println()
	}
}

func overview(pendingRequests []*gitlab.PendingRequest) {
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
}
