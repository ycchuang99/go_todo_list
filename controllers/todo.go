package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ycchuang99/todo-list/models"
)

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

	var todoList []models.Todo = []models.Todo{}
	var totalItems int

	err = models.DB.QueryRow("SELECT COUNT(*) FROM todo_list").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	rows, err := models.DB.Query("SELECT * FROM todo_list LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo

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
	var todo models.Todo

	title := c.PostForm("title")
	description := c.PostForm("description")

	stmt, err := models.DB.Prepare("INSERT INTO todo_list(title, description) VALUES(?, ?)")
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

	row := models.DB.QueryRow("SELECT * FROM todo_list WHERE id = ?", id)
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

	models.DB.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	stmt, err := models.DB.Prepare("DELETE FROM todo_list WHERE id = ?")
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
	err := models.DB.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")

	stmt, err := models.DB.Prepare("UPDATE todo_list SET title = ?, description = ? WHERE id = ?")
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

	var todo models.Todo

	row := models.DB.QueryRow("SELECT * FROM todo_list WHERE id = ?", id)
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
	err := models.DB.QueryRow("SELECT COUNT(*) FROM todo_list WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	stmt, err := models.DB.Prepare("UPDATE todo_list SET done_at = NOW() WHERE id = ?")
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
