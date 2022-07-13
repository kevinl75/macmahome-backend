package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/service-status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "OK",
		})
	})

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
