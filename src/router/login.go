package router

import "github.com/gin-gonic/gin"

// post : /login
func Login(c *gin.Context) {
	println("Login router!")
	c.JSON(200, gin.H{
		"message": "ping",
	})
}
