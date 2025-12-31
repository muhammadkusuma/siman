package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- FACULTY ---
func GetFaculties(c *gin.Context) {
	var faculties []models.Faculty
	models.DB.Preload("Departments").Find(&faculties)
	c.JSON(http.StatusOK, gin.H{"data": faculties})
}

func CreateFaculty(c *gin.Context) {
	var input models.Faculty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// --- DEPARTMENT ---
func GetDepartments(c *gin.Context) {
	var depts []models.Department
	models.DB.Preload("Faculty").Find(&depts)
	c.JSON(http.StatusOK, gin.H{"data": depts})
}

func CreateDepartment(c *gin.Context) {
	var input models.Department
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// --- BUILDING & ROOM ---
func GetBuildings(c *gin.Context) {
	var buildings []models.Building
	models.DB.Preload("Rooms").Find(&buildings)
	c.JSON(http.StatusOK, gin.H{"data": buildings})
}

func CreateBuilding(c *gin.Context) {
	var input models.Building
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func CreateRoom(c *gin.Context) {
	var input models.Room
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// --- ASSET CATEGORY ---
func GetCategories(c *gin.Context) {
	var cats []models.AssetCategory
	models.DB.Find(&cats)
	c.JSON(http.StatusOK, gin.H{"data": cats})
}

func CreateCategory(c *gin.Context) {
	var input models.AssetCategory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&input)
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// GetRoomsByBuildingID mengambil daftar ruangan berdasarkan ID Gedung
func GetRoomsByBuildingID(c *gin.Context) {
	buildingID := c.Param("buildingID")
	var rooms []models.Room

	if err := models.DB.Where("building_id = ?", buildingID).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rooms not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rooms})
}
