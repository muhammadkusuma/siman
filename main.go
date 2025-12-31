package main

import (
	"muhammadkusuma/siman/controllers"
	"muhammadkusuma/siman/middlewares"
	"muhammadkusuma/siman/models"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// 1. LOAD ENVIRONMENT VARIABLES DULU
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
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

	// --- ROOT ENDPOINT (HEALTH CHECK) ---
	// Diubah agar mengecek koneksi database
	router.GET("/", func(c *gin.Context) {
		// Default status
		dbStatus := "Connected"
		statusCode := http.StatusOK

		// Cek apakah object DB sudah ada
		if models.DB == nil {
			dbStatus = "Database Not Initialized"
			statusCode = http.StatusInternalServerError
		} else {
			// Cek ping ke database sql asli
			sqlDB, err := models.DB.DB()
			if err != nil {
				dbStatus = "Database Driver Error"
				statusCode = http.StatusInternalServerError
			} else if err := sqlDB.Ping(); err != nil {
				// Ping gagal berarti server DB mati atau jaringan putus
				dbStatus = "Connection Failed: " + err.Error()
				statusCode = http.StatusServiceUnavailable
			}
		}

		c.JSON(statusCode, gin.H{
			"app_name": "SIMAN API",
			"version":  "1.0.0",
			"status":   "Running",
			"database": dbStatus,
			"time":     time.Now(),
		})
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
		api.GET("/users", controllers.GetAllUsers)           // Endpoint List User

		// --- PROFILE ROUTES (PERBAIKAN DISINI) ---
		api.GET("/profile", controllers.GetProfile)
		api.PUT("/profile", controllers.UpdateProfile) // <--- Tambahkan ini
		api.PUT("/change-password", controllers.ChangePassword)

		api.GET("/faculties", controllers.GetFaculties)
		api.POST("/faculties", controllers.CreateFaculty)
		api.PUT("/faculties/:id", controllers.UpdateFaculty) // <--- TAMBAH
		api.DELETE("/faculties/:id", controllers.DeleteFaculty)

		api.GET("/departments", controllers.GetDepartments)
		api.POST("/departments", controllers.CreateDepartment)
		api.DELETE("/departments/:id", controllers.DeleteDepartment)

		api.GET("/buildings", controllers.GetBuildings)
		api.POST("/buildings", controllers.CreateBuilding)
		api.PUT("/buildings/:id", controllers.UpdateBuilding)    // <--- Tambah ini
		api.DELETE("/buildings/:id", controllers.DeleteBuilding) // <--- Tambah ini
		api.GET("/buildings/:buildingID/rooms", controllers.GetRoomsByBuildingID)

		api.POST("/rooms", controllers.CreateRoom)
		api.PUT("/rooms/:id", controllers.UpdateRoom) // <--- Tambah ini
		api.DELETE("/rooms/:id", controllers.DeleteRoom)

		api.GET("/categories", controllers.GetCategories)
		api.POST("/categories", controllers.CreateCategory)
		api.PUT("/categories/:id", controllers.UpdateCategory)    // Untuk Update
		api.DELETE("/categories/:id", controllers.DeleteCategory) // Untuk Delete

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
		api.PUT("/maintenances/:id", controllers.UpdateMaintenance) 
		api.DELETE("/maintenances/:id", controllers.DeleteMaintenance)

		api.GET("/audit-logs", controllers.GetAuditLogs)
		api.GET("/audit-logs/:id", controllers.GetAuditLogByID)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run(":" + port)
}
