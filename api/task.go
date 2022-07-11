package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/model"
)

func createTask(c *gin.Context) {
	var newTask model.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := newTask.CreateTask()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

func returnTask(c *gin.Context) {
	id := c.Param("id")

	task := model.ReturnTask(id)

	if (task == model.Task{}) {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.IndentedJSON(http.StatusOK, task)
}

func returnTasks(c *gin.Context) {

	tasks, err := model.ReturnTasks()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func updateTask(c *gin.Context) {
	var updatedTask model.Task

	if err := c.BindJSON(&updatedTask); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err := updatedTask.UpdateTask()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	taskToDelete := model.ReturnTask(id)

	if (taskToDelete == model.Task{}) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		taskToDelete.DeleteTask()
		c.IndentedJSON(http.StatusOK, taskToDelete)
	}
}
