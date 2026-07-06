package main

import (
	"log"
	"net/http"
	"sync"
	task "to-do/Task"
	todolist "to-do/ToDoList"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	list = todolist.TodoList{}
	mu   sync.RWMutex // защищает list от конкурентных записей
)

func main() {
	router := gin.Default()
	router.GET("/list", getList)
	router.POST("/task", addTask)
	router.PUT("/task/:id", editTask)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}

func getList(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, list)
}

func addTask(c *gin.Context) {
	var newTask task.Task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}
	mu.Lock()
	newTask.Id = uuid.New()
	list = append(list, newTask)
	mu.Unlock()
	c.JSON(http.StatusCreated, newTask)
}

func editTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	index := -1

	for i := range list {
		if list[i].Id == id {
			index = i
		}
	}

	if index < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NotFound"})
		return
	}

	var updatedTask task.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "InvalidObject"})
		return
	}
	mu.Lock()
	list[index].Title = updatedTask.Title
	list[index].Description = updatedTask.Description
	list[index].Priority = updatedTask.Priority
	list[index].Status = updatedTask.Status
	mu.Unlock()
	c.JSON(http.StatusCreated, list[index])
}
