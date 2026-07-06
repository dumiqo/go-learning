package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	task "to-do/Task"
	todolist "to-do/ToDoList"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// setupRouter настраивает роутер в тестовом режиме
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/list", getList)
	r.POST("/task", addTask)
	r.PUT("/task/:id", editTask)
	return r
}

// clearList очищает глобальный список перед каждым тестом
func clearList() {
	mu.Lock()
	list = todolist.TodoList{} // предполагаем, что это срез
	mu.Unlock()
}

// ----------------------------------------------------------------------------
// Тесты
// ----------------------------------------------------------------------------

// 1. GET /list — пустой список
func TestGetListEmpty(t *testing.T) {
	clearList()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response []task.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse JSON:", err)
	}
	if len(response) != 0 {
		t.Errorf("Expected empty list, got %d items", len(response))
	}
}

// 2. POST /task — успешное создание
func TestAddTaskSuccess(t *testing.T) {
	clearList()
	router := setupRouter()

	newTask := task.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    task.Normal,
		Status:      task.New,
	}
	body, _ := json.Marshal(newTask)

	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var created task.Task
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal("Failed to parse response:", err)
	}
	if created.Id == uuid.Nil {
		t.Error("ID was not set")
	}
	if created.Title != newTask.Title {
		t.Errorf("Expected title %q, got %q", newTask.Title, created.Title)
	}

	// Проверяем, что список действительно содержит новую задачу
	mu.RLock()
	if len(list) != 1 {
		t.Errorf("Expected list length 1, got %d", len(list))
	}
	mu.RUnlock()
}

// 3. POST /task — невалидный JSON (приводит к 400)
func TestAddTaskInvalidJSON(t *testing.T) {
	clearList()
	router := setupRouter()

	// Передаём строку вместо числа для Priority
	invalidBody := []byte(`{"title":"test","priority":"invalid"}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// 4. PUT /task/:id — успешное обновление
func TestEditTaskSuccess(t *testing.T) {
	clearList()
	router := setupRouter()

	// Сначала создаём задачу
	newTask := task.Task{
		Title:       "Original",
		Description: "Original desc",
		Priority:    task.Low,
		Status:      task.New,
	}
	body, _ := json.Marshal(newTask)
	reqPost, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	reqPost.Header.Set("Content-Type", "application/json")
	wPost := httptest.NewRecorder()
	router.ServeHTTP(wPost, reqPost)
	if wPost.Code != http.StatusCreated {
		t.Fatal("Failed to create task for update test")
	}
	var created task.Task
	if err := json.Unmarshal(wPost.Body.Bytes(), &created); err != nil {
		t.Fatal("Failed to parse created task:", err)
	}

	// Теперь обновляем
	updated := task.Task{
		Title:       "Updated",
		Description: "Updated desc",
		Priority:    task.High,
		Status:      task.Done,
	}
	updateBody, _ := json.Marshal(updated)
	reqPut, _ := http.NewRequest("PUT", "/task/"+created.Id.String(), bytes.NewBuffer(updateBody))
	reqPut.Header.Set("Content-Type", "application/json")
	wPut := httptest.NewRecorder()
	router.ServeHTTP(wPut, reqPut)

	if wPut.Code != http.StatusCreated { // в вашем коде используется 201 Created
		t.Errorf("Expected 201, got %d", wPut.Code)
	}

	var result task.Task
	if err := json.Unmarshal(wPut.Body.Bytes(), &result); err != nil {
		t.Fatal("Failed to parse updated task:", err)
	}
	if result.Title != "Updated" {
		t.Errorf("Title not updated: %s", result.Title)
	}
	if result.Id != created.Id {
		t.Errorf("ID changed: expected %v, got %v", created.Id, result.Id)
	}
}

// 5. PUT с невалидным UUID
func TestEditTaskInvalidUUID(t *testing.T) {
	clearList()
	router := setupRouter()

	req, _ := http.NewRequest("PUT", "/task/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// 6. PUT с валидным UUID, но задача не найдена
func TestEditTaskNotFound(t *testing.T) {
	clearList()
	router := setupRouter()

	fakeID := uuid.New()
	body := []byte(`{"title":"test"}`)
	req, _ := http.NewRequest("PUT", "/task/"+fakeID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Ваш код возвращает 400 при NotFound, хотя логичнее 404, но проверяем что есть
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
