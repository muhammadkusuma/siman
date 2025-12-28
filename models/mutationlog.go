package models

import (
	"time"
)

// MutationLog mencatat sejarah perpindahan aset antar unit/ruangan.
type MutationLog struct {
	ID      uint `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Log Mutasi"`
	AssetID uint `json:"asset_id" gorm:"not null;comment:ID Aset yang dipindahkan"`

	// --- ASAL (Bisa NULL) ---
	// Kita gunakan pointer *uint agar database menerima NULL.
	FromDepartmentID *uint `json:"from_department_id" gorm:"default:null;comment:Unit Asal"`
	FromRoomID       *uint `json:"from_room_id" gorm:"default:null;comment:Ruangan Asal"`

	// --- TUJUAN (Wajib Ada) ---
	ToDepartmentID uint `json:"to_department_id" gorm:"not null;comment:Unit Tujuan"`
	ToRoomID       uint `json:"to_room_id" gorm:"not null;comment:Ruangan Tujuan"`

	MutationDate time.Time `json:"mutation_date" gorm:"autoCreateTime;comment:Waktu mutasi terjadi"`
	ApprovedBy   string    `json:"approved_by" gorm:"size:255;comment:Nama Pejabat yang menyetujui"`
	Reason       string    `json:"reason" gorm:"type:text;comment:Alasan mutasi"`

	// --- RELASI DENGAN CONSTRAINT EKSPLISIT ---
	
	// Jika Aset dihapus, log ikut terhapus (CASCADE)
	Asset Asset `json:"asset" gorm:"foreignKey:AssetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Jika Unit Asal dihapus, set kolom ini jadi NULL (agar log tetap ada)
	FromDepartment *Department `json:"from_department" gorm:"foreignKey:FromDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	
	// Jika Ruangan Asal dihapus, set NULL
	FromRoom       *Room       `json:"from_room" gorm:"foreignKey:FromRoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Jika Unit Tujuan dihapus, tolak penghapusan (RESTRICT) karena masih tercatat di sini
	ToDepartment   Department  `json:"to_department" gorm:"foreignKey:ToDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	
	// Jika Ruangan Tujuan dihapus, tolak penghapusan
	ToRoom         Room        `json:"to_room" gorm:"foreignKey:ToRoomID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}