package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/model"
)

func postTask(c *gin.Context) {
	loggeur := log.Default()
	var newTask model.Task

	if err := c.BindJSON(&newTask); err != nil {
		fmt.Println(err)
		return
	}

	err := newTask.CreateTask()

	if err != nil {
		loggeur.Printf("an error occured during the insertion.")
		msg, _ := json.Marshal(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msg)
		return
	}

	loggeur.Printf("everything went well.")
	c.IndentedJSON(http.StatusCreated, newTask)
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
