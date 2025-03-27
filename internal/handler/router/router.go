package router

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	return r
}
