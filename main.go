package main

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v28/github"
	"gopkg.in/yaml.v2"
)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func getCurrentLabels(pr *github.PullRequest) []string {
	var labelSet []string
	for _, l := range pr.Labels {
		labelSet = append(labelSet, *l.Name)
	}
	return labelSet
}

func containsLabels(expected []string, current []string) bool {
	for _, e := range expected {
		if !contains(current, e) {
			return false
		}
	}
	return true
}

// Get files and labels matchers, output labels
func matchFiles(labelsMatch map[string][]string, files []*github.CommitFile) []string {
	var labelSet []string
	set := make(map[string]bool)
	for _, file := range files {
		for label, matchers := range labelsMatch {
			for _, pattern := range matchers {
				match, _ := path.Match(pattern, *file.Filename)
				if match && !set[label] {
					set[label] = true
					labelSet = append(labelSet, label)
					break
				}
			}
		}
	}
	return labelSet
}

func main() {
	var owner, repo, token, configpath string
	repoSlug, _ := os.LookupEnv("GITHUB_REPOSITORY")
	token, _ = os.LookupEnv("GITHUB_TOKEN")
	configpath, exists := os.LookupEnv("LABEL_MAPPINGS_FILE")
	if !exists {
		configpath = ".github/labeler.yml"
	}
	s := strings.Split(repoSlug, "/")
	owner = s[0]
	repo = s[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		//TODO: access token should be passed as CLI parameter
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	content, _, _, err := client.Repositories.GetContents(context.Background(), owner, repo, configpath, nil)
	if err != nil {
		glog.Fatal(err)
	}

	yamlFile, err := content.GetContent()
	if err != nil {
		glog.Fatal(err)
	}

	var labelMatchers map[string][]string
	err = yaml.Unmarshal([]byte(yamlFile), &labelMatchers)
	if err != nil {
		glog.Fatal(err)
	}

	opt := &github.PullRequestListOptions{State: "open", Sort: "updated"}
	// get all pages of results
	for {
		pulls, resp, err := client.PullRequests.List(context.Background(), owner, repo, opt)
		if err != nil {
			glog.Fatal(err)
		}
		for _, pull := range pulls {
			files, _, err := client.PullRequests.ListFiles(context.Background(), owner, repo, *pull.Number, nil)
			if err != nil {
				glog.Error(err)
			}
			expectedLabels := matchFiles(labelMatchers, files)
			if !containsLabels(expectedLabels, getCurrentLabels(pull)) {
				client.Issues.AddLabelsToIssue(context.Background(), owner, repo, *pull.Number, expectedLabels)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
