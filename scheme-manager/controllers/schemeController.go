package controllers

import (
	"net/http"

	"github.com/Aerospace-prog/scheme-manager/config"
	"github.com/Aerospace-prog/scheme-manager/models"

	"github.com/gin-gonic/gin"
)

// =================== 1. Create Scheme ===================
func PushScheme(c *gin.Context) {
	var scheme models.Scheme

	if err := c.ShouldBindJSON(&scheme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
		return
	}

	if scheme.Name == "" || scheme.Description == "" || scheme.Eligibility == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Description and Eligibility are required"})
		return
	}

	if err := config.DB.Create(&scheme).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scheme"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Scheme created", "scheme": scheme})
}

// =================== 2. Get All Applications ===================
func GetApplications(c *gin.Context) {
	status := c.Query("status")
	var applications []models.Application
	query := config.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// =================== 3. Update Application Status ===================
func UpdateApplicationStatus(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status field is required"})
		return
	}

	validStatuses := map[string]bool{
		"pending":  true,
		"approved": true,
		"rejected": true,
	}

	if !validStatuses[payload.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	var app models.Application
	if err := config.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	app.Status = payload.Status
	if err := config.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "application": app})
}

// =================== 4. Get Application Status By its ID===================

func GetApplicationByID(c *gin.Context) {
    id := c.Param("id")
    var application models.Application

    // Find application by ID
    if err := config.DB.First(&application, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "application": application,
    })
}