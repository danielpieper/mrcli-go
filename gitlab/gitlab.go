package gitlab

import (
	"fmt"
	"github.com/bclicn/color"
	g "github.com/xanzy/go-gitlab"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
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

// AuthenticatedUser is exported
func (client *Client) AuthenticatedUser() (*g.User, error) {
	user, _, err := client.gitlabClient.Users.CurrentUser()
	return user, err
}

// PendingRequests fetches a list of up to 100 pending
// merge requests and adds project and approvals information
func (client *Client) PendingRequests() (RankedPendingRequests, error) {
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

	approvals := RankedPendingRequests{}
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
	sort.Sort(approvals)

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

// RankedPendingRequests is exported
type RankedPendingRequests []*PendingRequest

func (r RankedPendingRequests) Len() int { return len(r) }
func (r RankedPendingRequests) Less(i, j int) bool {
	return r[i].Request.CreatedAt.Before(*r[j].Request.CreatedAt)
}
func (r RankedPendingRequests) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// PendingRequest is exported
type PendingRequest struct {
	Project   *g.Project
	Request   *g.MergeRequest
	Approvals *g.MergeRequestApprovals
}

// ApproverNames is exported
func (pr *PendingRequest) ApproverNames() (approverNames []string) {
	for _, approver := range pr.Approvals.Approvers {
		approverNames = append(approverNames, approver.User.Username)
	}
	return
}

// Color is exported
func (pr *PendingRequest) Color(value string) string {
	age := math.Round(time.Since(*pr.Request.CreatedAt).Hours() / 24)

	switch {
	case age > 2:
		return color.Red(value)
	case age > 1:
		return color.Yellow(value)
	}

	return color.Green(value)
}

// HumanReadableCreatedAtDiff is exported
func (pr *PendingRequest) HumanReadableCreatedAtDiff() string {
	return humanReadableTimeDiff(*pr.Request.CreatedAt)
}

// HumanReadableUpdatedAtDiff is exported
func (pr *PendingRequest) HumanReadableUpdatedAtDiff() string {
	return humanReadableTimeDiff(*pr.Request.UpdatedAt)
}

func humanReadableTimeDiff(value time.Time) string {
	duration := time.Since(value)

	age := []string{}

	weeks := math.Floor(duration.Hours() / 24 / 7)
	if weeks == 1 {
		age = append(age, "1 week")
	} else if weeks > 1 {
		age = append(age, fmt.Sprintf("%d weeks", int(weeks)))
	}

	days := math.Floor(duration.Hours()/24 - weeks*7)
	if days == 1 {
		age = append(age, "1 day")
	} else if days > 1 {
		age = append(age, fmt.Sprintf("%d days", int(days)))
	}

	hours := int(math.Floor(duration.Hours() - days*24 - weeks*7*24))
	if hours == 1 {
		age = append(age, "1 hour")
	} else if hours > 1 {
		age = append(age, fmt.Sprintf("%d hours", hours))
	}
	age = append(age, "ago")

	return strings.Join(age, " ")
}
