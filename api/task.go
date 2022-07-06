package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/model"
)

func postTask(c *gin.Context) {

	var newTask model.Task

	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println(err)
		return
	}

	newTask.CreateTask()
	c.IndentedJSON(http.StatusCreated, newTask)
}

func returnTask(c *gin.Context) {
	id := c.Param("id")

	task, err := model.ReturnTask(id)

	if err != nil {
		c.Error(err)
	}

	c.IndentedJSON(http.StatusOK, task)
}

func returnTasks(c *gin.Context) {

	tasks, err := model.ReturnTasks()

	if err != nil {
		c.Error(err)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func updateTask(c *gin.Context) {
	var updatedTask model.Task

	if err := c.BindJSON(&updatedTask); err != nil {
		fmt.Println(err)
		return
	}

	updatedTask.UpdateTask()
	c.IndentedJSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	taskToDelete, err := model.ReturnTask(id)

	if err != nil {
		c.Error(err)
	}

	taskToDelete.DeleteTask()
	c.IndentedJSON(http.StatusOK, taskToDelete)
}
