package main

import (
	"net/http"
	"sync"
	task "to-do/Task"
	todolist "to-do/ToDoList"

	"github.com/gin-gonic/gin"
)

var (
	list = todolist.TodoList{}
	mu   sync.RWMutex // защищает albums от конкурентных записей
)

func main() {
	router := gin.Default()
	router.GET("/list", getList)
	router.POST("/task", addTask)

	router.Run("localhost:8080")
}

func getList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, list)
}

func addTask(c *gin.Context) {
	var newTask task.Task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}
	mu.Lock()
	list = append(list, newTask)
	defer mu.Unlock()
	c.IndentedJSON(http.StatusCreated, list)
}
