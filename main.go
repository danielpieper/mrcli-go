package main

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
)

func main() {
	client := gitlab.NewClient(nil)
	pendingRequests := client.GetPendingRequests()
	for _, pr := range pendingRequests {
		fmt.Println(pr.Request.Title)
		// approvals := client.GetMergeRequestApprovals(mergeRequest)
		// fmt.Println(approvals.State)
	}
}
