package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	// _ "github.com/joho/godotenv/autoload"
	"github.com/joho/godotenv"

	"github.com/ycchuang99/todo-list/models"
	"github.com/ycchuang99/todo-list/controllers"
)

func main() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	models.ConnectDatabase()

	router.GET("/api/v1/todo-list", controllers.GetTodoList)
	router.POST("/api/v1/todo-list", controllers.PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", controllers.DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", controllers.PutTodoList)
	router.PUT("/api/v1/todo-list/:id/done", controllers.DoneTodoList)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(router.Run(":8000"))
}
