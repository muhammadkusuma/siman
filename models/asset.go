package models

import (
	"time"
)

// Asset adalah tabel utama yang menyimpan data barang.
type Asset struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	// --- Identitas Aset ---
	InventoryCode   string `json:"inventory_code" gorm:"size:100;unique;comment:Kode Barcode/QR Internal Kampus"`
	NUP             int    `json:"nup" gorm:"comment:Nomor Urut Pendaftaran"`
	AssetCategoryID uint   `json:"asset_category_id"`

	// --- Detail Fisik ---
	Name         string `json:"name" gorm:"size:255;not null"`
	Brand        string `json:"brand" gorm:"size:100"`
	Model        string `json:"model" gorm:"size:100"`
	SerialNumber string `json:"serial_number" gorm:"size:100"`

	// --- Status & Kondisi ---
	ConditionStatus   string `json:"condition_status" gorm:"default:'Baik'"`
	OperationalStatus string `json:"operational_status" gorm:"default:'Aktif'"`

	// --- Lokasi & Kepemilikan ---
	DepartmentID uint `json:"department_id"`
	RoomID       uint `json:"room_id"`

	// --- Data Keuangan ---
	AcquisitionDate time.Time `json:"acquisition_date"`
	Price           float64   `json:"price" gorm:"type:decimal(15,2)"`
	SourceOfFund    string    `json:"source_of_fund" gorm:"size:50"`
	PurchaseOrder   string    `json:"purchase_order" gorm:"size:100"`

	// --- TRACKING USER (Field Baru) ---
	CreatedByID uint `json:"created_by_id" gorm:"comment:User yang pertama kali input"`
	UpdatedByID uint `json:"updated_by_id" gorm:"comment:User yang terakhir kali edit"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// --- Relasi ---
	Category   AssetCategory `json:"category" gorm:"foreignKey:AssetCategoryID"`
	Department Department    `json:"department" gorm:"foreignKey:DepartmentID"`
	Room       Room          `json:"room" gorm:"foreignKey:RoomID"`
	
	// Relasi ke User (Tracking)
	CreatedBy User `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy User `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}