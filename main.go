package main

import (
	"fmt"
	"github.com/danielpieper/mrcli-go/gitlab"
)

func main() {
	client := gitlab.NewClient(nil)
	projects := client.GetProjects()
	fmt.Println(projects)

	// pendingRequests := client.GetPendingRequests()
	// for _, pr := range pendingRequests {
	// 	fmt.Println(pr.Request.Title)
	// 	fmt.Println(pr.Request.WebURL)
	// }
}
