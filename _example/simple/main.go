package main

import (
	"fmt"
	"log"

	"github.com/denizgursoy/inpu"
)

type ToDo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	filteredTodos := make([]ToDo, 0)

	err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&filteredTodos)).
		Send()
	if err != nil {
		log.Fatal(err)
	}

	for i := range filteredTodos {
		fmt.Println(i+1, "-", filteredTodos[i].Title)
	}

}
