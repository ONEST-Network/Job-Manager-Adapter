package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Aerospace-prog/scheme-manager/database"
	"github.com/Aerospace-prog/scheme-manager/models"
)

func CreateApplication(c *gin.Context) {
	var application models.Application

	// Bind the JSON payload to the struct
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// Check if the referenced scheme exists
	var scheme models.Scheme
	if err := database.DB.First(&scheme, application.SchemeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Scheme not found"})
		return
	}

	// Default status
	if application.Status == "" {
		application.Status = "Pending"
	}

	// Save the application to DB
	if err := database.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create application", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}