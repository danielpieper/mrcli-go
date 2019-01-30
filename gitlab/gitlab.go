package gitlab

import (
	"fmt"
	g "github.com/xanzy/go-gitlab"
	"log"
	"net/http"
	"os"
	"time"
)

// GitlabClient is exported
type Client struct {
	gitlabClient *g.Client
}

func NewClient(httpClient *http.Client) *Client {
	token := os.Getenv("GITLAB_TOKEN")
	gitlabClient := g.NewClient(httpClient, token)
	baseUrl, ok := os.LookupEnv("GITLAB_URL")
	if ok {
		gitlabClient.SetBaseURL(baseUrl)
	}

	c := &Client{gitlabClient: gitlabClient}

	return c
}

func (client *Client) GetMergeRequests() {
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
	for _, mergeRequest := range mergeRequests {
		fmt.Println(mergeRequest.Title)
	}

}
