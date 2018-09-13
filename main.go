package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <username/source_repo> <username/dest_repo>", os.Args[0])
		os.Exit(0)
	}

	src := strings.Split(os.Args[1], "/")
	dest := strings.Split(os.Args[2], "/")

	ghToken := os.Getenv("GH_TOKEN")
	if ghToken == "" {
		fmt.Printf("You have to pass GH_TOKEN pointing to valid personal access token")
		os.Exit(2)
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	opt := &github.IssueListByRepoOptions{}

	issues, _, err := ghClient.Issues.ListByRepo(ctx, src[0], src[1], opt)

	if err != nil {
		fmt.Printf("Could not get list of issues from the repo provided")
		os.Exit(1)
	}

	for _, issue := range issues {
		body := fmt.Sprintf("%s\n\n_[Source Issue](%s)_", *issue.Body, *issue.URL)
		var labels = []string{}
		for _, label := range issue.Labels {
			labels = append(labels, *label.Name)
		}

		var assignees = []string{}
		for _, user := range issue.Assignees {
			assignees = append(assignees, *user.Login)
		}

		issueReq := &github.IssueRequest{
			Title:     issue.Title,
			Body:      &body,
			Labels:    &labels,
			State:     issue.State,
			Assignees: &assignees,
		}

		fmt.Printf("Copying %s...\n", *issue.Title)
		_, _, err := ghClient.Issues.Create(ctx, dest[0], dest[1], issueReq)

		if err != nil {
			fmt.Printf("An error encountered while copying %s.\n", *issue.Title)
		}
	}
}
