package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Aerospace-prog/scheme-manager/config"
	"github.com/Aerospace-prog/scheme-manager/models"
	"github.com/Aerospace-prog/scheme-manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect DB
	config.ConnectDatabase()

	// Migrate
	err := config.DB.AutoMigrate(&models.Scheme{}, &models.Application{})
	if err != nil {
		log.Fatal("❌ Failed DB migration:", err)
	}
	fmt.Println("✅ Database migrated!")

	// Register all routes
	routes.RegisterRoutes(r)

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Scheme Manager API is running!"})
	})

	port := "8080"
	fmt.Println("Server is running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}