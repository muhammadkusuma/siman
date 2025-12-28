package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
)

// --- MUTATION LOGIC ---

// CreateMutation mencatat perpindahan aset dan MENGUPDATE lokasi aset secara otomatis
func CreateMutation(c *gin.Context) {
	var input models.MutationLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mulai Transaksi Database (agar aman jika error di tengah jalan)
	tx := models.DB.Begin()

	// 1. Ambil data aset saat ini untuk mengisi 'FromDepartment' dan 'FromRoom' jika belum diisi
	var asset models.Asset
	if err := tx.First(&asset, input.AssetID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// Isi data asal dari posisi aset sekarang
	input.FromDepartmentID = &asset.DepartmentID
	input.FromRoomID = &asset.RoomID
	input.MutationDate = time.Now()

	// 2. Simpan Log Mutasi
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mutation log: " + err.Error()})
		return
	}

	// 3. Update Lokasi Aset ke Lokasi Baru (ToDepartment & ToRoom)
	if err := tx.Model(&asset).Updates(map[string]interface{}{
		"department_id": input.ToDepartmentID,
		"room_id":       input.ToRoomID,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset location"})
		return
	}

	// Commit Transaksi
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Mutation successful", "data": input})
}

func GetMutations(c *gin.Context) {
	var mutations []models.MutationLog
	models.DB.Preload("Asset").Preload("FromDepartment").Preload("ToDepartment").Find(&mutations)
	c.JSON(http.StatusOK, gin.H{"data": mutations})
}

// --- MAINTENANCE LOGIC ---

// CreateMaintenance mencatat perbaikan dan bisa update status aset
func CreateMaintenance(c *gin.Context) {
	var input models.MaintenanceLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.DB.Begin()

	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Opsional: Jika status 'Proses', ubah status aset jadi 'Dalam Perbaikan'
	if input.Status == "Proses" {
		tx.Model(&models.Asset{}).Where("id = ?", input.AssetID).Update("operational_status", "Dalam Perbaikan")
	} else if input.Status == "Selesai" {
		tx.Model(&models.Asset{}).Where("id = ?", input.AssetID).Update("operational_status", "Aktif")
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": input})
}

func GetMaintenances(c *gin.Context) {
	var logs []models.MaintenanceLog
	models.DB.Preload("Asset").Find(&logs)
	c.JSON(http.StatusOK, gin.H{"data": logs})
}