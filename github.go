package main

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
