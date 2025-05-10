package routes

import (
	"github.com/ONEST-Network/scheme-manager-adapter/api/handlers"
	"github.com/ONEST-Network/scheme-manager-adapter/api/middleware"
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRouter sets up the API routes
func SetupRouter(db *pgxpool.Pool, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Apply middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())

	// Create base handler
	baseHandler := handlers.NewBaseHandler(db, cfg)

	// Create specific handlers
	orgHandler := handlers.NewOrganizationHandler(baseHandler)
	schemeHandler := handlers.NewSchemeHandler(baseHandler)
	appHandler := handlers.NewApplicationHandler(baseHandler)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Organization routes
		orgs := v1.Group("/organizations")
		{
			orgs.POST("", orgHandler.Create)
			orgs.GET("", orgHandler.List)
			orgs.GET("/:id", orgHandler.GetByID)
			orgs.GET("/api-key/:api_key", orgHandler.GetByAPIKey)
			orgs.PUT("/:id", orgHandler.Update)
			orgs.DELETE("/:id", orgHandler.Delete)

			// Scheme routes
			schemes := orgs.Group("/:organization_id/schemes")
			{
				schemes.POST("", schemeHandler.Create)
				schemes.GET("", schemeHandler.ListByOrganization)
				schemes.GET("/:scheme_id", schemeHandler.GetBySchemeID)
			}
		}

		// Standalone scheme routes
		schemes := v1.Group("/schemes")
		{
			schemes.GET("/:id", schemeHandler.GetBySchemeID)
			schemes.PUT("/:id/status", schemeHandler.UpdateStatus)
			schemes.DELETE("/:id", schemeHandler.Delete)

			// Application routes
			apps := schemes.Group("/:scheme_id/applications")
			{
				apps.POST("", appHandler.Create)
				apps.GET("", appHandler.ListByScheme)
			}
		}

		// Standalone application routes
		apps := v1.Group("/applications")
		{
			apps.GET("/:id", appHandler.GetByID)
			apps.PUT("/:id/status", appHandler.UpdateStatus)
			apps.DELETE("/:id", appHandler.Delete)
		}
	}

	return r
}
