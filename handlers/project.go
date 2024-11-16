package handlers

import (
	"encoding/json"
	"net/http"

	"portfolio-backend/models"
	"portfolio-backend/utils"

	"github.com/gin-gonic/gin"
)

var (
	s3Client = utils.NewS3Client("thanthtooaung-portfolio")
	s3Key    = "json/projects.json"
)

func GetProjects(c *gin.Context) {
	data, err := s3Client.GetObject(s3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects from S3", "details": err.Error()})
		return
	}

	var projects []models.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse projects JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func AddProject(c *gin.Context) {
	var newProject models.Project
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	data, err := s3Client.GetObject(s3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing projects"})
		return
	}

	var projects []models.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format"})
		return
	}

	projects = append(projects, newProject)
	updatedData, _ := json.Marshal(projects)
	if err := s3Client.PutObject(s3Key, updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update projects"})
		return
	}

	c.JSON(http.StatusCreated, newProject)
}

func DeleteProject(c *gin.Context) {
	title := c.Param("title")

	data, err := s3Client.GetObject(s3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing projects"})
		return
	}

	var projects []models.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format"})
		return
	}

	var updatedProjects []models.Project
	for _, project := range projects {
		if project.Title != title {
			updatedProjects = append(updatedProjects, project)
		}
	}

	updatedData, _ := json.Marshal(updatedProjects)
	if err := s3Client.PutObject(s3Key, updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}

func UpdateProject(c *gin.Context) {
	title := c.Param("title")

	var updatedProject models.Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	data, err := s3Client.GetObject(s3Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing projects", "details": err.Error()})
		return
	}

	var projects []models.Project
	if err := json.Unmarshal(data, &projects); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format"})
		return
	}

	projectUpdated := false
	for i, project := range projects {
		if project.Title == title {
			projects[i] = updatedProject
			projectUpdated = true
			break
		}
	}

	if !projectUpdated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	updatedData, err := json.Marshal(projects)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process updated projects", "details": err.Error()})
		return
	}

	if err := s3Client.PutObject(s3Key, updatedData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": updatedProject})
}
