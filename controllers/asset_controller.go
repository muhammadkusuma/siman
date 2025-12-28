package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

// GetAssets mengambil semua aset beserta relasinya
func GetAssets(c *gin.Context) {
	var assets []models.Asset
	// Preload digunakan untuk memuat data dari tabel relasi (Category, Department, Room)
	models.DB.Preload("Category").Preload("Department").Preload("Room").Find(&assets)

	c.JSON(http.StatusOK, gin.H{"data": assets})
}

// GetAssetByID mengambil detail satu aset
func GetAssetByID(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.Preload("Category").Preload("Department").Preload("Room").Preload("CreatedBy").Preload("UpdatedBy").First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// CreateAsset menambahkan aset baru
func CreateAsset(c *gin.Context) {
	var input models.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi sederhana: Pastikan kode inventaris unik (opsional, karena DB sudah unique)
	// Set CreatedAt otomatis oleh GORM
	
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- TAMBAHAN: CATAT KE AUDIT LOG ---
    // Pastikan Anda punya mekanisme untuk tahu UserID yang sedang login. 
    // Untuk contoh ini kita hardcode ID user = 1 (Admin) jika belum pakai JWT Middleware.
    currentUserID := uint(1) 
    
    // Panggil Helper dari Audit Controller tadi
    // Parameter: Context, UserID, Action, TableName, RecordID, Changes/Keterangan
    RecordLog(c, currentUserID, "CREATE", "assets", input.ID, "Menambahkan aset baru: "+input.Name)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

// UpdateAsset mengupdate data aset
func UpdateAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	var input models.Asset
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&asset).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// DeleteAsset menghapus aset
func DeleteAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	models.DB.Delete(&asset)
	c.JSON(http.StatusOK, gin.H{"data": true})
}