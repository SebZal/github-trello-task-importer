package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	GitHubToken         string `json:"gitHubToken"`
	TrelloBoardJsonPath string `json:"trelloBoardJsonPath"`
	User                string `json:"user"`
	Repository          string `json:"repository"`
}

func ReadConfig() Config {
	configJson, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configJson.Close()

	bytes, err := ioutil.ReadAll(configJson)
	if err != nil {
		fmt.Println(err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func main() {
	config := ReadConfig()

	board, err := ReadBoard(config.TrelloBoardJsonPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(board.Tasks); i++ {
		task := board.Tasks[i]

		fmt.Print("\n\n---------------------------------\n\n")

		fmt.Printf("Id: %d\n", task.Id)
		fmt.Println("Title: " + task.Title)

		list, _ := task.List(board.Lists)
		fmt.Println("List: " + list.Name)

		fmt.Println(board.CreateDescription(task))
	}
}
