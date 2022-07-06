package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/service-status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	router.GET("/task/:id", returnTask)
	router.GET("/task", returnTasks)
	router.POST("/task", postTask)
	router.PATCH("/task", updateTask)
	router.DELETE("/task/:id", deleteTask)

	return router
}
