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
		UseMiddlewares(inpu.LoggingMiddleware(true, false))

	filteredTodos := make([]ToDo, 0)
	err := client.Get("/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		OnReply(inpu.StatusIsOk, inpu.UnmarshalJson(&filteredTodos)).
		OnReply(inpu.StatusAnyExcept(http.StatusOK), inpu.ReturnDefaultError).
		Send()

	for i := range filteredTodos {
		fmt.Println(i+1, "-", filteredTodos[i].Title)
	}

	newTodo := ToDo{
		UserID:    22,
		Title:     "",
		Completed: true,
	}
	err = client.
		Post("/todos", inpu.BodyJson(newTodo)).
		OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&newTodo)).
		Send()
	if err != nil {
		log.Fatal(err)
	}
}
