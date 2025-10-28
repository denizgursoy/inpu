package main

import (
	"fmt"
	"log"
	"net/http"

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
		Use(inpu.LoggingMiddleware(false, false)).
		Use(inpu.RequestIDMiddleware())

	filteredTodos := make([]ToDo, 0)
	err := client.Get("/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		OnReplyIf(inpu.StatusIsOk, inpu.ThenUnmarshalJsonTo(&filteredTodos)).
		OnReplyIf(inpu.StatusAnyExcept(http.StatusOK), inpu.ThenReturnDefaultError).
		Send()

	for i := range filteredTodos {
		fmt.Println(i+1, "-", filteredTodos[i].Title)
	}

	newTodo := ToDo{
		UserID:    22,
		Title:     "",
		Completed: true,
	}
	createdToDo := CreatedTodo{}
	err = client.
		Post("/todos", inpu.BodyJson(newTodo)).
		OnReplyIf(inpu.StatusIsOk, inpu.ThenUnmarshalJsonTo(&createdToDo)).
		Send()
	if err != nil {
		log.Fatal(err)
	}
}
