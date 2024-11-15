package main

import (
	"log"

	"blog_server/cache"
	"blog_server/controller"
	"blog_server/db"
	"blog_server/server"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Faied to init database: %v", err)
	}

	defer db.DB.Close()

	cache.InitCache(32)
	server.Sv = server.InitSv()

	server.Sv.Gin.GET("/article/:name", controller.GetArticleFromName)
	server.Sv.Gin.POST("/webhook/update-knowledge", controller.GithubWebhook)
	server.Sv.Gin.Run(":3000")
}
