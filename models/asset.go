package models

import (
	"time"
)

// Asset adalah tabel utama yang menyimpan data barang inventaris.
type Asset struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Aset"`

	// --- Identitas Aset ---
	InventoryCode   string `json:"inventory_code" gorm:"size:100;unique;comment:Kode Barcode/QR Internal Kampus"`
	NUP             int    `json:"nup" gorm:"comment:Nomor Urut Pendaftaran BMN"`
	AssetCategoryID uint   `json:"asset_category_id" gorm:"comment:Foreign Key ke Kategori BMN"`

	// --- Detail Fisik ---
	Name         string `json:"name" gorm:"size:255;not null;comment:Nama Spesifik Barang"`
	Brand        string `json:"brand" gorm:"size:100;comment:Merk Pabrikan"`
	Model        string `json:"model" gorm:"size:100;comment:Tipe/Model Barang"`
	SerialNumber string `json:"serial_number" gorm:"size:100;comment:Nomor Seri Pabrikan"`

	// --- FOTO ASET (BARU) ---
	PhotoPath string `json:"photo_path" gorm:"size:255;comment:Lokasi file foto"`

	// --- Status & Kondisi ---
	ConditionStatus   string `json:"condition_status" gorm:"default:'Baik';comment:Kondisi Fisik"`
	OperationalStatus string `json:"operational_status" gorm:"default:'Aktif';comment:Status Operasional"`

	// --- Lokasi & Kepemilikan ---
	DepartmentID uint `json:"department_id" gorm:"comment:Foreign Key Unit Pemilik"`
	RoomID       uint `json:"room_id" gorm:"comment:Foreign Key Lokasi Ruangan"`

	// --- Data Keuangan ---
	AcquisitionDate time.Time `json:"acquisition_date" gorm:"comment:Tanggal Pembelian"`
	Price           float64   `json:"price" gorm:"type:decimal(15,2);comment:Nilai Perolehan"`
	SourceOfFund    string    `json:"source_of_fund" gorm:"size:50;comment:Sumber Dana"`
	PurchaseOrder   string    `json:"purchase_order" gorm:"size:100;comment:Nomor SPK"`

	// --- TRACKING USER ---
	CreatedByID uint `json:"created_by_id" gorm:"comment:ID User yang input"`
	UpdatedByID uint `json:"updated_by_id" gorm:"comment:ID User yang edit"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// --- Relasi ---
	Category   AssetCategory `json:"category" gorm:"foreignKey:AssetCategoryID"`
	Department Department    `json:"department" gorm:"foreignKey:DepartmentID"`
	Room       Room          `json:"room" gorm:"foreignKey:RoomID"`

	CreatedBy User `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy User `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}
