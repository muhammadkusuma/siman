package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAuditLogs (Hanya Admin/User login yang bisa lihat)
func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	models.DB.Preload("User").Order("created_at desc").Find(&logs)
	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func GetAuditLogByID(c *gin.Context) {
	var log models.AuditLog
	if err := models.DB.Preload("User").First(&log, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": log})
}

// --- SMART HELPER FUNCTION ---
// Sekarang tidak butuh parameter 'userID' karena diambil otomatis dari token
func RecordLog(c *gin.Context, action string, tableName string, recordID uint, changes string) {
	
	// 1. Ambil User ID dari Context (hasil set dari Middleware JWT)
	userIDVal, exists := c.Get("userID")
	var userID uint
	if exists {
		userID = userIDVal.(uint)
	} else {
		// Jika tidak ada user (misal sistem background), set 0 atau handle error
		userID = 0 
	}

	// 2. Ambil Info Tambahan
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 3. Simpan Log
	log := models.AuditLog{
		UserID:    userID,
		Action:    action,
		TableName: tableName,
		RecordID:  recordID,
		Changes:   changes,
		IPAddress: ip,
		UserAgent: userAgent,
	}

	// Eksekusi di Goroutine agar cepat
	go func() {
		models.DB.Create(&log)
	}()
}