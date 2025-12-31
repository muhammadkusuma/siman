package main

import (
	"muhammadkusuma/siman/controllers"
	"muhammadkusuma/siman/middlewares"
	"muhammadkusuma/siman/models"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"time"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// 1. LOAD ENVIRONMENT VARIABLES DULU
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	
	// 2. Koneksi Database sekarang akan baca dari env
	models.ConnectDatabase()

    // --- SETUP CORS ---
    // Ini mengizinkan semua domain mengakses API (Bisa diperketat nanti)
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Atau ganti "http://localhost:5173"
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
	// router := gin.Default()
	// models.ConnectDatabase()

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
		// --- BARU: Dashboard & User ---
		api.GET("/dashboard", controllers.GetDashboardStats) // Endpoint Dashboard
		api.GET("/profile", controllers.GetProfile)          // Endpoint Profil Saya
		api.GET("/users", controllers.GetAllUsers)           // Endpoint List User

		api.GET("/faculties", controllers.GetFaculties)
		api.POST("/faculties", controllers.CreateFaculty)
		api.GET("/departments", controllers.GetDepartments)
		api.POST("/departments", controllers.CreateDepartment)
		api.GET("/buildings", controllers.GetBuildings)
		api.POST("/buildings", controllers.CreateBuilding)

		// --- BARU: Filter Ruangan by Gedung ---
		api.GET("/buildings/:buildingID/rooms", controllers.GetRoomsByBuildingID)
		
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
