package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type IssueState string

const (
	Open   IssueState = "open"
	Closed IssueState = "closed"
)

type CreateIssueRequest struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Labels    []string `json:"labels"`
	Assignees []string `json:"assignees"`
}

type EditIssueRequest struct {
	State IssueState `json:"state"`
}

func CreateGitHubIssues(owner string, repo string, token string, issues []CreateIssueRequest) error {
	url := "http://api.github.com/repos/" + owner + "/" + repo + "/issues"

	for i := 0; i < len(issues); i++ {
		err := sendRequest(url, "POST", url, owner, issues[i])
		if err != nil {
			return nil
		}
	}

	return nil
}

func CloseGitHubIssues(owner string, repo string, token string, issueIds []int) error {
	for i := 0; i < len(issueIds); i++ {
		id := strconv.Itoa(issueIds[i])
		url := "http://api.github.com/repos/" + owner + "/" + repo + "/issues/" + id
		request := EditIssueRequest{State: Closed}
		err := sendRequest(url, "PATCH", token, owner, request)
		if err != nil {
			return nil
		}
	}

	return nil
}

func sendRequest(url string, method string, token string, user string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	time.Sleep(time.Second)

	return nil
}
