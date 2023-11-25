package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"

	"github.com/ycchuang99/todo-list/controllers"
	"github.com/ycchuang99/todo-list/models"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	InitEnv()
	router := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))

	models.ConnectDatabase()

	router.GET("/api/v1/todo-list", controllers.GetTodoList)
	router.POST("/api/v1/todo-list", controllers.PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", controllers.DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", controllers.PutTodoList)
	router.PUT("/api/v1/todo-list/:id/done", controllers.DoneTodoList)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(router.Run(":8000"))
}
