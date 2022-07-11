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

func returnProjects(c *gin.Context) {

	projects, err := model.ReturnProjects()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

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
