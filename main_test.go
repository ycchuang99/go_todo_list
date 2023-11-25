package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "os"
    "github.com/gin-gonic/gin"

    "github.com/stretchr/testify/assert"
)

func TestGetTodoList(t *testing.T) {
    r := setupRouter()

    req, _ := http.NewRequest("GET", "/api/v1/todo-list?page=1&perPage=10", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Add more assertions to check the response body if needed
}

func TestPostTodoList(t *testing.T) {
    r := setupRouter()

    requestBody := `{"title": "Test Todo", "description": "Test Description"}`
    req, _ := http.NewRequest("POST", "/api/v1/todo-list", strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    // Add more assertions to check the response body if needed
}

func TestDeleteTodoList(t *testing.T) {
    r := setupRouter()

    req, _ := http.NewRequest("DELETE", "/api/v1/todo-list/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPutTodoList(t *testing.T) {
    r := setupRouter()

    requestBody := `{"title": "Updated Todo", "description": "Updated Description"}`
    req, _ := http.NewRequest("PUT", "/api/v1/todo-list/1", strings.NewReader(requestBody))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Add more assertions to check the response body if needed
}

func TestDoneTodoList(t *testing.T) {
    r := setupRouter()

    req, _ := http.NewRequest("PUT", "/api/v1/todo-list/1/done", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNoContent, w.Code)
}

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/api/v1/todo-list", GetTodoList)
    r.POST("/api/v1/todo-list", PostTodoList)
    r.DELETE("/api/v1/todo-list/:id", DeleteTodoList)
    r.PUT("/api/v1/todo-list/:id", PutTodoList)
    r.PUT("/api/v1/todo-list/:id/done", DoneTodoList)
    return r
}

func TestMain(m *testing.M) {
    // Run the tests
    exitCode := m.Run()

    // Exit with the appropriate exit code
    os.Exit(exitCode)
}
