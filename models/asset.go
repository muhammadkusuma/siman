package models

import (
	"time"
)

// Asset adalah tabel utama yang menyimpan data barang inventaris.
type Asset struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Aset"`

	// --- Identitas Aset ---
	InventoryCode   string `json:"inventory_code" gorm:"size:100;unique;comment:Kode Barcode/QR Internal Kampus (cth: UIN-2024-001)"`
	NUP             int    `json:"nup" gorm:"comment:Nomor Urut Pendaftaran BMN (cth: 1)"`
	AssetCategoryID uint   `json:"asset_category_id" gorm:"comment:Foreign Key ke Kategori BMN"`

	// --- Detail Fisik ---
	Name         string `json:"name" gorm:"size:255;not null;comment:Nama Spesifik Barang (cth: Laptop Dell XPS 13)"`
	Brand        string `json:"brand" gorm:"size:100;comment:Merk Pabrikan (cth: Dell, Samsung)"`
	Model        string `json:"model" gorm:"size:100;comment:Tipe/Model Barang (cth: XPS 13 9310)"`
	SerialNumber string `json:"serial_number" gorm:"size:100;comment:Nomor Seri Pabrikan (cth: 8H7G6F5)"`

	// --- Status & Kondisi ---
	ConditionStatus   string `json:"condition_status" gorm:"default:'Baik';comment:Kondisi Fisik (cth: Baik, Rusak Ringan, Rusak Berat)"`
	OperationalStatus string `json:"operational_status" gorm:"default:'Aktif';comment:Status Operasional (cth: Aktif, Dipinjam, Dalam Perbaikan, Dihapuskan)"`

	// --- Lokasi & Kepemilikan ---
	DepartmentID uint `json:"department_id" gorm:"comment:Foreign Key Unit Pemilik Barang"`
	RoomID       uint `json:"room_id" gorm:"comment:Foreign Key Lokasi Ruangan Sekarang"`

	// --- Data Keuangan ---
	AcquisitionDate time.Time `json:"acquisition_date" gorm:"comment:Tanggal Pembelian/Perolehan (cth: 2023-01-20)"`
	Price           float64   `json:"price" gorm:"type:decimal(15,2);comment:Nilai Perolehan dalam Rupiah (cth: 15000000.00)"`
	SourceOfFund    string    `json:"source_of_fund" gorm:"size:50;comment:Sumber Dana Pengadaan (cth: APBN, PNBP, Hibah)"`
	PurchaseOrder   string    `json:"purchase_order" gorm:"size:100;comment:Nomor SPK/Kontrak/Kuitansi (cth: SPK/01/UIN/2023)"`

	// --- TRACKING USER ---
	CreatedByID uint `json:"created_by_id" gorm:"comment:ID User yang pertama kali input data ini"`
	UpdatedByID uint `json:"updated_by_id" gorm:"comment:ID User yang terakhir kali mengedit data ini"`

	CreatedAt time.Time `json:"created_at" gorm:"comment:Waktu data dibuat"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:Waktu data terakhir diubah"`

	// --- Relasi ---
	Category   AssetCategory `json:"category" gorm:"foreignKey:AssetCategoryID"`
	Department Department    `json:"department" gorm:"foreignKey:DepartmentID"`
	Room       Room          `json:"room" gorm:"foreignKey:RoomID"`
	
	CreatedBy User `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy User `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}