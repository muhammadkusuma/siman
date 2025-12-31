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

// UpdateMaintenance: Mengubah data perbaikan
func UpdateMaintenance(c *gin.Context) {
	id := c.Param("id")
	var maintenance models.MaintenanceLog

	// 1. Cari Data Lama
	if err := models.DB.First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance log not found"})
		return
	}

	// 2. Validasi Input JSON
	var input models.MaintenanceLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update Field
	maintenance.AssetID = input.AssetID
	maintenance.IssueDate = input.IssueDate
	maintenance.Status = input.Status
	maintenance.Description = input.Description
	maintenance.VendorName = input.VendorName
	maintenance.Cost = input.Cost
	maintenance.ActionTaken = input.ActionTaken

	// 4. Simpan Perubahan
	if err := models.DB.Save(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance log"})
		return
	}

	// 5. OTOMATIS UPDATE STATUS ASET
	// Jika status maintenance berubah, update juga status operasional aset
	if maintenance.Status == "Proses" {
		models.DB.Model(&models.Asset{}).Where("id = ?", maintenance.AssetID).Update("operational_status", "Dalam Perbaikan")
	} else if maintenance.Status == "Selesai" {
		models.DB.Model(&models.Asset{}).Where("id = ?", maintenance.AssetID).Update("operational_status", "Aktif")
	}

	c.JSON(http.StatusOK, gin.H{"data": maintenance, "message": "Maintenance updated successfully"})
}

// DeleteMaintenance: Menghapus data perbaikan
func DeleteMaintenance(c *gin.Context) {
	id := c.Param("id")
	var maintenance models.MaintenanceLog

	if err := models.DB.First(&maintenance, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance log not found"})
		return
	}

	if err := models.DB.Delete(&maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete maintenance log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance log deleted successfully"})
}
