package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Gin    *gin.Engine
	Status int
}

var Sv *Server

func TemplateMiddleware(tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("tmpl", tmpl)
		c.Next()
	}
}

func InitSv() *Server {
	engine := gin.Default()

	engine.Static("/assets", "./assets")
	engine.LoadHTMLGlob("templates/*")
	tmpl := template.Must(template.New("").ParseGlob("templates/*.tmpl"))
	engine.Use(TemplateMiddleware(tmpl))

	engine.GET("/hc", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Service working.",
		})
	})
	return &Server{
		Gin:    engine,
		Status: 0,
	}
}

func StartServerMaintenance(s *Server) {
	s.Status = 1
}

func EndServerMaintenance(s *Server) {
	s.Status = 0
}
