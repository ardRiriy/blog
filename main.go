package main

import (
	"log"

	"blog_server/cache"
	"blog_server/controller"
	"blog_server/db"
	"blog_server/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
	}

	if err := db.InitDB(); err != nil {
		log.Fatalf("Faied to init database: %v", err)
	}

	defer db.DB.Close()

	// initialize.InitializeDB()

	cache.InitCache(32)
	server.Sv = server.InitSv()

	server.Sv.Gin.GET("/", controller.GetIndex)
	server.Sv.Gin.GET("/index", controller.GetIndex)
	server.Sv.Gin.GET("/hc", controller.GetHealthCheck)
	server.Sv.Gin.GET("/article/:name", controller.GetArticleFromName)
	server.Sv.Gin.POST("/webhook/update-knowledge", controller.GithubWebhook)
	server.Sv.Gin.Run(":3000")
}
