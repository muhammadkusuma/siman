package models

import (
	"time"
)

// ===================================================================================
// 1. MASTER DATA ORGANISASI (STRUKTUR KAMPUS)
// ===================================================================================

// Faculty merepresentasikan induk unit: Fakultas, Sekolah Pascasarjana, atau Direktorat.
type Faculty struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Unik Fakultas/Unit Induk"`
	Code      string    `json:"code" gorm:"size:50;unique;comment:Kode Singkatan (misal: FT, SPS, REK)"`
	Name      string    `json:"name" gorm:"size:255;not null;comment:Nama Lengkap Fakultas/Unit"`
	Type      string    `json:"type" gorm:"type:enum('Fakultas','Sekolah Pascasarjana','Direktorat','Lembaga');not null;comment:Jenis Institusi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi: Satu Fakultas memiliki banyak Prodi/Departemen
	Departments []Department `json:"departments" gorm:"foreignKey:FacultyID"`
}

// Department merepresentasikan Prodi (S1/S2/S3) atau Bagian/Biro spesifik.
type Department struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FacultyID uint   `json:"faculty_id" gorm:"not null;comment:Foreign Key ke tabel Faculties"`
	
	Code      string `json:"code" gorm:"size:50;comment:Kode Prodi/Unit (misal: IF-S1)"`
	Name      string `json:"name" gorm:"size:255;not null;comment:Nama Prodi (misal: S1 Teknik Informatika)"`
	
	// Enum StudyLevel penting untuk membedakan aset Pascasarjana vs Sarjana
	StudyLevel string `json:"study_level" gorm:"type:enum('D3','S1','S2','S3','Profesi','Non-Akademik');default:'Non-Akademik';comment:Jenjang Studi. Pilih Non-Akademik untuk kantor admin."`
	
	HeadOfDepartment string `json:"head_of_department" gorm:"size:255;comment:Nama Kaprodi atau Kepala Unit"`

	// Relasi Belongs To
	Faculty Faculty `json:"faculty" gorm:"foreignKey:FacultyID"`
}

// ===================================================================================
// 2. MASTER DATA LOKASI & KATEGORI (STANDAR BMN)
// ===================================================================================

// Building menyimpan data Gedung fisik.
type Building struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string `json:"code" gorm:"size:50;comment:Kode Gedung (misal: G1)"`
	Name        string `json:"name" gorm:"size:255;comment:Nama Gedung (misal: Gedung Rektorat)"`
	TotalFloors int    `json:"total_floors" gorm:"comment:Jumlah Lantai dalam gedung"`

	// Relasi: Satu Gedung punya banyak Ruangan
	Rooms []Room `json:"rooms" gorm:"foreignKey:BuildingID"`
}

// Room menyimpan detail ruangan di dalam gedung.
type Room struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	BuildingID uint   `json:"building_id" gorm:"not null;comment:Foreign Key ke tabel Buildings"`
	
	RoomNumber string `json:"room_number" gorm:"size:50;comment:Nomor Pintu (misal: A.204)"`
	Name       string `json:"name" gorm:"size:255;comment:Fungsi Ruangan (misal: Lab Komputer)"`
	Floor      int    `json:"floor" gorm:"comment:Posisi Lantai"`

	// Relasi Belongs To
	Building Building `json:"building" gorm:"foreignKey:BuildingID"`
}

// AssetCategory menyimpan kode klasifikasi BMN (Barang Milik Negara).
type AssetCategory struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	KodeBarang  string `json:"kode_barang" gorm:"size:50;unique;comment:Kode Klasifikasi BMN (misal: 3.10.01.02.001)"`
	Name        string `json:"name" gorm:"size:255;not null;comment:Nama Golongan Barang"`
	Description string `json:"description" gorm:"type:text;comment:Keterangan detail kategori"`
}

// ===================================================================================
// 3. DATA UTAMA ASET (INVENTARIS)
// ===================================================================================

// Asset adalah tabel utama yang menyimpan data barang.
type Asset struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	// --- Identitas Aset ---
	InventoryCode   string `json:"inventory_code" gorm:"size:100;unique;comment:Kode Barcode/QR Internal Kampus"`
	NUP             int    `json:"nup" gorm:"comment:Nomor Urut Pendaftaran (Wajib BMN)"`
	AssetCategoryID uint   `json:"asset_category_id" gorm:"comment:Foreign Key ke Kategori BMN"`

	// --- Detail Fisik ---
	Name         string `json:"name" gorm:"size:255;not null;comment:Nama Spesifik Barang"`
	Brand        string `json:"brand" gorm:"size:100;comment:Merk"`
	Model        string `json:"model" gorm:"size:100;comment:Tipe/Model"`
	SerialNumber string `json:"serial_number" gorm:"size:100;comment:SN Pabrikan"`

	// --- Status & Kondisi ---
	ConditionStatus   string `json:"condition_status" gorm:"default:'Baik';comment:Kondisi Fisik: Baik, Rusak Ringan, Rusak Berat"`
	OperationalStatus string `json:"operational_status" gorm:"default:'Aktif';comment:Status: Aktif, Dipinjam, Dalam Perbaikan, Penghapusan"`

	// --- Lokasi & Kepemilikan ---
	DepartmentID uint `json:"department_id" gorm:"comment:Milik Prodi/Unit mana"`
	RoomID       uint `json:"room_id" gorm:"comment:Posisi sekarang di ruangan mana"`

	// --- Data Keuangan & Perolehan ---
	AcquisitionDate time.Time `json:"acquisition_date" gorm:"comment:Tanggal Pembelian/Perolehan"`
	Price           float64   `json:"price" gorm:"type:decimal(15,2);comment:Nilai Perolehan (Rupiah)"`
	SourceOfFund    string    `json:"source_of_fund" gorm:"size:50;comment:Sumber Dana (APBN/PNBP/Hibah)"`
	PurchaseOrder   string    `json:"purchase_order" gorm:"size:100;comment:Nomor SPK/Kontrak"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// --- Relasi (Preloading) ---
	Category   AssetCategory `json:"category" gorm:"foreignKey:AssetCategoryID"`
	Department Department    `json:"department" gorm:"foreignKey:DepartmentID"`
	Room       Room          `json:"room" gorm:"foreignKey:RoomID"`
}

// ===================================================================================
// 4. DATA RIWAYAT (LOGS)
// ===================================================================================

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

// MutationLog mencatat sejarah perpindahan aset antar unit/ruangan.
type MutationLog struct {
	ID      uint `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetID uint `json:"asset_id" gorm:"not null"`

	// Relasi Asal (Pointer *uint karena bisa NULL jika barang baru input)
	FromDepartmentID *uint `json:"from_department_id" gorm:"comment:Unit Asal"`
	FromRoomID       *uint `json:"from_room_id" gorm:"comment:Ruangan Asal"`

	// Relasi Tujuan (Tidak boleh NULL)
	ToDepartmentID uint `json:"to_department_id" gorm:"not null;comment:Unit Tujuan"`
	ToRoomID       uint `json:"to_room_id" gorm:"comment:Ruangan Tujuan"`

	MutationDate time.Time `json:"mutation_date" gorm:"default:CURRENT_TIMESTAMP"`
	ApprovedBy   string    `json:"approved_by" gorm:"size:255;comment:Nama Pejabat yang menyetujui"`
	Reason       string    `json:"reason" gorm:"type:text;comment:Alasan mutasi"`

	// Relasi Struct
	Asset          Asset       `json:"asset" gorm:"foreignKey:AssetID"`
	FromDepartment *Department `json:"from_department" gorm:"foreignKey:FromDepartmentID"`
	ToDepartment   Department  `json:"to_department" gorm:"foreignKey:ToDepartmentID"`
	FromRoom       *Room       `json:"from_room" gorm:"foreignKey:FromRoomID"`
	ToRoom         Room        `json:"to_room" gorm:"foreignKey:ToRoomID"`
}