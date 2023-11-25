package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/ycchuang99/todo-list/controllers"
	"github.com/ycchuang99/todo-list/models"
)

func SetUpRouter() *gin.Engine {
    router := gin.Default()

    return router
}

func SetUpTestEnvironment() {
    err := godotenv.Load(".env.test")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func TestGetTodoList(t *testing.T) {
    router := SetUpRouter()

    SetUpTestEnvironment()
    models.ConnectDatabase()

    router.GET("/api/v1/todo-list", controllers.GetTodoList)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/todo-list", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "{\"data\":[],\"page\":1,\"perPage\":10,\"totalItems\":0,\"totalPages\":0}", w.Body.String())
}
