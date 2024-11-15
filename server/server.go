package server

import (
	"html/template"

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
