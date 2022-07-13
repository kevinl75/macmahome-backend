package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/model"
	"gorm.io/gorm"
)

// createTask permit to handle a POST request to create a new Task.
// @summary      POST /task
// @description  Permit to create a new task entity from a JSON object. Task could be linked
//               to a project or not. Return the created entity
// @tags         task
// @accept       json
// @produce      json
// @param        task  body  model.Task  true  "New task to create"
// @Success      201  {object}  model.Task
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /task [post]
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

// returnTask permit to return a task based on its id
// @summary      GET /task/:id
// @description  Permit to return a Task entity based on its id.
// @tags         task
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the task to return"
// @Success      200  {object}  model.Task
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to return not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /task/:id [get]
func returnTask(c *gin.Context) {

	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := model.ReturnTask(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (task == model.Task{}) {
		errorMsg := fmt.Errorf("no task entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, task)
}

// returnTasks permit to return all the Task entities
// @summary      GET /task
// @description  Permit to return all the Task entities.
// @tags         task
// @accept       mpfd
// @produce      json
// @Success      200  {object}  []model.Task
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /task [get]
func returnTasks(c *gin.Context) {

	tasks, err := model.ReturnTasks()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// returnTasksByProjectId permit to return all the Task entities of a project with the project id
// @summary      GET /project/:id/task
// @description  Permit to return all the Task entities of a project with the project id.
// @tags         task
// @accept       mpfd
// @produce      json
// @Success      200  {object}  []model.Task
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project/:id/task [get]
func returnTasksByProjectId(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := model.ReturnTasksByProjectId(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// updateTask permit to handle a PATCH request to update a Task.
// @summary      PATCH /task
// @description  Permit to update a task entity with a JSON object. Return the updated entity
// @tags         task
// @accept       json
// @produce      json
// @param        task  body  model.Task  true  "Updated task entity"
// @Success      200  {object}  model.Task
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /task [patch]
func updateTask(c *gin.Context) {
	var updatedTask model.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := updatedTask.UpdateTask()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := fmt.Errorf("no task entity with id %d", updatedTask.TaskId).Error()
			c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// deleteTask permit to delete a task based on its id
// @summary      DELETE /task/:id
// @description  Permit to delete a task based on its id. Return the deleted entity.
// @tags         task
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the task to delete"
// @Success      200  {object}  model.Task
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /task/:id [delete]
func deleteTask(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskToDelete, err := model.ReturnTask(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (taskToDelete == model.Task{}) {
		errorMsg := fmt.Errorf("no task entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	taskToDelete.DeleteTask()
	c.JSON(http.StatusOK, taskToDelete)
}
