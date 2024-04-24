package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
)

// Function to check if a specific label exists in a PR
func hasLabel(pr *github.PullRequest, labelName string) bool {
	for _, label := range pr.Labels {
		if strings.ToLower(label.GetName()) == strings.ToLower(labelName) {
			return true
		}
	}
	return false
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	repoOwner := os.Getenv("REPO_OWNER")
	repoName := os.Getenv("REPO_NAME")
	prNumber := os.Getenv("PR_NUMBER")
	labelName := os.Getenv("LABEL_NAME")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the PR to check its labels
	prNumberInt, _ := strconv.Atoi(prNumber)
	pr, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, prNumberInt)
	if err != nil {
		fmt.Println("Error fetching PR:", err)
		return
	}

	// If the PR has the specified label, dismiss all reviews
	if hasLabel(pr, labelName) {
		reviews, _, err := client.PullRequests.ListReviews(ctx, repoOwner, repoName, prNumberInt, nil)
		if err != nil {
			fmt.Println("Error fetching reviews:", err)
			return
		}

		for _, review := range reviews {
			_, err := client.PullRequests.DismissReview(ctx, repoOwner, repoName, prNumberInt, review.GetID(), &github.PullRequestReviewDismissalRequest{
				Message: "Review dismissed due to new commit with specific label",
			})
			if err != nil {
				fmt.Println("Error dismissing review:", err)
			}
		}
	}
}

