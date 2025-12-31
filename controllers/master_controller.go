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

// UpdateFaculty
func UpdateFaculty(c *gin.Context) {
	id := c.Param("id")
	var faculty models.Faculty

	if err := models.DB.First(&faculty, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	}

	var input models.Faculty
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	faculty.Code = input.Code
	faculty.Name = input.Name
	faculty.Type = input.Type

	models.DB.Save(&faculty)
	c.JSON(http.StatusOK, gin.H{"data": faculty})
}

// DeleteFaculty
func DeleteFaculty(c *gin.Context) {
	id := c.Param("id")
	var faculty models.Faculty

	if err := models.DB.First(&faculty, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Faculty not found"})
		return
	}

	// Hapus (hati-hati, ini bisa menghapus prodi/department terkait jika on delete cascade)
	if err := models.DB.Delete(&faculty).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete faculty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Faculty deleted successfully"})
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

// DeleteDepartment
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	var dept models.Department

	if err := models.DB.First(&dept, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	if err := models.DB.Delete(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
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

// UpdateBuilding: Mengubah data gedung
func UpdateBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building

	if err := models.DB.First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		TotalFloors int    `json:"total_floors"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	building.Code = input.Code
	building.Name = input.Name
	building.TotalFloors = input.TotalFloors

	models.DB.Save(&building)
	c.JSON(http.StatusOK, gin.H{"data": building})
}

// DeleteBuilding: Menghapus gedung
func DeleteBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building

	if err := models.DB.First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}

	// Hapus gedung (hati-hati jika ada relasi rooms, pastikan constraint database ON DELETE CASCADE)
	if err := models.DB.Delete(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete building"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Building deleted successfully"})
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

// UpdateCategory: Mengubah data kategori
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.AssetCategory

	// 1. Cari Data
	if err := models.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// 2. Validasi Input
	var input struct {
		KodeBarang  string `json:"kode_barang"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update & Simpan
	category.KodeBarang = input.KodeBarang
	category.Name = input.Name
	category.Description = input.Description

	if err := models.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory: Menghapus kategori
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.AssetCategory

	// 1. Cari Data
	if err := models.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// 2. Hapus
	if err := models.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
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

// UpdateRoom: Mengubah data ruangan
func UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room

	if err := models.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	var input struct {
		BuildingID uint   `json:"building_id"`
		RoomNumber string `json:"room_number"`
		Name       string `json:"name"`
		Floor      int    `json:"floor"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room.BuildingID = input.BuildingID
	room.RoomNumber = input.RoomNumber
	room.Name = input.Name
	room.Floor = input.Floor

	models.DB.Save(&room)
	c.JSON(http.StatusOK, gin.H{"data": room})
}

// DeleteRoom: Menghapus ruangan (Ini yang menyebabkan error 404 Anda)
func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room

	if err := models.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if err := models.DB.Delete(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
