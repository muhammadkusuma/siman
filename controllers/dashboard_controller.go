package controllers

import (
	"muhammadkusuma/siman/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats memberikan ringkasan data untuk halaman Home/Dashboard
func GetDashboardStats(c *gin.Context) {
	var totalAssets int64
	var totalValue float64
	var totalMaintenance int64
	var totalMutations int64

	// 1. Hitung Total Aset
	models.DB.Model(&models.Asset{}).Count(&totalAssets)

	// 2. Hitung Total Nilai Aset (Sum Price)
	// Menggunakan Scan untuk memasukkan hasil query ke variabel pointer
	models.DB.Model(&models.Asset{}).Select("COALESCE(SUM(price), 0)").Scan(&totalValue)

	// 3. Hitung Aset yang sedang dalam Perbaikan (Status != Selesai)
	models.DB.Model(&models.MaintenanceLog{}).Where("status != ?", "Selesai").Count(&totalMaintenance)

	// 4. Hitung Total Mutasi bulan ini (Opsional, contoh filter date)
	models.DB.Model(&models.MutationLog{}).Count(&totalMutations)

	// 5. Statistik Kondisi Aset (Group By)
	// Contoh Output: [{"condition_status": "Baik", "total": 50}, {"condition_status": "Rusak", "total": 5}]
	type ConditionStat struct {
		ConditionStatus string `json:"condition_status"`
		Total           int    `json:"total"`
	}
	var conditionStats []ConditionStat
	models.DB.Model(&models.Asset{}).Select("condition_status, count(*) as total").Group("condition_status").Scan(&conditionStats)

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard data fetched",
		"data": gin.H{
			"total_assets":      totalAssets,
			"total_asset_value": totalValue,
			"active_maintenance": totalMaintenance,
			"total_mutations":   totalMutations,
			"asset_conditions":  conditionStats,
		},
	})
}