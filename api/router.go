package api

import "github.com/gin-gonic/gin"

// serviceStatus permit to determine if the service is up or not.
// @summary      GET /service-status
// @description  Simple route to determine if the service is up or not.
// @tags         admin
// @produce      json
// @Success      200  {object}  map[string]string
// @Router       /service-status [get]
func serviceStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func NewRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/service-status", serviceStatus)

	router.GET("/task/:id", returnTask)
	router.GET("/task", returnTasks)
	router.POST("/task", createTask)
	router.PATCH("/task", updateTask)
	router.DELETE("/task/:id", deleteTask)

	router.GET("/note/:id", returnNote)
	router.GET("/note", returnNotes)
	router.POST("/note", createNote)
	router.PATCH("/note", updateNote)
	router.DELETE("/note/:id", deleteNote)

	router.GET("/project/:id", returnProject)
	router.GET("/project", returnProjects)
	router.POST("/project", createProject)
	router.PATCH("/project", updateProject)
	router.DELETE("/project/:id", deleteProject)
	router.GET("/project/:id/task", returnTasksByProjectId)
	router.GET("/project/:id/note", returnNotesByProjectId)

	return router
}
