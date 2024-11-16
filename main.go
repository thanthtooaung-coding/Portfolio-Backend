package main

import (
	"portfolio-backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/projects", handlers.GetProjects)
	r.POST("/projects", handlers.AddProject)
	r.DELETE("/projects/:title", handlers.DeleteProject)
	r.PUT("/projects/:title", handlers.UpdateProject)

	r.Run(":8080")
}
