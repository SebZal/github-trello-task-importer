package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type Config struct {
	GitHubToken         string `json:"gitHubToken"`
	TrelloBoardJsonPath string `json:"trelloBoardJsonPath"`
	User                string `json:"user"`
	Repository          string `json:"repository"`
}

func ReadConfig() (Config, error) {
	var config Config

	configJson, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer configJson.Close()

	bytes, err := ioutil.ReadAll(configJson)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	board, err := ReadBoard(config.TrelloBoardJsonPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var issues []CreateIssueRequest
	var issueIdsToClose []int

	sort.Slice(board.Tasks, func(i, j int) bool {
		return board.Tasks[i].Id < board.Tasks[j].Id
	})

	for i := 0; i < len(board.Tasks); i++ {
		task := board.Tasks[i]

		labels := []string{}
		for j := 0; j < len(task.Labels); j++ {
			labels = append(labels, task.Labels[j].Name)
		}

		request := CreateIssueRequest{
			Title:     task.Title,
			Body:      board.CreateDescription(task),
			Assignees: []string{config.User},
			Labels:    labels,
		}

		issues = append(issues, request)

		list, err := task.List(board.Lists)
		if err != nil {
			fmt.Println(err)
			return
		}

		if list.Name == "Done" || list.Name == "Testing" {
			issueIdsToClose = append(issueIdsToClose, task.Id)
		}
	}

	err = CreateGitHubIssues(config.User, config.Repository, config.GitHubToken, issues)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = CloseGitHubIssues(config.User, config.Repository, config.GitHubToken, issueIdsToClose)
	if err != nil {
		fmt.Println(err)
		return
	}
}
