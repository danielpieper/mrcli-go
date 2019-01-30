package main

import (
	"github.com/danielpieper/mrcli-go/gitlab"
)

func main() {
	client := gitlab.NewClient(nil)
	client.GetMergeRequests()
}

// func getApprovals(client *gitlab.Client, mergeRequest *gitlab.MergeRequest) {
// 	approvals, _, err := client.MergeRequests.GetMergeRequestApprovals(mergeRequest.ProjectID, mergeRequest.IID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(approvals)
// }
