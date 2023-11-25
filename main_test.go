package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
    "os"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/ycchuang99/todo-list/controllers"
	"github.com/ycchuang99/todo-list/models"
)

func SetUpRouter() *gin.Engine {
    router := gin.Default()
    gin.SetMode(os.Getenv("GIN_MODE"))

    return router
}

func SetUpTestEnvironment() {
    err := godotenv.Load(".env.test")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func TearDownTestDatabase() {
    _, err := models.DB.Exec("TRUNCATE TABLE todo_list")
    if err != nil {
        log.Fatal("Error truncating todo_list table")
    }
}

func TestGetTodoList(t *testing.T) {
    router := SetUpRouter()

    SetUpTestEnvironment()
    models.ConnectDatabase()

    _, err := models.DB.Exec("INSERT INTO todo_list (title, description) VALUES (?, ?)", "Test Title", "Test Description")
    if err != nil {
        t.Errorf("Failed to insert todo_list: %v", err)
    }
    _, err = models.DB.Exec("INSERT INTO todo_list (title, description) VALUES (?, ?)", "Test Title 2", "Test Description 2")
    if err != nil {
        t.Errorf("Failed to insert todo_list: %v", err)
    }
    _, err = models.DB.Exec("INSERT INTO todo_list (title, description) VALUES (?, ?)", "Test Title 3", "Test Description 3")
    if err != nil {
        t.Errorf("Failed to insert todo_list: %v", err)
    }

    router.GET("/api/v1/todo-list", controllers.GetTodoList)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/todo-list", nil)
    router.ServeHTTP(w, req)

    var jsonResponse map[string]interface{}
    err = json.Unmarshal(w.Body.Bytes(), &jsonResponse)
    if err != nil {
        t.Errorf("Failed to unmarshal response body: %v", err)
    }

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, 3, len(jsonResponse["data"].([]interface{})))
    
    assert.Equal(t, "Test Title", jsonResponse["data"].([]interface{})[0].(map[string]interface{})["title"])
    assert.Equal(t, "Test Description", jsonResponse["data"].([]interface{})[0].(map[string]interface{})["description"])
    assert.Equal(t, "Test Title 2", jsonResponse["data"].([]interface{})[1].(map[string]interface{})["title"])
    assert.Equal(t, "Test Description 2", jsonResponse["data"].([]interface{})[1].(map[string]interface{})["description"])
    assert.Equal(t, "Test Title 3", jsonResponse["data"].([]interface{})[2].(map[string]interface{})["title"])
    assert.Equal(t, "Test Description 3", jsonResponse["data"].([]interface{})[2].(map[string]interface{})["description"])

    assert.Equal(t, 1, int(jsonResponse["page"].(float64)))
    assert.Equal(t, 10, int(jsonResponse["perPage"].(float64)))
    assert.Equal(t, 3, int(jsonResponse["totalItems"].(float64)))
    assert.Equal(t, 1, int(jsonResponse["totalPages"].(float64)))

    TearDownTestDatabase()
}

func TestPostTodoList(t *testing.T) {
    SetUpTestEnvironment()

    router := SetUpRouter()

    models.ConnectDatabase()

    router.POST("/api/v1/todo-list", controllers.PostTodoList)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/todo-list", nil)
    req.PostForm = map[string][]string{
        "title": {"Test Title"},
        "description": {"Test Description"},
    }
    router.ServeHTTP(w, req)

    var jsonResponse map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
    if err != nil {
        t.Errorf("Failed to unmarshal response body: %v", err)
    }

    assert.Equal(t, 201, w.Code)
    assert.Equal(t, "Test Title", jsonResponse["data"].(map[string]interface{})["title"])
    assert.Equal(t, "Test Description", jsonResponse["data"].(map[string]interface{})["description"])

    TearDownTestDatabase()
}

func TestDeleteTodoList(t *testing.T) {
    SetUpTestEnvironment()

    router := SetUpRouter()

    models.ConnectDatabase()

    _, err := models.DB.Exec("INSERT INTO todo_list (title, description) VALUES (?, ?)", "Test Title", "Test Description")
    if err != nil {
        t.Errorf("Failed to insert todo_list: %v", err)
    }

    router.DELETE("/api/v1/todo-list/:id", controllers.DeleteTodoList)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/api/v1/todo-list/1", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 204, w.Code)

    var count int
    models.DB.QueryRow("SELECT COUNT(*) FROM todo_list").Scan(&count)
    assert.Equal(t, 0, count)

    TearDownTestDatabase()
}

func TestPutTodoList(t *testing.T) {
    SetUpTestEnvironment()

    router := SetUpRouter()

    models.ConnectDatabase()

    models.DB.Exec("INSERT INTO todo_list (title, description) VALUES (?, ?)", "Test Title", "Test Description")

    router.PUT("/api/v1/todo-list/:id", controllers.PutTodoList)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/api/v1/todo-list/1", nil)
    req.PostForm = map[string][]string{
        "title": {"Test Title Updated"},
        "description": {"Test Description Updated"},
    }
    router.ServeHTTP(w, req)

    var jsonResponse map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
    if err != nil {
        t.Errorf("Failed to unmarshal response body: %v", err)
    }

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "Test Title Updated", jsonResponse["data"].(map[string]interface{})["title"])
    assert.Equal(t, "Test Description Updated", jsonResponse["data"].(map[string]interface{})["description"])

    TearDownTestDatabase()
}
