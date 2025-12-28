package models

import (
	"time"
)

// MutationLog mencatat sejarah perpindahan aset antar unit/ruangan.
type MutationLog struct {
	ID      uint `json:"id" gorm:"primaryKey;autoIncrement;comment:ID Log Mutasi"`
	AssetID uint `json:"asset_id" gorm:"not null;comment:ID Aset yang dipindahkan"`

	// Asal
	FromDepartmentID *uint `json:"from_department_id" gorm:"comment:Unit Asal (Bisa NULL jika barang baru)"`
	FromRoomID       *uint `json:"from_room_id" gorm:"comment:Ruangan Asal (Bisa NULL jika barang baru)"`

	// Tujuan
	ToDepartmentID uint `json:"to_department_id" gorm:"not null;comment:Unit Tujuan pemindahan"`
	ToRoomID       uint `json:"to_room_id" gorm:"comment:Ruangan Tujuan pemindahan"`

	MutationDate time.Time `json:"mutation_date" gorm:"default:CURRENT_TIMESTAMP;comment:Waktu mutasi terjadi"`
	ApprovedBy   string    `json:"approved_by" gorm:"size:255;comment:Nama Pejabat yang menyetujui mutasi (cth: Kabag Umum)"`
	Reason       string    `json:"reason" gorm:"type:text;comment:Alasan mutasi (cth: Peminjaman untuk acara seminar)"`

	// Relasi
	Asset          Asset       `json:"asset" gorm:"foreignKey:AssetID"`
	FromDepartment *Department `json:"from_department" gorm:"foreignKey:FromDepartmentID"`
	ToDepartment   Department  `json:"to_department" gorm:"foreignKey:ToDepartmentID"`
	FromRoom       *Room       `json:"from_room" gorm:"foreignKey:FromRoomID"`
	ToRoom         Room        `json:"to_room" gorm:"foreignKey:ToRoomID"`
}