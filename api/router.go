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

	global := router.Group("/api/v0")
	{
		global.GET("/service-status", serviceStatus)

		taskRoute := global.Group("/task")
		{
			taskRoute.GET("/:id", returnTask)
			taskRoute.GET("", returnTasks)
			taskRoute.POST("", createTask)
			taskRoute.PATCH("", updateTask)
			taskRoute.DELETE("/:id", deleteTask)
		}

		noteRoute := global.Group("/note")
		{
			noteRoute.GET("/:id", returnNote)
			noteRoute.GET("", returnNotes)
			noteRoute.POST("", createNote)
			noteRoute.PATCH("", updateNote)
			noteRoute.DELETE("/:id", deleteNote)
		}

		projectRoute := global.Group("/project")
		{
			projectRoute.GET("/:id", returnProject)
			projectRoute.GET("", returnProjects)
			projectRoute.POST("", createProject)
			projectRoute.PATCH("", updateProject)
			projectRoute.DELETE("/:id", deleteProject)
			projectRoute.GET("/:id/task", returnTasksByProjectId)
			projectRoute.GET("/:id/note", returnNotesByProjectId)
		}
	}

	return router
}
