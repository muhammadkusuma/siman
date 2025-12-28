package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Ambil User ID dari Token untuk field CreatedBy & UpdatedBy
	// Kita gunakan helper c.MustGet karena sudah dipastikan ada oleh Middleware
	if userID, exists := c.Get("userID"); exists {
		input.CreatedByID = userID.(uint)
		input.UpdatedByID = userID.(uint)
	}

	// Validasi & Simpan DB
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- LOG AUDIT OTOMATIS ---
	// Perbaikan: Tidak perlu kirim User ID lagi, cukup Context
	RecordLog(c, "CREATE", "assets", input.ID, "Menambahkan aset baru: "+input.Name)

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

	// Update field UpdatedBy dari token
	if userID, exists := c.Get("userID"); exists {
		input.UpdatedByID = userID.(uint)
	}

	models.DB.Model(&asset).Updates(input)

	// --- LOG AUDIT UPDATE ---
	RecordLog(c, "UPDATE", "assets", asset.ID, "Mengupdate data aset")

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

	// --- LOG AUDIT DELETE ---
	RecordLog(c, "DELETE", "assets", asset.ID, "Menghapus aset: "+asset.Name)

	c.JSON(http.StatusOK, gin.H{"data": true})
}