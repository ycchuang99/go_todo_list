package main

import (
	"net/http"
	"fmt"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	IsDone bool `json:"is_done"`
}

var todoList = []Todo{}

func GetTodoList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": todoList})
}

func PostTodoList(c *gin.Context) {
	var todo Todo
	
	todo.ID = len(todoList) + 1
	todo.Title = c.PostForm("title")
	todo.Description = c.PostForm("description")
	todo.IsDone = false

	todoList = append(todoList, todo)

	c.JSON(http.StatusCreated, gin.H{"data": todo})
}

func DeleteTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for index, todo := range todoList {
		if todo.ID == id {
			todoList = append(todoList[:index], todoList[index+1:]...)
			c.Status(http.StatusNoContent)

			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}

func PutTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for index, todo := range todoList {
		if todo.ID == id {
			todoList[index].Title = c.PostForm("title")
			todoList[index].Description = c.PostForm("description")
			todoList[index].IsDone = c.PostForm("is_done") == "true"

			c.JSON(http.StatusOK, gin.H{"data": todoList[index]})

			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}

func main() {
	router := gin.Default()

	router.GET("/api/v1/todo-list", GetTodoList)
	router.POST("/api/v1/todo-list", PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", PutTodoList)

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(router.Run(":8080"))
}