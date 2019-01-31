package gitlab

import (
	g "github.com/xanzy/go-gitlab"
	"log"
	"net/http"
	"os"
	"time"
)

// Client is exported
type Client struct {
	gitlabClient *g.Client
}

// NewClient is exported
func NewClient(httpClient *http.Client) *Client {
	token := os.Getenv("GITLAB_TOKEN")
	gClient := g.NewClient(httpClient, token)
	baseURL, ok := os.LookupEnv("GITLAB_URL")
	if ok {
		gClient.SetBaseURL(baseURL)
	}

	c := &Client{gitlabClient: gClient}

	return c
}

// GetMergeRequests is exported
func (client *Client) GetMergeRequests() []*g.MergeRequestApprovals {
	lastMonth := time.Now().AddDate(0, -1, 0)
	mergeRequestOptions := &g.ListMergeRequestsOptions{
		State:        g.String("opened"),
		OrderBy:      g.String("created_at"),
		Scope:        g.String("all"),
		CreatedAfter: &lastMonth,
	}

	mergeRequests, _, err := client.gitlabClient.MergeRequests.ListMergeRequests(mergeRequestOptions)
	if err != nil {
		log.Fatal(err)
	}
	mergeRequestCount := len(mergeRequests)

	jobs := make(chan *g.MergeRequest, mergeRequestCount)
	results := make(chan *g.MergeRequestApprovals, mergeRequestCount)
	go client.approvalsWorker(jobs, results)
	go client.approvalsWorker(jobs, results)

	for _, mergeRequest := range mergeRequests {
		jobs <- mergeRequest
	}
	close(jobs)

	var approvals []*g.MergeRequestApprovals
	for i := 0; i < mergeRequestCount; i++ {
		approvals = append(approvals, <-results)
	}

	return approvals
}

func (client *Client) approvalsWorker(jobs <-chan *g.MergeRequest, results chan<- *g.MergeRequestApprovals) {
	for mergeRequest := range jobs {
		approvals, _, err := client.gitlabClient.MergeRequests.GetMergeRequestApprovals(mergeRequest.ProjectID, mergeRequest.IID)
		if err != nil {
			log.Fatal(err)
		}
		results <- approvals
		// mockApproval := g.MergeRequestApprovals{Title: mergeRequest.Title}
		// results <- &mockApproval
	}
}
