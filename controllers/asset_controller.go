package controllers

import (
	"fmt"
	"muhammadkusuma/siman/models"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Update fungsi GetAssets agar bisa search
func GetAssets(c *gin.Context) {
    var assets []models.Asset
    
    // Ambil parameter query dari URL: ?search=laptop&status=Baik
    search := c.Query("search")
    status := c.Query("status")

    // Mulai Query
    query := models.DB.Preload("Category").Preload("Department").Preload("Room")

    // Filter by Name (Search)
    if search != "" {
        query = query.Where("name LIKE ?", "%"+search+"%")
    }

    // Filter by Condition
    if status != "" {
        query = query.Where("condition_status = ?", status)
    }

    // Eksekusi
    query.Find(&assets)

    c.JSON(http.StatusOK, gin.H{"data": assets})
}

// GetAssetByID (Tidak berubah)
func GetAssetByID(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.Preload("Category").Preload("Department").Preload("Room").Preload("CreatedBy").Preload("UpdatedBy").First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// CreateAsset MENAMBAHKAN Upload File
func CreateAsset(c *gin.Context) {
	var input models.Asset

	// 1. Ubah Binding dari ShouldBindJSON ke ShouldBind agar bisa baca Form Data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Proses Upload File
	// "photo" adalah nama key di form-data (postman/curl)
	file, err := c.FormFile("photo")
	if err == nil {
		// Pastikan folder uploads/assets ada
		path := "uploads/assets"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}

		// Buat nama file unik (timestamp + nama asli)
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filepath := fmt.Sprintf("%s/%s", path, filename)

		// Simpan file ke server
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Simpan path ke struct input
		input.PhotoPath = filepath
	}

	// 3. Set User Tracking
	if userID, exists := c.Get("userID"); exists {
		input.CreatedByID = userID.(uint)
		input.UpdatedByID = userID.(uint)
	}

	// 4. Simpan ke Database
	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	RecordLog(c, "CREATE", "assets", input.ID, "Menambahkan aset baru: "+input.Name)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

// UpdateAsset MENAMBAHKAN kemampuan ganti foto
func UpdateAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}

	var input models.Asset
	// Gunakan ShouldBind untuk Form Data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah user mengupload foto baru
	file, err := c.FormFile("photo")
	if err == nil {
		// Logic simpan file sama seperti Create
		path := "uploads/assets"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		newFilePath := fmt.Sprintf("%s/%s", path, filename)

		if err := c.SaveUploadedFile(file, newFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Hapus foto lama jika ada (optional, agar tidak menuhin server)
		if asset.PhotoPath != "" {
			os.Remove(asset.PhotoPath)
		}

		// Update path di struct input (agar tersimpan ke DB)
		input.PhotoPath = newFilePath
	}

	if userID, exists := c.Get("userID"); exists {
		input.UpdatedByID = userID.(uint)
	}

	models.DB.Model(&asset).Updates(input)

	RecordLog(c, "UPDATE", "assets", asset.ID, "Mengupdate data aset")

	c.JSON(http.StatusOK, gin.H{"data": asset})
}

// DeleteAsset (Tidak berubah)
func DeleteAsset(c *gin.Context) {
	var asset models.Asset
	if err := models.DB.First(&asset, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found!"})
		return
	}
	
	// Hapus file fisik foto jika ada
	if asset.PhotoPath != "" {
		os.Remove(asset.PhotoPath)
	}

	models.DB.Delete(&asset)
	RecordLog(c, "DELETE", "assets", asset.ID, "Menghapus aset: "+asset.Name)
	c.JSON(http.StatusOK, gin.H{"data": true})
}