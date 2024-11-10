package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/hc", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Service working.",
		})
	})
	engine.Run(":3000")
}
