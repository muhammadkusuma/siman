package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- PUBLIC API ENDPOINTS (READ ONLY) ---

// GetAuditLogs mengambil daftar sejarah aktivitas sistem
func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	
	// Preload User agar kita tahu nama user yang melakukan aksi
	// Order by created_at desc agar log terbaru muncul paling atas
	models.DB.Preload("User").Order("created_at desc").Find(&logs)

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetAuditLogByID melihat detail satu log spesifik
func GetAuditLogByID(c *gin.Context) {
	var log models.AuditLog
	if err := models.DB.Preload("User").First(&log, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": log})
}

// --- HELPER FUNCTION (INTERNAL USE) ---
// Fungsi ini tidak dijadikan route API, tapi dipanggil oleh Controller lain
func RecordLog(c *gin.Context, userID uint, action string, tableName string, recordID uint, changes string) {
	// Ambil IP dan UserAgent otomatis dari Context Gin
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	log := models.AuditLog{
		UserID:    userID,
		Action:    action,
		TableName: tableName,
		RecordID:  recordID,
		Changes:   changes, // JSON string perubahan data
		IPAddress: ip,
		UserAgent: userAgent,
	}

	// Simpan ke database di background (goroutine) agar tidak memperlambat response user
	go func() {
		models.DB.Create(&log)
	}()
}