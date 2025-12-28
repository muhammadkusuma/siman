package main

import (
	"muhammadkusuma/siman/controllers" // Import package controllers
	"muhammadkusuma/siman/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	models.ConnectDatabase()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "SIMAN API is Running!"})
	})

	// --- Auth Routes ---
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	// --- Master Data Routes ---
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		// Faculties & Depts
		api.GET("/faculties", controllers.GetFaculties)
		api.POST("/faculties", controllers.CreateFaculty)
		api.GET("/departments", controllers.GetDepartments)
		api.POST("/departments", controllers.CreateDepartment)

		// Locations
		api.GET("/buildings", controllers.GetBuildings)
		api.POST("/buildings", controllers.CreateBuilding)
		api.POST("/rooms", controllers.CreateRoom)

		// Categories
		api.GET("/categories", controllers.GetCategories)
		api.POST("/categories", controllers.CreateCategory)

		// --- Assets ---
		api.GET("/assets", controllers.GetAssets)
		api.GET("/assets/:id", controllers.GetAssetByID)
		api.POST("/assets", controllers.CreateAsset)
		api.PUT("/assets/:id", controllers.UpdateAsset)
		api.DELETE("/assets/:id", controllers.DeleteAsset)

		// --- Transactions ---
		api.POST("/mutations", controllers.CreateMutation) // Pindah Aset
		api.GET("/mutations", controllers.GetMutations)
		
		api.POST("/maintenances", controllers.CreateMaintenance) // Lapor Rusak
		api.GET("/maintenances", controllers.GetMaintenances)

		// --- Audit Logs (Read Only) ---
		api.GET("/audit-logs", controllers.GetAuditLogs)
		api.GET("/audit-logs/:id", controllers.GetAuditLogByID)
	}

	router.Run(":3000")
}