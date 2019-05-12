package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/issues"

	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		fmt.Printf("Creating issue (%d/%d) %s\n", i+1, len(issues), issue.Title)
		err := sendRequest(url, "POST", token, owner, issue)
		if err != nil {
			return err
		}
	}

	return nil
}

func CloseGitHubIssues(owner string, repo string, token string, issueIds []int) error {
	for i := 0; i < len(issueIds); i++ {
		id := strconv.Itoa(issueIds[i])
		fmt.Printf("Closing issue (%d/%d)\n", i+1, len(issueIds))
		url := "https://api.github.com/repos/" + owner + "/" + repo + "/issues/" + id
		request := EditIssueRequest{State: Closed}
		err := sendRequest(url, "PATCH", token, owner, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendRequest(url string, method string, token string, user string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
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

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Status code was " + strconv.Itoa(resp.StatusCode) + "\n" + string(body))
	}

	time.Sleep(time.Second)

	return nil
}
