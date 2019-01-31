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

// PendingRequest is exported
type PendingRequest struct {
	Project   *g.Project
	Request   *g.MergeRequest
	Approvals *g.MergeRequestApprovals
}

// NewClient is exported
func NewClient(httpClient *http.Client) *Client {
	token := os.Getenv("GITLAB_TOKEN")
	gClient := g.NewClient(httpClient, token)
	baseURL, ok := os.LookupEnv("GITLAB_URL")
	if ok {
		gClient.SetBaseURL(baseURL)
	}

	return &Client{gitlabClient: gClient}
}

func (client *Client) getProjects() ([]*g.Project, error) {
	listProjectOptions := &g.ListProjectsOptions{
		Archived:                 g.Bool(false),
		WithMergeRequestsEnabled: g.Bool(true),
		ListOptions:              g.ListOptions{Page: 1, PerPage: 100},
	}
	result, _, err := client.gitlabClient.Projects.ListProjects(listProjectOptions)

	return result, err
}

func (client *Client) getMergeRequests() ([]*g.MergeRequest, error) {
	lastMonth := time.Now().AddDate(0, -1, 0)
	mergeRequestOptions := &g.ListMergeRequestsOptions{
		State:        g.String("opened"),
		OrderBy:      g.String("created_at"),
		Scope:        g.String("all"),
		CreatedAfter: &lastMonth,
		ListOptions:  g.ListOptions{Page: 1, PerPage: 100},
	}

	result, _, err := client.gitlabClient.MergeRequests.ListMergeRequests(mergeRequestOptions)
	if err != nil {
		return nil, err
	}
	var mergeRequests []*g.MergeRequest
	for _, mergeRequest := range result {
		if !mergeRequest.WorkInProgress {
			mergeRequests = append(mergeRequests, mergeRequest)
		}
	}
	return mergeRequests, nil
}

// GetPendingRequests fetches a list of up to 100 pending
// merge requests and adds project and approvals information
func (client *Client) GetPendingRequests() ([]*PendingRequest, error) {
	mergeRequests, err := client.getMergeRequests()
	if err != nil {
		return nil, err
	}
	mergeRequestCount := len(mergeRequests)

	projects, err := client.getProjects()
	if err != nil {
		return nil, err
	}

	jobs := make(chan *g.MergeRequest, mergeRequestCount)
	results := make(chan *PendingRequest, mergeRequestCount)
	go client.approvalsWorker(jobs, results)
	go client.approvalsWorker(jobs, results)

	for _, mergeRequest := range mergeRequests {
		jobs <- mergeRequest
	}
	close(jobs)

	var approvals []*PendingRequest
	for i := 0; i < mergeRequestCount; i++ {
		result := <-results
		if result.Approvals.ApprovalsLeft > 0 {
			for _, project := range projects {
				if project.ID == result.Request.ProjectID {
					result.Project = project
					break
				}
			}
			approvals = append(approvals, result)
		}
	}

	return approvals, nil
}

func (client *Client) approvalsWorker(jobs <-chan *g.MergeRequest, results chan<- *PendingRequest) {
	for mergeRequest := range jobs {
		approvals, _, err := client.gitlabClient.MergeRequests.GetMergeRequestApprovals(mergeRequest.ProjectID, mergeRequest.IID)
		if err != nil {
			log.Fatal(err)
		}
		results <- &PendingRequest{Request: mergeRequest, Approvals: approvals}
	}
}
