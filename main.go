package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Board struct {
	Tasks      []Task      `json:"cards"`
	Checklists []Checklist `json:"checklists"`
	Actions    []Action    `json:"actions"`
	Lists      []List      `json:"lists"`
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

func main() {
	boardJson, err := os.Open("example-trello-tasks.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer boardJson.Close()

	bytes, _ := ioutil.ReadAll(boardJson)

	var board Board
	json.Unmarshal(bytes, &board)

	for i := 0; i < len(board.Tasks); i++ {
		task := board.Tasks[i]

		fmt.Print("\n\n---------------------------------\n\n")

		fmt.Printf("Id: %d\n", task.Id)
		fmt.Println("Title: " + task.Title)
		fmt.Println("Description: " + task.Description)
		fmt.Println("Url: " + task.Url)

		list, _ := task.List(board.Lists)
		fmt.Println("List: " + list.Name)

		checklists := task.Checklists(board.Checklists)
		for j := 0; j < len(checklists); j++ {
			for k := 0; k < len(checklists[j].CheckItems); k++ {
				fmt.Println("Checklist item: " + checklists[j].CheckItems[k].State + " " + checklists[j].CheckItems[k].Name)
			}
		}

		for j := 0; j < len(task.Attachments); j++ {
			fmt.Println("Attachment: " + task.Attachments[j].Name + " " + task.Attachments[j].Url)
		}

		actions := task.Actions(board.Actions)
		for j := 0; j < len(actions); j++ {
			fmt.Println("Action: " + actions[j].Type + " " + actions[j].Date)

			if actions[j].Type == "commentCard" {
				fmt.Println("Comment: " + actions[j].Data.Comment)
			}
		}
	}
}