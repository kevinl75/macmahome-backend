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

// createProject permit to handle a POST request to create a new Project.
// @summary      POST /project
// @description  Permit to create a new project entity from a JSON object. Return the created entity. If
//               the JSON provided in request data included project tasks and notes, they will not be
//               inserted. Insertion of linked notes and tasks should be done separately.
// @tags         project
// @accept       json
// @produce      json
// @param        project  body  model.Project  true  "New project to create"
// @Success      201  {object}  model.Project
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project [post]
func createProject(c *gin.Context) {
	var newProject model.Project

	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := newProject.CreateProject()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newProject)
}

// returnProject permit to return a project based on its id
// @summary      GET /project/:id
// @description  Permit to return a Project entity based on its id. Return
//               also the tasks and the notes linked to the project.
// @tags         project
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the project to return"
// @Success      200  {object}  model.Project
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to return not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project/:id [get]
func returnProject(c *gin.Context) {

	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := model.ReturnProject(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if project.IsEqual(model.Project{}) {
		errorMsg := fmt.Errorf("no project entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, project)
}

// returnProjects permit to return all the Project entities
// @summary      GET /project
// @description  Permit to return all the Project entities.
// @tags         project
// @accept       mpfd
// @produce      json
// @Success      200  {object}  []model.Project
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project [get]
func returnProjects(c *gin.Context) {

	projects, err := model.ReturnProjects()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// updateProject permit to handle a PATCH request to update a Project.
// @summary      PATCH /project
// @description  Permit to update a project entity with a JSON object. Return the updated entity
// @tags         project
// @accept       json
// @produce      json
// @param        project  body  model.Project  true  "Updated project entity"
// @Success      200  {object}  model.Project
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project [patch]
func updateProject(c *gin.Context) {
	var updatedProject model.Project

	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := updatedProject.UpdateProject()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := fmt.Errorf("no project entity with id %d", updatedProject.ProjectId).Error()
			c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

// deleteProject permit to delete a project based on its id
// @summary      DELETE /project/:id
// @description  Permit to delete a project based on its id. Return the deleted entity.
// @tags         project
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the project to delete"
// @Success      200  {object}  model.Project
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project/:id [delete]
func deleteProject(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectToDelete, err := model.ReturnProject(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if projectToDelete.IsEqual(model.Project{}) {
		errorMsg := fmt.Errorf("no project entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	projectToDelete.DeleteProject()
	c.JSON(http.StatusOK, projectToDelete)
}
