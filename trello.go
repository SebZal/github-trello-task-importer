package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Board struct {
	Tasks      []Task      `json:"cards"`
	Checklists []Checklist `json:"checklists"`
	Actions    []Action    `json:"actions"`
	Lists      []List      `json:"lists"`
}

func (b Board) CreateDescription(task Task) string {
	description := task.Description + "\n\n"

	checklists := task.Checklists(b.Checklists)
	for j := 0; j < len(checklists); j++ {
		for k := 0; k < len(checklists[j].CheckItems); k++ {
			if checklists[j].CheckItems[k].State == "complete" {
				description += "- [x] "
			} else {
				description += "- [ ] "
			}
			description += checklists[j].CheckItems[k].Name + "\n"
		}
	}

	for j := 0; j < len(task.Attachments); j++ {
		description += "\n![" + task.Attachments[j].Name + "](" + task.Attachments[j].Url + ")"
	}

	var comments []string
	description += "\n\n__Actions:__\n"
	actions := task.Actions(b.Actions)
	for j := 0; j < len(actions); j++ {
		description += "- " + actions[j].Type + " (" + actions[j].Date + ")\n"

		if actions[j].Type == "commentCard" {
			comments = append(comments, actions[j].Data.Comment)
		}
	}

	if len(comments) > 0 {
		description += "\n__Comments:__\n"
	}
	for j := 0; j < len(comments); j++ {
		description += "- " + comments[j]
	}

	description += "\n\n_This task was imported from [" + task.Url + "](" + task.Url + ")_"

	return description
}

type Task struct {
	Id           int          `json:"idShort"`
	Description  string       `json:"desc"`
	Title        string       `json:"name"`
	Url          string       `json:"url"`
	ChecklistIds []string     `json:"idChecklists"`
	Labels       []Label      `json:"labels"`
	Attachments  []Attachment `json:"attachments"`
	ListId       string       `json:"idList"`
}

func (t Task) List(lists []List) (List, error) {
	var list List

	for i := 0; i < len(lists); i++ {
		if lists[i].Id == t.ListId {
			return lists[i], nil
		}
	}

	return list, errors.New("List not found.")
}

func (t Task) Checklists(checklists []Checklist) []Checklist {
	var taskChecklists []Checklist

	for i := 0; i < len(checklists); i++ {
		for j := 0; j < len(t.ChecklistIds); j++ {
			if checklists[i].Id == t.ChecklistIds[j] {
				taskChecklists = append(taskChecklists, checklists[i])
			}
		}
	}

	return taskChecklists
}

func (t Task) Actions(actions []Action) []Action {
	var taskActions []Action

	for i := 0; i < len(actions); i++ {
		if actions[i].Data.Task.Id == t.Id {
			taskActions = append(taskActions, actions[i])
		}
	}

	return taskActions
}

type Checklist struct {
	Id         string      `json:"id"`
	CheckItems []CheckItem `json:"checkItems"`
}

type CheckItem struct {
	State string `json:"state"`
	Name  string `json:"name"`
}

type Label struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Attachment struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Action struct {
	Type string     `json:"type"`
	Date string     `json:"date"`
	Data ActionData `json:"data"`
}

type ActionData struct {
	Task    ActionDataTask `json:"card"`
	Comment string         `json:"text"`
}

type ActionDataTask struct {
	Id int `json:"idShort"`
}

type List struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ReadBoard(jsonPath string) (Board, error) {
	var board Board

	boardJson, err := os.Open(jsonPath)
	if err != nil {
		return board, err
	}
	defer boardJson.Close()

	bytes, err := ioutil.ReadAll(boardJson)
	if err != nil {
		return board, err
	}

	err = json.Unmarshal(bytes, &board)
	if err != nil {
		return board, err
	}

	return board, nil
}
