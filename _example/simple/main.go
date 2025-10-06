package main

import (
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
	response, err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		Send()
	if err != nil {
		log.Fatal(err)
	}

	filteredTodos := make([]ToDo, 0)
	if response.IsSuccess() {
		if err := response.UnmarshalJson(&filteredTodos); err != nil {
			log.Fatal(err)
		}
	} else if response.IsServerError() {
		log.Fatal("server failed")
	}
}
