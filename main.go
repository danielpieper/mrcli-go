package main

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
)

func main() {
	client := gitlab.NewClient(nil)
	mergeRequests := client.GetMergeRequests()
	for _, mergeRequest := range mergeRequests {
		fmt.Println(mergeRequest.Title)
		approvals := client.GetMergeRequestApprovals(mergeRequest)
		fmt.Println(approvals.State)
	}
}
