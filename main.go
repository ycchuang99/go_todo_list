package main

import (
	"net/http"
	"fmt"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DoneAt      *string `json:"done_at"`
}

var todoList = []Todo{}
var db *sql.DB
var err error

func GetTodoList(c *gin.Context) {
	var todoList = []Todo{}

	result, err := db.Query("SELECT * FROM todo_list")

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var todo Todo

		err := result.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.DoneAt)

		if err != nil {
			panic(err.Error())
		}

		todoList = append(todoList, todo)
	}

	c.JSON(http.StatusOK, gin.H{"data": todoList})
}

func PostTodoList(c *gin.Context) {
	var todo Todo
	stmt, err := db.Prepare("INSERT INTO todo_list(title, description) VALUES(?, ?)")

	if err != nil {
		panic(err.Error())
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

	result, err := stmt.Exec(title, description)

	if err != nil {
		panic(err.Error())
	}

	id, _ := result.LastInsertId()

	stmt, err = db.Prepare("SELECT * FROM todo_list WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	row := stmt.QueryRow(id)

	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.DoneAt)

	if err != nil {
		panic(err.Error())
	}

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
			doneAt := c.PostForm("done_at")
			todoList[index].DoneAt = &doneAt

			c.JSON(http.StatusOK, gin.H{"data": todoList[index]})

			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})

}

func main() {
	initDB()
	router := gin.Default()

	// router.Use(Logger())

	router.GET("/api/v1/todo-list", GetTodoList)
	router.POST("/api/v1/todo-list", PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", PutTodoList)

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(router.Run(":8000"))
}

func initDB() {
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo_list")

	if err != nil {
		panic(err.Error())
	}

	// defer db.Close()
}
