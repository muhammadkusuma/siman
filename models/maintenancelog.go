package models

import (
	"time"
)

// MaintenanceLog mencatat sejarah perbaikan aset (Service History).
type MaintenanceLog struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID     uint      `json:"asset_id" gorm:"not null"`
	
	IssueDate   time.Time `json:"issue_date" gorm:"comment:Tanggal kerusakan dilaporkan"`
	Description string    `json:"description" gorm:"type:text;comment:Keluhan kerusakan"`
	ActionTaken string    `json:"action_taken" gorm:"type:text;comment:Tindakan perbaikan yang dilakukan"`
	Cost        float64   `json:"cost" gorm:"type:decimal(15,2);comment:Biaya perbaikan"`
	VendorName  string    `json:"vendor_name" gorm:"size:255;comment:Pihak ketiga/Bengkel"`
	Status      string    `json:"status" gorm:"default:'Pending';comment:Status Perbaikan: Pending, Proses, Selesai"`

	// Relasi ke Aset
	Asset Asset `json:"asset" gorm:"foreignKey:AssetID"`
}