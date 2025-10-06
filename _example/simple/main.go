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
	response, err := inpu.Get("https://jsonplaceholder.typicode.com/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		Send()
	if err != nil {
		log.Fatal(err)
	}

	if response.IsSuccess() {
		filteredTodos := make([]ToDo, 0)
		if err := response.UnmarshalJson(&filteredTodos); err != nil {
			log.Fatal(err)
		}
		for i := range filteredTodos {
			fmt.Println(i+1, "-", filteredTodos[i].Title)
		}
	} else if response.IsServerError() {
		log.Fatal("server failed")
	}
}
