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
func (client *Client) GetMergeRequests() []*g.MergeRequest {
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
	return mergeRequests
}

// GetMergeRequestApprovals is exported
func (client *Client) GetMergeRequestApprovals(mergeRequest *g.MergeRequest) *g.MergeRequestApprovals {
	approvals, _, err := client.gitlabClient.MergeRequests.GetMergeRequestApprovals(mergeRequest.ProjectID, mergeRequest.IID)
	if err != nil {
		log.Fatal(err)
	}
	return approvals
}
