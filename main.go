package main

import (
	"fmt"
)

func main() {
	board, err := ReadBoard("example-trello-tasks.json")
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
