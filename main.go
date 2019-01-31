package main

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
)

func main() {
	client := gitlab.NewClient(nil)

	pendingRequests, err := client.GetPendingRequests()
	if err != nil {
		fmt.Println("An error occured: %v", err)
		return
	}
	for _, pr := range pendingRequests {
		fmt.Println(pr.Request.Title)
		fmt.Println(pr.Request.WebURL)
		fmt.Println(pr.Project.Name)
	}
}
