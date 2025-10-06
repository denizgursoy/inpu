package main

import (
	"fmt"
	"log"

	"github.com/denizgursoy/inpu"
)

type ToDo struct {
	UserID    int    `json:"userID"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
type CreatedTodo struct {
	ID int `json:"id"`
}

func main() {
	client := inpu.
		New().
		BasePath("https://jsonplaceholder.typicode.com").
		UseMiddlewares(inpu.LoggingMiddleware(true))

	response, err := client.Get("/todos").
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

	newTodo := ToDo{
		UserID:    22,
		Title:     "",
		Completed: true,
	}
	send, err := client.Post("/todos", newTodo).Send()
	if err != nil {
		log.Fatal(err)
	}
	if send.IsSuccess() {
		newTodo := CreatedTodo{}
		if err := send.UnmarshalJson(&newTodo); err != nil {
			log.Fatal(err)
		}
		fmt.Println("created with id", newTodo.ID)
	}
}
