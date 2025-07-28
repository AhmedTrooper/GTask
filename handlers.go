package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getTasks(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	taskList := []Task{}
	for _, task := range tasks {
		taskList = append(taskList, task)
	}
	c.JSON(http.StatusOK, taskList)
}

func getTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	mu.Lock()
	defer mu.Unlock()
	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context) {
	var newTask Task
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	newTask.ID = nextID
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()
	tasks[nextID] = newTask
	nextID++
	mu.Unlock()
	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedTask Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()
	existingTask, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	updatedTask.ID = id
	updatedTask.CreatedAt = existingTask.CreatedAt
	updatedTask.UpdatedAt = time.Now()
	tasks[id] = updatedTask
	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	mu.Lock()
	defer mu.Unlock()
	if _, exists := tasks[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	delete(tasks, id)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
