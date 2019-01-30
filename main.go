package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"time"
)

func main() {
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	gitlabClient := gitlab.NewClient(nil, gitlabToken)
	baseUrl, ok := os.LookupEnv("GITLAB_URL")
	if ok {
		gitlabClient.SetBaseURL(baseUrl)
	}

	lastMonth := time.Now().AddDate(0, -1, 0)
	mergeRequestOptions := &gitlab.ListMergeRequestsOptions{
		State:        gitlab.String("opened"),
		OrderBy:      gitlab.String("created_at"),
		CreatedAfter: &lastMonth,
	}

	mergeRequests, _, err := gitlabClient.MergeRequests.ListMergeRequests(mergeRequestOptions)
	if err != nil {
		log.Fatal(err)
	}
	for _, mergeRequest := range mergeRequests {
		fmt.Println(mergeRequest.Title)
	}
}

func getApprovals(client *gitlab.Client, mergeRequest *gitlab.MergeRequest) {
	approvals, _, err := client.MergeRequests.GetMergeRequestApprovals(mergeRequest.ProjectID, mergeRequest.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(approvals.State)
}
