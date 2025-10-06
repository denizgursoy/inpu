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
		UseMiddlewares(inpu.LoggingMiddleware(true))

	filteredTodos := make([]ToDo, 0)
	err := client.Get("/todos").
		QueryBool("completed", true).
		QueryInt("userId", 2).
		OnReply(inpu.StatusIs(http.StatusCreated), inpu.ReturnError(inpu.ErrConnectionFailed)).
		OnReply(inpu.StatusIsOneOf(http.StatusMovedPermanently, http.StatusAccepted), inpu.UnmarshalJson(filteredTodos)).
		OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&filteredTodos)).
		OnReply(inpu.StatusIs(http.StatusOK), inpu.UnmarshalJson(&filteredTodos)).
		OnReply(inpu.StatusAny, inpu.UnmarshalJson(&filteredTodos)).
		OnReply(inpu.StatusAnyExcept(http.StatusBadRequest), inpu.UnmarshalJson(&filteredTodos)).
		OnReply(inpu.StatusAnyExceptOneOf(http.StatusMultipleChoices, http.StatusBadRequest), inpu.UnmarshalJson(&filteredTodos)).
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
		Post("/todos", newTodo).
		OnReply(inpu.StatusIsSuccess, inpu.UnmarshalJson(&newTodo)).
		Send()
	if err != nil {
		log.Fatal(err)
	}
}
