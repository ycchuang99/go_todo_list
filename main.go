package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ycchuang99/todo-list/models"
	"github.com/ycchuang99/todo-list/controllers"
)

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DoneAt      *string `json:"done_at"`
}

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	router.GET("/api/v1/todo-list", controllers.GetTodoList)
	router.POST("/api/v1/todo-list", controllers.PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", controllers.DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", controllers.PutTodoList)
	router.PUT("/api/v1/todo-list/:id/done", controllers.DoneTodoList)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(router.Run(":8000"))
}
