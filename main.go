package main

import (
	"muhammadkusuma/siman/controllers"
	"muhammadkusuma/siman/middlewares"
	"muhammadkusuma/siman/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	models.ConnectDatabase()

	// --- Static File Server ---
	// Ini membuat file di dalam folder "./uploads" bisa diakses via URL "http://localhost:3000/uploads/..."
	router.Static("/uploads", "./uploads")

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
		// ... (kode route lainnya tetap sama) ...
		api.GET("/faculties", controllers.GetFaculties)
		api.POST("/faculties", controllers.CreateFaculty)
		api.GET("/departments", controllers.GetDepartments)
		api.POST("/departments", controllers.CreateDepartment)
		api.GET("/buildings", controllers.GetBuildings)
		api.POST("/buildings", controllers.CreateBuilding)
		api.POST("/rooms", controllers.CreateRoom)
		api.GET("/categories", controllers.GetCategories)
		api.POST("/categories", controllers.CreateCategory)

		// --- Assets ---
		api.GET("/assets", controllers.GetAssets)
		api.GET("/assets/:id", controllers.GetAssetByID)
		
		// Update logic Create dan Update di controller (lihat file berikutnya)
		api.POST("/assets", controllers.CreateAsset) 
		api.PUT("/assets/:id", controllers.UpdateAsset) 
		
		api.DELETE("/assets/:id", controllers.DeleteAsset)

		// ... (kode transaction & audit log tetap sama) ...
		api.POST("/mutations", controllers.CreateMutation)
		api.GET("/mutations", controllers.GetMutations)
		api.POST("/maintenances", controllers.CreateMaintenance)
		api.GET("/maintenances", controllers.GetMaintenances)
		api.GET("/audit-logs", controllers.GetAuditLogs)
		api.GET("/audit-logs/:id", controllers.GetAuditLogByID)
	}

	router.Run(":3000")
}