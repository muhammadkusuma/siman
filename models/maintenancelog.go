package models

import (
	"time"
)

// MaintenanceLog mencatat sejarah perbaikan aset (Service History).
type MaintenanceLog struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Log Perbaikan"`
	AssetID     uint      `json:"asset_id" gorm:"not null;comment:ID Aset yang rusak"`
	
	IssueDate   time.Time `json:"issue_date" gorm:"comment:Tanggal kerusakan dilaporkan (cth: 2024-05-10)"`
	Description string    `json:"description" gorm:"type:text;comment:Detail keluhan kerusakan (cth: Layar berkedip dan mati total)"`
	ActionTaken string    `json:"action_taken" gorm:"type:text;comment:Tindakan perbaikan yang dilakukan (cth: Penggantian LCD Panel)"`
	Cost        float64   `json:"cost" gorm:"type:decimal(15,2);comment:Total biaya perbaikan (cth: 1500000.00)"`
	VendorName  string    `json:"vendor_name" gorm:"size:255;comment:Pihak ketiga/Bengkel yang mengerjakan (cth: Service Center Asus)"`
	Status      string    `json:"status" gorm:"default:'Pending';comment:Status Tiket: Pending, Proses, Selesai"`

	// Relasi
	Asset Asset `json:"asset" gorm:"foreignKey:AssetID"`
}