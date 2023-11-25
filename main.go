package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DoneAt      *string `json:"done_at"`
}

var db *sql.DB
var err error

func main() {
	initDB()
	router := gin.Default()

	router.GET("/api/v1/todo-list", GetTodoList)
	router.POST("/api/v1/todo-list", PostTodoList)
	router.DELETE("/api/v1/todo-list/:id", DeleteTodoList)
	router.PUT("/api/v1/todo-list/:id", PutTodoList)
	router.PUT("/api/v1/todo-list/:id/done", DoneTodoList)

	fmt.Println("Starting server on port 8000...")
	log.Fatal(router.Run(":8000"))
}

func initDB() {
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/todo_list")
	if err != nil {
		panic(err.Error())
	}
}

func GetTodoList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid perPage"})
		return
	}

	offset := (page - 1) * perPage

	var todoList []Todo
	var totalItems int

	err = db.QueryRow("SELECT COUNT(*) FROM todo_list").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	rows, err := db.Query("SELECT * FROM todo_list LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.DoneAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}

		todoList = append(todoList, todo)
	}

	totalPages := (totalItems + perPage - 1) / perPage

	response := gin.H{
		"data":       todoList,
		"page":       page,
		"perPage":    perPage,
		"totalPages": totalPages,
		"totalItems": totalItems,
	}

	c.JSON(http.StatusOK, response)
}

func PostTodoList(c *gin.Context) {
	var todo Todo

	title := c.PostForm("title")
	description := c.PostForm("description")

	stmt, err := db.Prepare("INSERT INTO todo_list(title, description) VALUES(?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	id, _ := result.LastInsertId()

	row := db.QueryRow("SELECT * FROM todo_list WHERE id = ?", id)
	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.DoneAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": todo})
}

func DeleteTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var count int

	err = db.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	stmt, err := db.Prepare("DELETE FROM todo_list WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func PutTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

	stmt, err := db.Prepare("UPDATE todo_list SET title = ?, description = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, description, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	var todo Todo

	row := db.QueryRow("SELECT * FROM todo_list WHERE id = ?", id)
	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.DoneAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func DoneTodoList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	stmt, err := db.Prepare("UPDATE todo_list SET done_at = NOW() WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.Status(http.StatusNoContent)
}
